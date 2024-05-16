package lambdahandler

import (
	"context"
	"log"
	"net/http"
	"os"

	book "main/src/books/application/handler"
	"main/utils/apigateway"

	"github.com/aws/aws-lambda-go/events"
)

var (
	BOOKS_TABLE = os.Getenv("BOOKS_TABLE")
)

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	bookMicro := book.MicroAWSBookDynamoDB{
		Ctx:        ctx,
		TableName:  BOOKS_TABLE,
	}

	book_record, errBookMicro := bookMicro.GetAllBooks() 
	if errBookMicro != nil {
		log.Printf("Error while saving book file, %s", errBookMicro.ToString())
		return apigateway.APIGatewayError(errBookMicro.Code, errBookMicro.ToString())
	}

	return apigateway.APIGatewayDataResponse(http.StatusOK, book_record)
}