package service

import (
	"main/src/books/domain/model"
	appError "main/utils/error"

)

type BookService interface {
	GetAllBooks() ([]model.Book, *appError.Error)
	CreateBook(*model.Book) (*model.Book, *appError.Error)
	CreateBatchBooks([]model.Book) *appError.Error
	GetBookByID(string) (*model.Book, *appError.Error)
	UpdateBookByID(string, *model.Book) (*model.Book, *appError.Error)
	DeleteBookByID(string) *appError.Error
}
