package adapter

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"main/src/books/domain/model"
	appError "main/utils/error"
)

type BookDynamoDBRepository struct {
	ctx    context.Context
	client *dynamodb.Client
	table  string
}

func NewBookDynamoDBRepository(ctx context.Context, client *dynamodb.Client, table string) *BookDynamoDBRepository {
	return &BookDynamoDBRepository{
		ctx:    ctx,
		client: client,
		table:  table,
	}
}

func (r *BookDynamoDBRepository) GetAllBooks() ([]model.Book, *appError.Error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(r.table),
	}
	result, err := r.client.Scan(r.ctx, input)
	if err != nil {
		log.Printf("Error scanning DynamoDB table: %v, table: %s", err, r.table)
		return nil, appError.NewUnexpectedError(err.Error())
	}

	var books []model.Book
	for _, item := range result.Items {
		var book model.Book
		err := attributevalue.UnmarshalMap(item, &book)
		if err != nil {
			log.Printf("Error unmarshaling item from DynamoDB: %v, item: %+v", err, item)
			return nil, appError.NewUnexpectedError(err.Error())
		}
		books = append(books, book)
	}
	log.Println("Retrieved all books successfully")
	return books, nil
}

func (r *BookDynamoDBRepository) CreateBook(book *model.Book) (*model.Book, *appError.Error) {
	av, err := attributevalue.MarshalMap(book)
	if err != nil {
		log.Printf("Error marshaling book: %v, book: %+v", err, book)
		return &model.Book{}, appError.NewUnexpectedError(err.Error())
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(r.table),
	}
	_, err = r.client.PutItem(r.ctx, input)
	if err != nil {
		log.Printf("Error putting item in DynamoDB: %v, table: %s", err, r.table)
		return &model.Book{}, appError.NewUnexpectedError(err.Error())
	}

	log.Printf("Book creation completed successfully, book: %+v", book)
	return book, nil
}

func (r *BookDynamoDBRepository) CreateBatchBooks(books []model.Book) *appError.Error {
	var writeRequests []types.WriteRequest
	for _, book := range books {
		av, err := attributevalue.MarshalMap(book)
		if err != nil {
			log.Printf("Error while marshalling book: %s, book: %+v", err, book)
			return appError.NewUnexpectedError(err.Error())
		}
		writeRequests = append(writeRequests, types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: av,
			},
		})
	}

	const maxBatchSize = 25
	for i := 0; i < len(writeRequests); i += maxBatchSize {
		end := i + maxBatchSize
		if end > len(writeRequests) {
			end = len(writeRequests)
		}
		batch := writeRequests[i:end]

		input := &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				r.table: batch,
			},
		}

		_, err := r.client.BatchWriteItem(r.ctx, input)
		if err != nil {
			log.Printf("Error while batch writing items: %s, batch size: %d", err, len(batch))
			return appError.NewUnexpectedError(err.Error())
		}
	}
	log.Println("Batch books creation completed successfully")
	return nil
}

func (r *BookDynamoDBRepository) GetBookByID(id string) (*model.Book, *appError.Error) {
	key := map[string]types.AttributeValue{
		"ID": &types.AttributeValueMemberS{Value: id},
	}
	
	input := &dynamodb.GetItemInput{
		Key:       key,
		TableName: aws.String(r.table),
	}
	result, err := r.client.GetItem(r.ctx, input)
	if err != nil {
		log.Printf("Error getting item from DynamoDB: %v, table: %s", err, r.table)
		return &model.Book{}, appError.NewUnexpectedError(err.Error())
	}

	var book model.Book
	if result.Item == nil {
		log.Println("No book found with ID:", id)
		return &model.Book{}, nil // No error but no data
	}

	err = attributevalue.UnmarshalMap(result.Item, &book)
	if err != nil {
		log.Printf("Error unmarshaling item from DynamoDB: %v, item: %+v", err, result.Item)
		return &model.Book{}, appError.NewUnexpectedError(err.Error())
	}

	log.Printf("Retrieved book successfully, ID: %s, book: %+v", id, book)
	return &book, nil
}

func (r *BookDynamoDBRepository) UpdateBookByID(id string, book *model.Book) (*model.Book, *appError.Error) {
	keyCond := map[string]types.AttributeValue{
		"ID": &types.AttributeValueMemberS{Value: id},
	}
	
	update := expression.Set(
		expression.Name("name"), expression.Value(book.Name),
	).Set(
		expression.Name("description"), expression.Value(book.Description),
	).Set(
		expression.Name("img_url"), expression.Value(book.ImgURL),
	)

	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		log.Printf("Error building expression for update: %v, ID: %s", err, id)
		return &model.Book{}, appError.NewUnexpectedError(err.Error())
	}

	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(r.table),
		Key:                       keyCond,
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ReturnValues:              types.ReturnValueUpdatedNew,
	}

	result, err := r.client.UpdateItem(r.ctx, input)
	if err != nil {
		log.Printf("Error updating item in DynamoDB: %v, table: %s", err, r.table)
		return &model.Book{}, appError.NewUnexpectedError(err.Error())
	}

	var updatedBook model.Book
	err = attributevalue.UnmarshalMap(result.Attributes, &updatedBook)
	if err != nil {
		log.Printf("Error unmarshaling updated item from DynamoDB: %v, item: %+v", err, result.Attributes)
		return &model.Book{}, appError.NewUnexpectedError(err.Error())
	}

	log.Printf("Updated book successfully, ID: %s, book: %+v", id, updatedBook)
	return &updatedBook, nil
}

func (r *BookDynamoDBRepository) DeleteBookByID(id string) *appError.Error {
	key := map[string]types.AttributeValue{
		"ID": &types.AttributeValueMemberS{Value: id},
	}

	input := &dynamodb.DeleteItemInput{
		Key:       key,
		TableName: aws.String(r.table),
	}

	result, err := r.client.DeleteItem(r.ctx, input)
	if err != nil {
		log.Printf("Error deleting item from DynamoDB: %v, table: %s", err, r.table)
		return appError.NewUnexpectedError(err.Error())
	}

	var deletedBook model.Book
	if result.Attributes != nil {
		err = attributevalue.UnmarshalMap(result.Attributes, &deletedBook)
		if err != nil {
			log.Printf("Error unmarshaling deleted item from DynamoDB: %v, item: %+v", err, result.Attributes)
			return appError.NewUnexpectedError(err.Error())
		}
	}

	log.Printf("Deleted book successfully, book_id: %s, book: %+v", id, deletedBook)
	return nil
}
