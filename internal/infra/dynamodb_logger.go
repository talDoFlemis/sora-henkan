package infra

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type DynamoDBLogger struct {
	client *dynamodb.Client
	table  string
}

func NewDynamoDBLogger(cfg aws.Config, table string) *DynamoDBLogger {
	return &DynamoDBLogger{
		client: dynamodb.NewFromConfig(cfg),
		table:  table,
	}
}

func (l *DynamoDBLogger) LogRequest(ctx context.Context, method, path, userAgent, ip string, statusCode int) {
	id := uuid.New().String()
	timestamp := time.Now().UTC().Format(time.RFC3339)

	item := map[string]types.AttributeValue{
		"id":          &types.AttributeValueMemberS{Value: id},
		"timestamp":   &types.AttributeValueMemberS{Value: timestamp},
		"method":      &types.AttributeValueMemberS{Value: method},
		"path":        &types.AttributeValueMemberS{Value: path},
		"status_code": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", statusCode)},
		"user_agent":  &types.AttributeValueMemberS{Value: userAgent},
		"ip":          &types.AttributeValueMemberS{Value: ip},
		"action":      &types.AttributeValueMemberS{Value: "api_request"},
	}

	slog.InfoContext(ctx, "DynamoDB Logger: Attempting to log request", "table", l.table, "id", id, "method", method, "path", path)
	_, err := l.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(l.table),
		Item:      item,
	})
	if err != nil {
		slog.ErrorContext(ctx, "DynamoDB Logger: Failed to log request", "table", l.table, "id", id, "error", err)
	} else {
		slog.InfoContext(ctx, "DynamoDB Logger: Successfully logged request", "table", l.table, "id", id)
	}
}

func (l *DynamoDBLogger) LogAction(ctx context.Context, action, details string) {
	id := uuid.New().String()
	timestamp := time.Now().UTC().Format(time.RFC3339)

	item := map[string]types.AttributeValue{
		"id":        &types.AttributeValueMemberS{Value: id},
		"timestamp": &types.AttributeValueMemberS{Value: timestamp},
		"action":    &types.AttributeValueMemberS{Value: action},
		"details":   &types.AttributeValueMemberS{Value: details},
	}

	slog.InfoContext(ctx, "DynamoDB Logger: Attempting to log action", "table", l.table, "id", id, "action", action)
	_, err := l.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(l.table),
		Item:      item,
	})
	if err != nil {
		slog.ErrorContext(ctx, "DynamoDB Logger: Failed to log action", "table", l.table, "id", id, "error", err)
	} else {
		slog.InfoContext(ctx, "DynamoDB Logger: Successfully logged action", "table", l.table, "id", id)
	}
}
