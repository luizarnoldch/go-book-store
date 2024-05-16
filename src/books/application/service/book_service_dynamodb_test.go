package service_test

import (
	"testing"

	"main/src/books/application/service"
	"main/src/books/domain/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	repoMock "main/mocks"
)

type BookServiceDynamoDBSuite struct {
	suite.Suite
	bookRepository *repoMock.BookRepository
	bookService    service.BookService
	testBook       *model.Book
	uuidGlobal     string
}

const (
	MethodGetAllBooks      = "GetAllBooks"
	MethodCreateBook       = "CreateBook"
	MethodCreateBatchBooks = "CreateBatchBooks"
	MethodGetBookByID      = "GetBookByID"
	MethodUpdateBookByID   = "UpdateBookByID"
	MethodDeleteBookByID   = "DeleteBookByID"
)

func (suite *BookServiceDynamoDBSuite) SetupTest() {
	suite.bookRepository = new(repoMock.BookRepository)
	suite.bookService = service.NewBookServiceDynamoDB(suite.bookRepository)
	suite.uuidGlobal = uuid.NewString()
	suite.testBook = &model.Book{
		ID:          suite.uuidGlobal,
		Name:        "Test Book",
		Description: "A descriptive text for testing",
		ImgURL:      "https://example.com/book.jpg",
	}
}

func (suite *BookServiceDynamoDBSuite) TestCreateBook() {
	suite.bookRepository.On(MethodCreateBook, suite.testBook).Return(suite.testBook, nil)
	createdBook, err := suite.bookService.CreateBook(suite.testBook)
	suite.Nil(err)
	suite.NotEmpty(createdBook.ID)
	suite.Equal(createdBook.Name, suite.testBook.Name)
	suite.uuidGlobal = createdBook.ID
	suite.bookRepository.AssertExpectations(suite.T())
}

func (suite *BookServiceDynamoDBSuite) TestGetAllBooks() {
	suite.bookRepository.On(MethodGetAllBooks).Return([]model.Book{*suite.testBook}, nil)
	books, err := suite.bookService.GetAllBooks()
	suite.Nil(err)
	suite.Len(books, 1)
	suite.Equal("Test Book", books[0].Name)
	suite.bookRepository.AssertExpectations(suite.T())
}

func (suite *BookServiceDynamoDBSuite) TestGetBookByID() {
	suite.bookRepository.On(MethodGetBookByID, suite.uuidGlobal).Return(suite.testBook, nil)
	book, err := suite.bookService.GetBookByID(suite.uuidGlobal)
	suite.Nil(err)
	suite.Equal("Test Book", book.Name)
	suite.bookRepository.AssertExpectations(suite.T())
}

func (suite *BookServiceDynamoDBSuite) TestUpdateBookByID() {
	updatedBook := &model.Book{
		ID:          suite.uuidGlobal,
		Name:        "Updated Test Book",
		Description: "A book used for testing with updated content",
		ImgURL:      "https://example.com/updated.jpg",
	}
	suite.bookRepository.On(MethodUpdateBookByID, suite.uuidGlobal, updatedBook).Return(updatedBook, nil)
	book, err := suite.bookService.UpdateBookByID(suite.uuidGlobal, updatedBook)
	suite.Nil(err)
	suite.NotNil(book, "Book should not be nil")
	suite.Equal("Updated Test Book", book.Name, "Book name should be updated")
	suite.bookRepository.AssertExpectations(suite.T())
}

func (suite *BookServiceDynamoDBSuite) TestDeleteBookByID() {
	suite.bookRepository.On(MethodDeleteBookByID, suite.uuidGlobal).Return(nil)
	err := suite.bookService.DeleteBookByID(suite.uuidGlobal)
	suite.Nil(err)
	suite.bookRepository.AssertExpectations(suite.T())
}

func TestBookServiceDynamoDBSuite(t *testing.T) {
	suite.Run(t, new(BookServiceDynamoDBSuite))
}
