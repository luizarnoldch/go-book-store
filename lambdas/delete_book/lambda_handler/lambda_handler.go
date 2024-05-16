package lambdahandler

import (
	"context"
	"log"
	"main/utils/apigateway"
	"net/http"
	"os"
	"strings"

	book "main/src/books/application/handler"

	"github.com/aws/aws-lambda-go/events"
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

	book_record, errBookMicro := bookMicro.GetBookByID(bookId) 
	if errBookMicro != nil {
		log.Printf("Error while saving book file, %s", errBookMicro.ToString())
		return apigateway.APIGatewayError(errBookMicro.Code, errBookMicro.ToString())
	}

	url_parts := strings.Split(book_record.ImgURL, BUCKET_KEY)
	fileName := url_parts[1]
	customKey := BUCKET_KEY + fileName

	errBookMicro = bookMicro.DeleteBookFile(customKey)
	if errBookMicro != nil {
		log.Printf("Error while saving book file, %s", errBookMicro.ToString())
		return apigateway.APIGatewayError(errBookMicro.Code, errBookMicro.ToString())
	}

	errBookMicro = bookMicro.DeleteBookByID(bookId)
	if errBookMicro != nil {
		log.Printf("Error while creating book, %s", errBookMicro.ToString())
		return apigateway.APIGatewayError(errBookMicro.Code, errBookMicro.ToString())
	}

	message := "Book " + bookId + " deleted"
	return apigateway.APIGatewayMessageResponse(http.StatusOK, message)
}
