package apigateway

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

var HeadersJSON = map[string]string{
	"Access-Control-Allow-Origin":  "*",
	"Access-Control-Allow-Methods": "DELETE,GET,HEAD,POST,PUT",
	"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
	"Content-Type": "application/json",
}

func apiGatewayResponse(statusCode int, body interface{}, headers map[string]string) (events.APIGatewayProxyResponse, error) {
	responseData, err := json.Marshal(body)
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		return events.APIGatewayProxyResponse{
			Body:       string(`{"error": "Error marshaling response"}`),
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
		}, nil
	}
	log.Printf("Response: %s", responseData)
	return events.APIGatewayProxyResponse{
		Body:       string(responseData),
		StatusCode: statusCode,
		Headers:    headers,
	}, nil
}

func APIGatewayMessageResponse(statusCode int, message string) (events.APIGatewayProxyResponse, error) {
	return apiGatewayResponse(statusCode, map[string]string{"message": message}, HeadersJSON)
}

func APIGatewayDataResponse(statusCode int, data interface{}) (events.APIGatewayProxyResponse, error) {
	return apiGatewayResponse(statusCode, data, HeadersJSON)
}

func APIGatewayError(statusCode int, err string) (events.APIGatewayProxyResponse, error) {
	return apiGatewayResponse(statusCode, map[string]string{"error": err}, HeadersJSON)
}
