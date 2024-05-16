package service

import (
	"bytes"
	appError "main/utils/error"
)

type BookFileService interface {
	DeleteBookFile(string) *appError.Error
	SaveBookFile(*bytes.Reader, string, string) *appError.Error
}
