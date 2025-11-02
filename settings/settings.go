package settings

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/url"
	"strings"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-aws/sqs"
	watermillSqs "github.com/ThreeDotsLabs/watermill-aws/sqs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	awsCreds "github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/go-playground/validator/v10"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	_ "embed"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	amazonsqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	transport "github.com/aws/smithy-go/endpoints"
	wotelfloss "github.com/dentech-floss/watermill-opentelemetry-go-extra/pkg/opentelemetry"
	"github.com/spf13/viper"
	wotel "github.com/voi-oss/watermill-opentelemetry/pkg/opentelemetry"
)

//go:embed base.yaml
var BaseSettings []byte

type CORSSettings struct {
	Origins []string `mapstructure:"origins" validate:"min=1,dive,url"`
	Methods []string `mapstructure:"methods" validate:"min=1,dive,oneof=GET POST PUT DELETE OPTIONS PATCH HEAD"`
	Headers []string `mapstructure:"headers" validate:"min=1"`
}

type HTTPSettings struct {
	Port             string       `mapstructure:"port" validate:"required,numeric"`
	Prefix           string       `mapstructure:"prefix" validate:"required"`
	IP               string       `mapstructure:"ip" validate:"required,ip"`
	CORS             CORSSettings `mapstructure:"cors" validate:"required"`
	Timeout          int          `mapstructure:"timeout" validate:"gte=1"`
	SwaggerUIEnabled bool         `mapstructure:"swagger-ui-enabled"`
}

type OpenTelemetryLogSettings struct {
	TimeoutInSec  int64 `mapstructure:"timeout"`
	IntervalInSec int64 `mapstructure:"interval"`
	MaxQueueSize  int   `mapstructure:"maxqueuesize"`
	BatchSize     int   `mapstructure:"batchsize"`
}

type OpenTelemetryTraceSettings struct {
	TimeoutInSec int64 `mapstructure:"timeout"`
	MaxQueueSize int   `mapstructure:"maxqueuesize"`
	BatchSize    int   `mapstructure:"batchsize"`
	SampleRate   int   `mapstructure:"samplerate"`
}

type OpenTelemetryMetricSettings struct {
	IntervalInSec int64 `mapstructure:"interval"`
	TimeoutInSec  int64 `mapstructure:"timeout"`
}

type OpenTelemetrySettings struct {
	Enabled  bool                        `mapstructure:"enabled"`
	Endpoint string                      `mapstructure:"endpoint"`
	Metrics  OpenTelemetryMetricSettings `mapstructure:"metrics"`
	Traces   OpenTelemetryTraceSettings  `mapstructure:"traces"`
	Logs     OpenTelemetryLogSettings    `mapstructure:"logs"`
	Interval int                         `mapstructure:"interval"`
}

type ImageProcessorSettings struct {
	BucketName string `mapstructure:"bucket-name" validate:"required"`
}

type ObjectStorerSettings struct {
	Endpoint        string `mapstructure:"endpoint" validate:"required"`
	AccessKeyID     string `mapstructure:"access-key-id"`
	SecretAccessKey string `mapstructure:"secret-access-key"`
	UseSSL          bool   `mapstructure:"use-ssl"`
	Region          string `mapstructure:"region"`
}

