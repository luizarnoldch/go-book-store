package configuration

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)


func GetDynamoDBBookTable() string {
	tableName := os.Getenv("BOOKS_TABLE")
	if tableName == "" {
		log.Printf("Local DynamoDB Database")
		return "Test_Book_Table"
	}
	log.Printf("AWS DynamoDB Database: %s", tableName)
	return tableName
}

func GetDynamoDBClient(ctx context.Context) (*dynamodb.Client, error) {
	tableName := os.Getenv("BOOKS_TABLE")
	if tableName == "" {
		return GetLocalDynamoDBClient(ctx)
	}
	return GetAWSDynamoDBClient(ctx)
}

func CreateLocalDynamoDBBookTable(ctx context.Context, client *dynamodb.Client, tableName string) error {
	_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("book_id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("book_id"),
				KeyType:       types.KeyTypeHash,
			},
			// {
			// 	AttributeName: aws.String("date"),
			// 	KeyType:       types.KeyTypeRange,
			// },
		},
		TableName:   aws.String(tableName),
		BillingMode: types.BillingModePayPerRequest,
	})

	if err != nil {
		log.Printf("Error creating table %s: %s", tableName, err)
		return err
	}

	log.Printf("Table %s created successfully", tableName)
	return nil
}

func DescribeBookTable(ctx context.Context, client *dynamodb.Client, tableName string) (bool, error) {
	_, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})

	if err != nil {
		var notFoundErr *types.ResourceNotFoundException
		if errors.As(err, &notFoundErr) {
			log.Printf("Table %s does not exist.", tableName)
			return false, nil
		}
		log.Printf("Unexpected error occurred while describing table %s: %s", tableName, err)
		return false, err
	}

	log.Printf("Table %s exists.", tableName)
	return true, nil
}

func DeleteLocalDynamoDBBookTable(ctx context.Context, client *dynamodb.Client, tableName string) error {
	_, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})

	if err != nil {
		return fmt.Errorf("table %s does not exist, no need to delete", tableName)
	}
	client.DeleteTable(ctx, &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	})

	log.Printf("Table %s deleted successfully", tableName)
	return nil
}

func GetAWSDynamoDBClient(ctx context.Context) (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("Error setting AWS configuration: %v", err)
		return nil, err
	}
	log.Printf("AWS Client connected successfully")
	return dynamodb.NewFromConfig(cfg), nil
}

func GetAWSS3Client(ctx context.Context) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("Error setting AWS configuration: %v", err)
		return nil, err
	}
	log.Printf("AWS Client connected successfully")
	return s3.NewFromConfig(cfg), nil
}

func GetLocalEndpoint(service, region string, options ...interface{}) (aws.Endpoint, error) {
	return aws.Endpoint{URL: "http://localhost:8000"}, nil
}

func GetLocalDynamoDBClient(ctx context.Context) (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(
				GetLocalEndpoint,
			),
		),
	)
	if err != nil {
		log.Printf("Error setting Local configurationb: %v", err)
		return nil, err
	}
	log.Printf("Local Client connected successfully")
	return dynamodb.NewFromConfig(cfg), nil
}
