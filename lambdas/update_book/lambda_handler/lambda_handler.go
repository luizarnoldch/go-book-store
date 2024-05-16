package lambdahandler

import (
	"bytes"
	"context"
	"encoding/base64"
	"log"
	"main/utils/apigateway"
	"main/utils/lib"
	"net/http"
	"os"
	"path/filepath"

	book "main/src/books/application/handler"
	"main/src/books/domain/model"

	"github.com/aws/aws-lambda-go/events"
	"github.com/mitchellh/mapstructure"
)

var (
	BOOKS_TABLE = os.Getenv("BOOKS_TABLE")
	BUCKET_NAME = os.Getenv("BUCKET_NAME")
	BUCKET_KEY  = os.Getenv("BUCKET_KEY")
)

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	bookMicro := book.MicroAWSBookDynamoDB{
		Ctx:        ctx,
		TableName:  BOOKS_TABLE,
		BucketName: BUCKET_NAME,
		BucketKey:  BUCKET_KEY,
	}

	bookId, errApi := apigateway.ParseAPIGatewayRequestParameters(request, "bookId")
	if errApi != nil {
		log.Printf("Error parsing request parameters: %v", errApi)
		return apigateway.APIGatewayError(http.StatusBadRequest, "Error parsing request parameters.")
	}
	
	decodedBody, err := base64.StdEncoding.DecodeString(request.Body)
	if err != nil {
		log.Printf("Error decoding base64 body: %v", err)
		return apigateway.APIGatewayError(http.StatusBadRequest, "Error decoding base64 body.")
	}

	boundary, errMultipart := lib.GetBoundaryFromMultipart(request.Headers["Content-Type"])
	if errMultipart != nil {
		log.Printf("Error getting boundary from Mutlipart: %v", errMultipart.ToString())
		return apigateway.APIGatewayError(errMultipart.Code, errMultipart.ToString())
	}

	reader := lib.GetFileReader(decodedBody, boundary)
	fileName, fileContent, formData, errForm := lib.GetFormDataFromDecodedBody(reader)
	if errForm != nil {
		log.Printf("Error getting data from fileReader: %v", errForm.ToString())
		return apigateway.APIGatewayError(errForm.Code, errForm.ToString())
	}

	fileExt := filepath.Ext(fileName)
	customKey := BUCKET_KEY + bookId + fileExt
	imgURL := "https://" + BUCKET_NAME + ".s3.amazonaws.com/" + customKey

	book := model.Book{
		ID:     bookId,
		ImgURL: imgURL,
	}
	errMap := mapstructure.Decode(formData, &book)
	if errMap != nil {
		log.Println("Error decoding formName to Name:", err)
		return apigateway.APIGatewayError(http.StatusInternalServerError, "Error decoding formName to Name.")
	}

	bookFile := bytes.NewReader(fileContent.Bytes())
	errBookMicro := bookMicro.SaveBookFile(bookFile, customKey, fileExt)
	if errBookMicro != nil {
		log.Printf("Error while saving book file, %s", errBookMicro.ToString())
		return apigateway.APIGatewayError(errBookMicro.Code, errBookMicro.ToString())
	}

	newBook, errBookMicro := bookMicro.UpdateBookByID(bookId, &book)
	if errBookMicro != nil {
		log.Printf("Error while creating book, %s", errBookMicro.ToString())
		return apigateway.APIGatewayError(errBookMicro.Code, errBookMicro.ToString())
	}

	return apigateway.APIGatewayDataResponse(http.StatusOK, newBook)
}
