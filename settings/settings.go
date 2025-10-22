package settings

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/go-playground/validator/v10"

	_ "embed"

	"github.com/spf13/viper"
)

//go:embed base.yaml
var BaseSettings []byte

type CORSSettings struct {
	Origins []string `mapstructure:"origins" validate:"min=1,dive,url"`
	Methods []string `mapstructure:"methods" validate:"min=1,dive,oneof=GET POST PUT DELETE OPTIONS PATCH HEAD"`
	Headers []string `mapstructure:"headers" validate:"min=1"`
}

type HTTPSettings struct {
	Port   string       `mapstructure:"port" validate:"required,numeric"`
	Prefix string       `mapstructure:"prefix" validate:"required"`
	IP     string       `mapstructure:"ip" validate:"required,ip"`
	CORS   CORSSettings `mapstructure:"cors" validate:"required"`
}

type ObservabilitySettings struct {
	Enabled  bool   `mapstructure:"enabled"`
	Endpoint string `mapstructure:"endpoint" validate:"required_if=Enabled true,url"`
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
	BucketName string `mapstructure:"bucket_name" validate:"required"`
}

type DatabaseSettings struct {
	Host                   string `mapstructure:"host" validate:"required"`
	Port                   int    `mapstructure:"port" validate:"required,gte=1,lte=65535"`
	User                   string `mapstructure:"user" validate:"required"`
	Password               string `mapstructure:"password" validate:"required"`
	Database               string `mapstructure:"database" validate:"required"`
	Schema                 string `mapstructure:"schema"`
	SSLMode                string `mapstructure:"ssl_mode" validate:"oneof=disable require verify-ca verify-full"`
	MaxOpenConns           int    `mapstructure:"max_open_conns" validate:"gte=1"`
	MaxIdleConns           int    `mapstructure:"max_idle_conns" validate:"gte=1"`
	ConnMaxLifetimeMinutes int    `mapstructure:"conn_max_lifetime_minutes" validate:"gte=1"`
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
