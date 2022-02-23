package main

import (
	"bytes"
	"context"
	"encoding/json"
	_ "errors"
	_ "regexp"
	_ "strings"
	_ "time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// For debugging purposes, use the following logger.
// It's instance of log.Logger (https://golang.org/pkg/log/)
var logger = GetLogger()

func CreateResponseWithStatus(status int) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
	}
}

func GetAuthenticatedUserEmail(token string) (email string, ok bool) {
	dynamoDBClient := GetDynamoDBClient()

	tableName := "token-email-lookup"

	result, err := dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		/* ... */
	})

	if err != nil {
		logger.Println("DynamoDB Error!", err)
		return "", false
	}

	item := TokenLookupItem{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)

	// Validate the given token with one from the database
	// and return user email if the tokens match ...

	return "", false
}

func QueryUserNotes(email string) []UserNote {
	dynamoDBClient := GetDynamoDBClient()

	// User the following date format for "now"
	// dateNow := time.Now().Format(time.RFC3339)

	userNotes := []UserNote{}

	_, err := dynamoDBClient.Query(&dynamodb.QueryInput{
		TableName: aws.String("user-notes"),
		/* ... */
	})

	if err != nil {
		logger.Println("error getting notes", err.Error())
		return nil
	}

	return userNotes
}

func AuthenticateUser(headers map[string]string) (string, error) {
	// You can get Authentication header in the following way:
	// authenticationHeader := headers["Authentication"]

	// Validate the Authentication header and retrieve token
	token := ""

	email, _ := GetAuthenticatedUserEmail(token)

	return email, nil
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := request.Headers

	email, err := AuthenticateUser(headers)

	if err != nil {
		// Return appropriate responses on failed authentication
	}

	userNotes := QueryUserNotes(email)

	var buf bytes.Buffer

	body, err := json.Marshal(userNotes)

	if err != nil {
		return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 400}, nil
	}

	json.HTMLEscape(&buf, body)

	resp := events.APIGatewayProxyResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
