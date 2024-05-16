package service

import (
	"bytes"
	"main/src/books/domain/repository"
	appError "main/utils/error"
)

type BookFileServiceS3 struct {
	repo repository.BookFileRepository
}

func NewBookFileServiceS3(repo repository.BookFileRepository) BookFileService {
	return &BookFileServiceS3{
		repo: repo,
	}
}

func (service *BookFileServiceS3) SaveBookFile(file *bytes.Reader, bucketKey, fileExt string) *appError.Error {
	return service.repo.SaveBookFile(file, bucketKey, fileExt)
}

func (service *BookFileServiceS3) DeleteBookFile(bucketKey string) *appError.Error {
	return service.repo.DeleteBookFile(bucketKey)
}