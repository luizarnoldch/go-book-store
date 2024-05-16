package lib

import (
	"bytes"
	"io"
	"log"
	appError "main/utils/error"
	"mime/multipart"
)

func GetFileReader(decodedBody []byte, boundary string) *multipart.Reader {
	// Read the multipart message
	// []bytes -> decodedBody
	// Convert []bytes to Reader -> bytes.NewReader(decodedBody)
	// Get the bounday of bytes params["boundary"]
	return multipart.NewReader(bytes.NewReader(decodedBody), boundary)
}

func GetFormDataFromDecodedBody(fileReader *multipart.Reader) (string, bytes.Buffer, map[string]interface{}, *appError.Error) {
	// Read the binary information based on their fileReader(Reader and boundary)
	var fileName string
	var fileContent bytes.Buffer
	formData := make(map[string]interface{})

	for {
		part, err := fileReader.NextPart()
		if err == io.EOF {
			
			break
		}
		if err != nil {
			log.Printf("Error reading multipart part: %v", err)
			return "", bytes.Buffer{}, nil, appError.NewBadRequestError("Error reading multipart part")
		}

		if part.FormName() == "file" {
			fileName = part.FileName()

			log.Println("partFileNameInside ", part.FileName())
			if _, err := io.Copy(&fileContent, part); err != nil {
				log.Printf("Error reading file part: %v", err)
				return "", bytes.Buffer{}, nil, appError.NewBadRequestError("Error reading file part")
			}
		} else {
			data, err := io.ReadAll(part)
			if err != nil {
				log.Println("Error reading formName:", err)
				return "", bytes.Buffer{}, nil, appError.NewBadRequestError("Error reading formName")
			}
			formData[part.FormName()] = string(data)
		}

	}

	log.Println("fileName Inside: ", fileName)

	return fileName, fileContent, formData, nil
}
