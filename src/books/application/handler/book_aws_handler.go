package handler

import (
	"bytes"
	"context"
	"log"

	"main/src/books/application/service"
	"main/src/books/domain/model"
	"main/src/books/infrastructure/adapter"
	"main/src/books/infrastructure/configuration"
	appError "main/utils/error"
)

type MicroAWSBookDynamoDB struct {
	Ctx        context.Context
	TableName  string
	BucketName string
	BucketKey  string
}

func (micro *MicroAWSBookDynamoDB) GetAllBooks() ([]model.Book, *appError.Error) {
	dynamoClient, err := configuration.GetDynamoDBClient(micro.Ctx)
	if err != nil {
		log.Println("Error while defining local/AWS database")
		return nil, appError.NewUnexpectedError(err.Error())
	}
	if micro.TableName == "" {
		micro.TableName = configuration.GetDynamoDBBookTable()
	}
	bookInfrastructure := adapter.NewBookDynamoDBRepository(micro.Ctx, dynamoClient, micro.TableName)
	bookService := service.NewBookServiceDynamoDB(bookInfrastructure)

	return bookService.GetAllBooks()
}

func (micro *MicroAWSBookDynamoDB) CreateBook(book *model.Book) (*model.Book, *appError.Error) {
	dynamoClient, err := configuration.GetDynamoDBClient(micro.Ctx)
	if err != nil {
		log.Println("Error while defining local/AWS database")
		return nil, appError.NewUnexpectedError(err.Error())
	}
	if micro.TableName == "" {
		micro.TableName = configuration.GetDynamoDBBookTable()
	}
	bookInfrastructure := adapter.NewBookDynamoDBRepository(micro.Ctx, dynamoClient, micro.TableName)
	bookService := service.NewBookServiceDynamoDB(bookInfrastructure)

	return bookService.CreateBook(book)
}

func (micro *MicroAWSBookDynamoDB) CreateBatchBooks(books []model.Book) *appError.Error {
	dynamoClient, err := configuration.GetDynamoDBClient(micro.Ctx)
	if err != nil {
		log.Println("Error while defining local/AWS database")
		return appError.NewUnexpectedError(err.Error())
	}
	if micro.TableName == "" {
		micro.TableName = configuration.GetDynamoDBBookTable()
	}
	bookInfrastructure := adapter.NewBookDynamoDBRepository(micro.Ctx, dynamoClient, micro.TableName)
	bookService := service.NewBookServiceDynamoDB(bookInfrastructure)

	return bookService.CreateBatchBooks(books)
}

func (micro *MicroAWSBookDynamoDB) GetBookByID(bookID string) (*model.Book, *appError.Error) {
	dynamoClient, err := configuration.GetDynamoDBClient(micro.Ctx)
	if err != nil {
		log.Println("Error while defining local/AWS database")
		return nil, appError.NewUnexpectedError(err.Error())
	}
	if micro.TableName == "" {
		micro.TableName = configuration.GetDynamoDBBookTable()
	}
	bookInfrastructure := adapter.NewBookDynamoDBRepository(micro.Ctx, dynamoClient, micro.TableName)
	bookService := service.NewBookServiceDynamoDB(bookInfrastructure)

	return bookService.GetBookByID(bookID)
}

func (micro *MicroAWSBookDynamoDB) UpdateBookByID(bookID string, book *model.Book) (*model.Book, *appError.Error) {
	dynamoClient, err := configuration.GetDynamoDBClient(micro.Ctx)
	if err != nil {
		log.Println("Error while defining local/AWS database")
		return nil, appError.NewUnexpectedError(err.Error())
	}
	if micro.TableName == "" {
		micro.TableName = configuration.GetDynamoDBBookTable()
	}
	bookInfrastructure := adapter.NewBookDynamoDBRepository(micro.Ctx, dynamoClient, micro.TableName)
	bookService := service.NewBookServiceDynamoDB(bookInfrastructure)

	return bookService.UpdateBookByID(bookID, book)
}

func (micro *MicroAWSBookDynamoDB) DeleteBookByID(bookID string) *appError.Error {
	dynamoClient, err := configuration.GetDynamoDBClient(micro.Ctx)
	if err != nil {
		log.Println("Error while defining local/AWS database")
		return appError.NewUnexpectedError(err.Error())
	}
	if micro.TableName == "" {
		micro.TableName = configuration.GetDynamoDBBookTable()
	}
	bookInfrastructure := adapter.NewBookDynamoDBRepository(micro.Ctx, dynamoClient, micro.TableName)
	bookService := service.NewBookServiceDynamoDB(bookInfrastructure)

	return bookService.DeleteBookByID(bookID)
}

func (micro *MicroAWSBookDynamoDB) SaveBookFile(file *bytes.Reader, bucketKey, fileExt string) *appError.Error {
	s3Client, err := configuration.GetAWSS3Client(micro.Ctx)
	if err != nil {
		log.Println("Error while defining local/AWS database")
		return appError.NewUnexpectedError(err.Error())
	}

	bookInfrastructure := adapter.NewBookFileRepositoryS3(micro.Ctx, s3Client, micro.BucketName, micro.BucketKey)
	bookService := service.NewBookFileServiceS3(bookInfrastructure)

	return bookService.SaveBookFile(file, bucketKey, fileExt)
}

func (micro *MicroAWSBookDynamoDB) DeleteBookFile(bucketKey string) *appError.Error {
	s3Client, err := configuration.GetAWSS3Client(micro.Ctx)
	if err != nil {
		log.Println("Error while defining local/AWS database")
		return appError.NewUnexpectedError(err.Error())
	}

	bookInfrastructure := adapter.NewBookFileRepositoryS3(micro.Ctx, s3Client, micro.BucketName, micro.BucketKey)
	bookService := service.NewBookFileServiceS3(bookInfrastructure)

	return bookService.DeleteBookFile(bucketKey)
}
