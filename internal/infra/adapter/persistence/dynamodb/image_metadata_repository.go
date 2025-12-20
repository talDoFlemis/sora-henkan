package dynamodb

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/taldoflemis/sora-henkan/internal/core/domain/images"
	"github.com/taldoflemis/sora-henkan/internal/core/ports"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var tracer = otel.Tracer("")

// metadataModel represents the DynamoDB persistence model for image metadata
type metadataModel struct {
	ID                    string `dynamodbav:"id"`
	OriginalImageURL      string `dynamodbav:"original_image_url"`
	ObjectStorageImageKey string `dynamodbav:"object_storage_image_key"`
	TransformedImageKey   string `dynamodbav:"transformed_image_key"`
	MimeType              string `dynamodbav:"mime_type"`
	Status                string `dynamodbav:"status"`
	Checksum              string `dynamodbav:"checksum"`
	ErrorMessage          string `dynamodbav:"error_message"`
	TransformationCount   int    `dynamodbav:"transformation_count"`
	UpdatedAt             string `dynamodbav:"updated_at"`
	CreatedAt             string `dynamodbav:"created_at"`
}

// toDomain converts a persistence model to domain model
func (m *metadataModel) toDomain() (*images.ImageMetadata, error) {
	updatedAt, err := time.Parse(time.RFC3339, m.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse updated_at: %w", err)
	}

	createdAt, err := time.Parse(time.RFC3339, m.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", err)
	}

	return &images.ImageMetadata{
		ID:                    m.ID,
		OriginalImageURL:      m.OriginalImageURL,
		ObjectStorageImageKey: m.ObjectStorageImageKey,
		TransformedImageKey:   m.TransformedImageKey,
		MimeType:              m.MimeType,
		Status:                m.Status,
		Checksum:              m.Checksum,
		ErrorMessage:          m.ErrorMessage,
		TransformationCount:   m.TransformationCount,
		UpdatedAt:             updatedAt,
		CreatedAt:             createdAt,
	}, nil
}

// fromDomain converts a domain model to persistence model
func fromMetadataDomain(meta *images.ImageMetadata) *metadataModel {
	return &metadataModel{
		ID:                    meta.ID,
		OriginalImageURL:      meta.OriginalImageURL,
		ObjectStorageImageKey: meta.ObjectStorageImageKey,
		TransformedImageKey:   meta.TransformedImageKey,
		MimeType:              meta.MimeType,
		Status:                meta.Status,
		Checksum:              meta.Checksum,
		ErrorMessage:          meta.ErrorMessage,
		TransformationCount:   meta.TransformationCount,
		UpdatedAt:             meta.UpdatedAt.Format(time.RFC3339),
		CreatedAt:             meta.CreatedAt.Format(time.RFC3339),
	}
}

type DynamoDBImageMetadataRepository struct {
	client    *dynamodb.Client
	tableName string
}

var _ ports.ImageMetadataRepository = (*DynamoDBImageMetadataRepository)(nil)

// NewDynamoDBImageMetadataRepository creates a new DynamoDB image metadata repository
func NewDynamoDBImageMetadataRepository(client *dynamodb.Client, tableName string) *DynamoDBImageMetadataRepository {
	return &DynamoDBImageMetadataRepository{
		client:    client,
		tableName: tableName,
	}
}

// SaveMetadata implements ports.ImageMetadataRepository.
func (d *DynamoDBImageMetadataRepository) SaveMetadata(ctx context.Context, image *images.Image) error {
	ctx, span := tracer.Start(ctx, "DynamoDBImageMetadataRepository.SaveMetadata")
	defer span.End()

	metadata := image.ToMetadata()
	model := fromMetadataDomain(metadata)

	item, err := attributevalue.MarshalMap(model)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	_, err = d.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(d.tableName),
		Item:      item,
	})
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to save metadata: %w", err)
	}

	span.SetAttributes(attribute.String("image.id", image.ID.String()))
	return nil
}

// UpdateMetadata implements ports.ImageMetadataRepository.
func (d *DynamoDBImageMetadataRepository) UpdateMetadata(ctx context.Context, image *images.Image) error {
	ctx, span := tracer.Start(ctx, "DynamoDBImageMetadataRepository.UpdateMetadata")
	defer span.End()

	metadata := image.ToMetadata()
	updatedAt := time.Now().Format(time.RFC3339)

	_, err := d.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: metadata.ID},
		},
		UpdateExpression: aws.String(`SET 
			original_image_url = :original_url,
			object_storage_image_key = :storage_key,
			transformed_image_key = :transformed_key,
			mime_type = :mime,
			#status = :status,
			checksum = :checksum,
			error_message = :error_msg,
			transformation_count = :trans_count,
			updated_at = :updated_at`),
		ExpressionAttributeNames: map[string]string{
			"#status": "status", // status is a reserved word in DynamoDB
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":original_url":    &types.AttributeValueMemberS{Value: metadata.OriginalImageURL},
			":storage_key":     &types.AttributeValueMemberS{Value: metadata.ObjectStorageImageKey},
			":transformed_key": &types.AttributeValueMemberS{Value: metadata.TransformedImageKey},
			":mime":            &types.AttributeValueMemberS{Value: metadata.MimeType},
			":status":          &types.AttributeValueMemberS{Value: metadata.Status},
			":checksum":        &types.AttributeValueMemberS{Value: metadata.Checksum},
			":error_msg":       &types.AttributeValueMemberS{Value: metadata.ErrorMessage},
			":trans_count":     &types.AttributeValueMemberN{Value: strconv.Itoa(metadata.TransformationCount)},
			":updated_at":      &types.AttributeValueMemberS{Value: updatedAt},
		},
	})
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to update metadata: %w", err)
	}

	span.SetAttributes(attribute.String("image.id", image.ID.String()))
	return nil
}

// DeleteMetadata implements ports.ImageMetadataRepository.
func (d *DynamoDBImageMetadataRepository) DeleteMetadata(ctx context.Context, id string) error {
	ctx, span := tracer.Start(ctx, "DynamoDBImageMetadataRepository.DeleteMetadata")
	defer span.End()

	_, err := d.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to delete metadata: %w", err)
	}

	span.SetAttributes(attribute.String("image.id", id))
	return nil
}

// GetMetadata implements ports.ImageMetadataRepository.
func (d *DynamoDBImageMetadataRepository) GetMetadata(ctx context.Context, id string) (*images.ImageMetadata, error) {
	ctx, span := tracer.Start(ctx, "DynamoDBImageMetadataRepository.GetMetadata")
	defer span.End()

	result, err := d.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to get metadata: %w", err)
	}

	if result.Item == nil {
		return nil, ports.ErrMetadataNotFound
	}

	var model metadataModel
	err = attributevalue.UnmarshalMap(result.Item, &model)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	metadata, err := model.toDomain()
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to convert to domain: %w", err)
	}

	span.SetAttributes(attribute.String("image.id", id))
	return metadata, nil
}
