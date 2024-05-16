package service

import (
	"sync"
	"github.com/google/uuid"
	"main/src/books/domain/model"
	"main/src/books/domain/repository"
	appError "main/utils/error"
	"main/utils/lib"
)

type BookServiceDynamoDB struct {
	repo repository.BookRepository
}

func NewBookServiceDynamoDB(repo repository.BookRepository) BookService {
	return &BookServiceDynamoDB{
		repo: repo,
	}
}

func (service *BookServiceDynamoDB) GetAllBooks() ([]model.Book, *appError.Error) {
	return service.repo.GetAllBooks()
}

func (service *BookServiceDynamoDB) CreateBook(book *model.Book) (*model.Book, *appError.Error) {
	if book.ID == "" {
		book.ID = uuid.NewString()
	}
	if err := book.Validate(); err != nil {
		return nil, err
	}
	return service.repo.CreateBook(book)
}

func (service *BookServiceDynamoDB) CreateBatchBooks(books []model.Book) *appError.Error {
	var wg sync.WaitGroup
	errorChan := make(chan *appError.Error, len(books))
	for i, book := range books {
		wg.Add(1)
		if book.ID == "" {
			book.ID = uuid.NewString()
		}
		books[i] = book
		go func(b model.Book) {
			defer wg.Done()
			if err := b.Validate(); err != nil {
				errorChan <- err
				return
			}
			errorChan <- nil
		}(book)
	}
	wg.Wait()
	close(errorChan)
	for err := range errorChan {
		if err != nil {
			return err
		}
	}
	return service.repo.CreateBatchBooks(books)
}

func (service *BookServiceDynamoDB) GetBookByID(bookID string) (*model.Book, *appError.Error) {
	if err := lib.ValidateUUID(bookID); err != nil {
		return nil, err
	}
	return service.repo.GetBookByID(bookID)
}

func (service *BookServiceDynamoDB) UpdateBookByID(bookID string, book *model.Book) (*model.Book, *appError.Error) {
	if err := lib.ValidateUUID(bookID); err != nil {
		return nil, err
	}
	if err := book.Validate(); err != nil {
		return nil, err
	}
	return service.repo.UpdateBookByID(bookID, book)
}

func (service *BookServiceDynamoDB) DeleteBookByID(bookID string) *appError.Error {
	if err := lib.ValidateUUID(bookID); err != nil {
		return err
	}
	return service.repo.DeleteBookByID(bookID)
}
