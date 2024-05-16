package model_test

import (
	"main/src/books/domain/model"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type BookModelSuite struct {
	suite.Suite
}

func (s *BookModelSuite) SetupSuite()    {}
func (s *BookModelSuite) TearDownSuite() {}

func (s *BookModelSuite) TestValidate() {
	var tests = []struct {
		book     model.Book
		expected bool // true if no error is expected, false otherwise
	}{
		{model.Book{ID: "123e4567-e89b-12d3-a456-426614174000", Name: "The Great Gatsby", Description: "A classic novel.", ImgURL: "https://example.com/image.jpg"}, true},
		{model.Book{ID: "invalid-uuid", Name: "The Great Gatsby", Description: "A classic novel.", ImgURL: "https://example.com/image.jpg"}, false},
		{model.Book{ID: "123e4567-e89b-12d3-a456-426614174000", Name: "", Description: "A classic novel.", ImgURL: "https://example.com/image.jpg"}, false},
		{model.Book{ID: "123e4567-e89b-12d3-a456-426614174000", Name: "The Great Gatsby", Description: strings.Repeat("A", 501), ImgURL: "https://example.com/image.jpg"}, false},
		{model.Book{ID: "123e4567-e89b-12d3-a456-426614174000", Name: "The Great Gatsby", Description: "A classic novel.", ImgURL: "ftp://example.com/image.jpg"}, false},
	}

	for _, tt := range tests {
		s.Run(tt.book.Name+"_"+tt.book.ID, func() {
			err := tt.book.Validate()
			if tt.expected {
				s.Nil(err)
			} else {
				s.NotNil(err)
			}
		})
	}
}

func TestBookModelSuite(t *testing.T) {
	suite.Run(t, new(BookModelSuite))
}
