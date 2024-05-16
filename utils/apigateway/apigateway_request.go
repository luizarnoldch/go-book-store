package apigateway

import (
	"encoding/json"

	appError "main/utils/error"

	"github.com/aws/aws-lambda-go/events"
)

func ParseAPIGatewayRequestBody(request events.APIGatewayProxyRequest, entity interface{}) *appError.Error {
	err := json.Unmarshal([]byte(request.Body), entity)
	if err != nil {
		return appError.NewBadRequestError(err.Error())
	}
	return nil
}

func ParseAPIGatewayRequestParameters(request events.APIGatewayProxyRequest, parameter string) (string, *appError.Error) {
	param := request.PathParameters[parameter]
	if param == "" {
		return "", appError.NewBadRequestError("path parameter from APIGateway doesn't exits")
	}
	return param, nil
}
