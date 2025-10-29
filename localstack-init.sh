#!/bin/bash

# Wait for LocalStack to be ready
echo "Waiting for LocalStack to be ready..."
sleep 5

# Create DynamoDB table with the same name as in Terraform
# Using simple configuration with just the hash key (id)
echo "Creating DynamoDB table..."
awslocal dynamodb create-table \
    --table-name sora-henkan-dev-main-table \
    --attribute-definitions \
        AttributeName=id,AttributeType=S \
    --key-schema \
        AttributeName=id,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --region us-east-1

echo "DynamoDB table created successfully!"

# List tables to verify
echo "Listing DynamoDB tables..."
awslocal dynamodb list-tables --region us-east-1

echo "LocalStack initialization complete!"
