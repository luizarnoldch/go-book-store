package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	index "main/lambdas/create_book/lambda_handler"
)

func main() {
	lambda.Start(index.Handler)
}