// NewMinioClient creates a new MinIO client from the settings
func (o *ObjectStorerSettings) NewMinioClient(ctx context.Context) (*minio.Client, error) {
	// Initialize MinIO client
	var creds *credentials.Credentials

	// Use IAM credentials if connecting to AWS S3 or if credentials are empty
	if strings.Contains(o.Endpoint, ".amazonaws.com") || (o.AccessKeyID == "" && o.SecretAccessKey == "") {
		slog.InfoContext(ctx, "Using IAM for connecting to s3")
		creds = credentials.NewIAM("")
	} else {
		creds = credentials.NewStaticV4(o.AccessKeyID, o.SecretAccessKey, "")
	}

	minioClient, err := minio.New(o.Endpoint, &minio.Options{
		Creds:  creds,
		Secure: o.UseSSL,
		Region: o.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to initialize MinIO client: %w", err)
	}

	// Verify connection by listing buckets
	_, err = minioClient.ListBuckets(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to MinIO: %w", err)
	}

	return minioClient, nil
}

type DatabaseSettings struct {
	Host                   string `mapstructure:"host" validate:"required"`
	Port                   int    `mapstructure:"port" validate:"required,gte=1,lte=65535"`
	User                   string `mapstructure:"user" validate:"required"`
	Password               string `mapstructure:"password" validate:"required"`
	Database               string `mapstructure:"database" validate:"required"`
	Schema                 string `mapstructure:"schema"`
	SSLMode                string `mapstructure:"ssl-mode" validate:"oneof=disable require verify-ca verify-full"`
	MaxOpenConns           int    `mapstructure:"max-open-conns" validate:"gte=1"`
	MaxIdleConns           int    `mapstructure:"max-idle-conns" validate:"gte=1"`
	ConnMaxLifetimeMinutes int    `mapstructure:"conn-max-lifetime-minutes" validate:"gte=1"`
}

// BuildConnectionString builds a PostgreSQL connection string from the settings
func (d *DatabaseSettings) BuildConnectionString() string {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.Database,
	)

	// Add query parameters
	params := []string{}

	if d.SSLMode != "" {
		params = append(params, fmt.Sprintf("sslmode=%s", d.SSLMode))
	}

	if d.Schema != "" {
		params = append(params, fmt.Sprintf("search_path=%s", d.Schema))
	}

	if len(params) > 0 {
		connStr += "?" + strings.Join(params, "&")
	}

	return connStr
}

type AWSSettings struct {
	Endpoint string `mapstructure:"endpoint" validate:"url"`
	// AccessKeyID and SecretAccessKey are used for explicit key-based authentication.
	// If both are empty, the application should assume an implicit credential
	// mechanism (e.g., EC2 IAM role, ECS task role, EKS Service Account).
	AccessKey string `mapstructure:"access-key"`
	SecretKey string `mapstructure:"secret-key"`
	Region    string `json:"region"`
	Anonymous bool   `mapstructure:"anonymous"`
}

func (a *AWSSettings) NewAWSConfig() (awsConfig.Config, error) {
	if a.Anonymous {
		slog.InfoContext(context.TODO(), "Using anonymous AWS credentials provider")
		return awsConfig.LoadDefaultConfig(
			context.TODO(),
			awsConfig.WithRegion(a.Region),
			awsConfig.WithCredentialsProvider(aws.AnonymousCredentials{}),
		)
	}
	if a.AccessKey == "" && a.SecretKey == "" {
		slog.InfoContext(context.TODO(), "Using default AWS credentials provider chain")

		return awsConfig.LoadDefaultConfig(
			context.TODO(),
			awsConfig.WithRegion(a.Region),
		)
	}

	return awsConfig.LoadDefaultConfig(
		context.TODO(),
		awsConfig.WithCredentialsProvider(
			awsCreds.NewStaticCredentialsProvider(a.AccessKey, a.SecretKey, ""),
		),
		awsConfig.WithRegion(a.Region),
	)
}

func (a *AWSSettings) GetEndpointResolver() (func(*amazonsqs.Options), error) {
	endpoint, err := url.Parse(a.Endpoint)
	if err != nil {
		return nil, err
	}

	return amazonsqs.WithEndpointResolverV2(watermillSqs.OverrideEndpointResolver{
		Endpoint: transport.Endpoint{
			URI: *endpoint,
		},
	}), nil
}

type (
	WatermillPublisherSettings  struct{}
	WatermillSubscriberSettings struct{}
)

type WatermillBrokerSettings struct {
	Kind       string                      `mapstructure:"kind" validate:"required,oneof=sqs nats"`
	AWS        AWSSettings                 `mapstructure:"aws" validate:"required_if=Kind sqs"`
	Publisher  WatermillPublisherSettings  `mapstructure:"publisher" validate:"omitempty"`
	Subscriber WatermillSubscriberSettings `mapstructure:"subscriber" validate:"omitempty"`
}

