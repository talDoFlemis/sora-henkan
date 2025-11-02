package telemetry

import (
	"context"
	"log/slog"
	"math/rand"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	slogmulti "github.com/samber/slog-multi"
	"github.com/taldoflemis/sora-henkan/settings"
)

func errorFormattingMiddleware(
	ctx context.Context,
	record slog.Record,
	next func(context.Context, slog.Record) error,
) error {
	attrs := []slog.Attr{}

	record.Attrs(func(attr slog.Attr) bool {
		key := attr.Key
		value := attr.Value
		kind := attr.Value.Kind()

		if (key == "error" || key == "err") && kind == slog.KindAny {
			if err, ok := value.Any().(error); ok {
				errType := reflect.TypeOf(err).String()
				msg := err.Error()

				attrs = append(
					attrs,
					slog.Group("error",
						slog.String("type", errType),
						slog.String("message", msg),
					),
				)
			} else {
				attrs = append(attrs, attr)
			}
		} else {
			attrs = append(attrs, attr)
		}

		return true
	})

	record = slog.NewRecord(record.Time, record.Level, record.Message, record.PC)
	record.AddAttrs(attrs...)

	return next(ctx, record)
}

func DynamoDBSlogHandler(client *dynamodb.Client, dynamoDBSettings settings.DynamoDBLogsSettings) (slog.Handler, error) {
	mdw := slogmulti.NewHandleInlineHandler(
		func(ctx context.Context, groups []string, attrs []slog.Attr, record slog.Record) error {
			sample := rand.Float64()
			if sample >= 0.2 {
				return nil
			}

			// Sample
			item := make(map[string]types.AttributeValue, 0)

			for _, attr := range attrs {
				item[attr.Key] = &types.AttributeValueMemberS{Value: attr.Value.String()}
			}

			item["id"] = &types.AttributeValueMemberS{Value: uuid.New().String()}
			item["timestamp"] = &types.AttributeValueMemberS{Value: record.Time.String()}
			item["message"] = &types.AttributeValueMemberS{Value: record.Message}
			item["level"] = &types.AttributeValueMemberS{Value: record.Level.String()}

			go func() {
				_, err := client.PutItem(context.Background(), &dynamodb.PutItemInput{
					TableName: aws.String(dynamoDBSettings.Table),
					Item:      item,
				},
				)
				if err != nil {
					slog.ErrorContext(context.Background(), "failed to send log to dynamodb", slog.Any("err", err))
				}
			}()
			return nil
		},
	)

	return mdw, nil
}
