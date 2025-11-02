package http

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func CloudFlareExtractClientIPfunc(req *http.Request) string {
	cf := req.Header["CF-Connecting-IP"]
	if len(cf) > 0 {
		return cf[0]
	}

	return echo.ExtractIPFromXFFHeader()(req)
}

func SSETimeoutSkipper(c echo.Context) bool {
	// Skip timeout for SSE endpoints
	path := c.Path()

	if c.Request().Header.Get("Accept") == "text/event-stream" || strings.Contains(path, "/sse") {
		return true
	}

	return false
}

func DynamoDBAuditLogger(client *dynamodb.Client, tableName string) echo.MiddlewareFunc {
	return middleware.BodyDump(func(c echo.Context, requestBody, responseBody []byte) {
		item := make(map[string]types.AttributeValue, 0)

		item["id"] = &types.AttributeValueMemberS{Value: uuid.New().String()}
		item["timestamp"] = &types.AttributeValueMemberS{Value: time.Now().String()}
		item["message"] = &types.AttributeValueMemberS{Value: "New request"}
		item["request-body"] = &types.AttributeValueMemberS{Value: string(requestBody)}
		item["response-body"] = &types.AttributeValueMemberS{Value: string(responseBody)}
		item["method"] = &types.AttributeValueMemberS{Value: c.Request().Method}
		item["path"] = &types.AttributeValueMemberS{Value: c.Request().URL.Path}
		item["user-agent"] = &types.AttributeValueMemberS{Value: c.Request().UserAgent()}
		item["ip"] = &types.AttributeValueMemberS{Value: c.RealIP()}
		item["response-status"] = &types.AttributeValueMemberS{Value: http.StatusText(c.Response().Status)}

		go func() {
			_, err := client.PutItem(context.Background(), &dynamodb.PutItemInput{
				TableName: aws.String(tableName),
				Item:      item,
			},
			)
			if err != nil {
				slog.ErrorContext(context.Background(), "failed to send log to dynamodb", slog.Any("err", err))
			}
		}()
	})
}