func (broker *WatermillBrokerSettings) NewPublisher() (message.Publisher, error) {
	wattermilLogger := watermill.NewSlogLogger(slog.Default())

	var publisher message.Publisher

	switch broker.Kind {
	case "sqs":
		endpointResolver, err := broker.AWS.GetEndpointResolver()
		if err != nil {
			return nil, err
		}
		cfg, err := broker.AWS.NewAWSConfig()

		slog.InfoContext(context.TODO(), "casting aws config")
		castedConfig := (cfg).(aws.Config)

		optFns := make([]func(*amazonsqs.Options), 0)

		if broker.AWS.Endpoint != "" {
			optFns = append(optFns, endpointResolver)
		}

		sqsConfig := sqs.PublisherConfig{
			OptFns:    optFns,
			AWSConfig: castedConfig,
		}

		sqsPublisher, err := sqs.NewPublisher(sqsConfig, wattermilLogger)
		if err != nil {
			return nil, err
		}

		publisher = sqsPublisher
	case "nats":
		println("nats")
	}

	tracePropagatingPublisherDecorator := wotelfloss.NewTracePropagatingPublisherDecorator(publisher)
	return wotel.NewNamedPublisherDecorator("pubsub.Publish", tracePropagatingPublisherDecorator), nil
}

func (broker *WatermillBrokerSettings) NewSubscriber() (message.Subscriber, error) {
	wattermilLogger := watermill.NewSlogLogger(slog.Default())
	var subscriber message.Subscriber

	switch broker.Kind {
	case "sqs":
		endpointResolver, err := broker.AWS.GetEndpointResolver()
		if err != nil {
			return nil, err
		}
		cfg, err := broker.AWS.NewAWSConfig()

		slog.InfoContext(context.TODO(), "casting aws config")
		castedConfig := (cfg).(aws.Config)

		optFns := make([]func(*amazonsqs.Options), 0)

		if broker.AWS.Endpoint != "" {
			optFns = append(optFns, endpointResolver)
		}

		sqsConfig := sqs.SubscriberConfig{
			OptFns:    optFns,
			AWSConfig: castedConfig,
		}

		sqsSubs, err := sqs.NewSubscriber(sqsConfig, wattermilLogger)
		if err != nil {
			return nil, err
		}

		subscriber = sqsSubs
	case "nats":
		println("nats")
	}

	return subscriber, nil
}

type WatermillSettings struct {
	Broker     WatermillBrokerSettings `mapstructure:"broker" validate:"required"`
	ImageTopic string                  `mapstructure:"image-topic" validate:"required"`
}

type DynamoDBLogsSettings struct {
	Enabled bool        `mapstructure:"enabled"`
	Table   string      `mapstructure:"table" validate:"required_if=Enabled true"`
	AWS     AWSSettings `mapstructure:"aws" validate:"required_if=Enabled true"`
}

func (d *DynamoDBLogsSettings) NewDynamoDBClient() (*dynamodb.Client, error) {
	cfg, err := d.AWS.NewAWSConfig()
	if err != nil {
		return nil, err
	}

	awsCfg, ok := cfg.(aws.Config)
	if !ok {
		return nil, fmt.Errorf("expected to create aws config for dynamodb client")
	}

	client := dynamodb.NewFromConfig(awsCfg,
		func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String(d.AWS.Endpoint)
		},
	)
	return client, nil
}

type AppSettings struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Env     string `mapstructure:"env"`
}

func LoadConfig[T any](prefix string, baseConfig []byte) (*T, error) {
	var cfg *T

	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewReader(baseConfig))
	if err != nil {
		log.Println("Failed to read config from yaml")
		return nil, err
	}

	viper.SetEnvPrefix(prefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", ""))
	viper.AutomaticEnv()

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
