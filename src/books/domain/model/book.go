package model

import (
	appError "main/utils/error"
	"main/utils/lib"
	"strings"
)

type Book struct {
	ID          string `json:"ID,omitempty" dynamodbav:"ID,omitempty" mapstructure:"ID"`
	Name        string `json:"name,omitempty" dynamodbav:"name,omitempty" mapstructure:"name"`
	Description string `json:"description,omitempty" dynamodbav:"description,omitempty" mapstructure:"description"`
	ImgURL      string `json:"img_url,omitempty" dynamodbav:"img_url,omitempty" mapstructure:"img_url"`
}

func (b *Book) Validate() *appError.Error {
	if err := lib.ValidateUUID(b.ID); err != nil {
		return err
	}
	if err := lib.ValidateStringNotEmpty(b.Name); err != nil {
		return err
	}
	if err := lib.ValidateMaxStringCharacteres(b.Description, 200); err != nil {
		return err
	}
	if !strings.HasPrefix(b.ImgURL, "http://") && !strings.HasPrefix(b.ImgURL, "https://") {
		return appError.NewValidationError("Image URL must start with 'http://' or 'https://'.")
	}
	return nil
}
