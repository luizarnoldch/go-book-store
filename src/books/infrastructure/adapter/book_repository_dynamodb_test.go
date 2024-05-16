package adapter_test

import (
	"context"
	"testing"

	"main/src/books/domain/model"
	"main/src/books/domain/repository"
	"main/src/books/infrastructure/adapter"
	"main/src/books/infrastructure/configuration"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type BookDynamoDBSuite struct {
	suite.Suite
	tableName      string
	initBooks      []model.Book
	dynamoClient   *dynamodb.Client
	bookRepository repository.BookRepository
}

func (suite *BookDynamoDBSuite) SetupSuite() {
	ctx := context.TODO()
	client, err := configuration.GetLocalDynamoDBClient(ctx) // Assuming similar helper functions exist
	suite.NoError(err)
	suite.dynamoClient = client

	suite.tableName = "Test_Book_Table" // Name of the DynamoDB table for tests

	suite.bookRepository = adapter.NewBookDynamoDBRepository(ctx, client, suite.tableName)

	configuration.CreateLocalDynamoDBBookTable(ctx, suite.dynamoClient, suite.tableName)

	suite.initBooks = []model.Book{ // Initializing with some books
		{ID: uuid.NewString(), Name: "Book One", Description: "A first book", ImgURL: "url1"},
		{ID: uuid.NewString(), Name: "Book Two", Description: "A second book", ImgURL: "url2"},
	}

	suite.NotNil(suite.bookRepository.CreateBatchBooks(suite.initBooks))
}

func (suite *BookDynamoDBSuite) TearDownSuite() {
	for _, book := range suite.initBooks {
		suite.NotNil(suite.bookRepository.DeleteBookByID(book.ID))
	}
}

func (suite *BookDynamoDBSuite) TestGetAllBooks() {
	books, err := suite.bookRepository.GetAllBooks()
	suite.NotNil(err)
	suite.Len(books, len(suite.initBooks))
}

func (suite *BookDynamoDBSuite) TestCreateBook() {
	newBook := model.Book{ID: uuid.NewString(), Name: "Book Three", Description: "A third book", ImgURL: "url3"}
	createdBook, err := suite.bookRepository.CreateBook(&newBook)
	suite.NotNil(err)
	suite.Equal(newBook.Name, createdBook.Name)
}

func (suite *BookDynamoDBSuite) TestUpdateBookByID() {
	book_id := suite.initBooks[0].ID
	update := model.Book{Name: "Updated Book One", Description: "Updated description", ImgURL: "updated_url1"}
	updatedBook, err := suite.bookRepository.UpdateBookByID(book_id, &update)
	suite.NotNil(err)
	suite.Equal("Updated Book One", updatedBook.Name)
}

func (suite *BookDynamoDBSuite) TestDeleteBookByID() {
	err := suite.bookRepository.DeleteBookByID("1")
	suite.NotNil(err)
}

func (suite *BookDynamoDBSuite) TestGetBookByID() {
	book, err := suite.bookRepository.GetBookByID("1")
	suite.NotNil(err)
	suite.NotNil(book)
}

func TestBookDynamoDBSuite(t *testing.T) {
	suite.Run(t, new(BookDynamoDBSuite))
}
