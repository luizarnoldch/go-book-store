package repository

import (
	"bytes"
	appError "main/utils/error"
)

type BookFileRepository interface {
	DeleteBookFile(string) *appError.Error
	SaveBookFile(*bytes.Reader, string, string) *appError.Error
}
