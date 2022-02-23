package main

import (
	"bytes"
	"context"
	"encoding/json"
	_ "errors"
	"fmt"
	"log"
	_ "regexp"
	"strings"
	_ "strings"
	"time"
	_ "time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type TokenLookupItem struct {
	Token string
	Email string
}
type UserNote struct {
	Id          string `json:"id"`
	User        string `json:"user"`
	Create_date string `json:"create_date"`
	Text        string `json:"text"`
}

func GetDynamoDBClient() *dynamodb.DynamoDB {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	dynamodb.New(sess)
}

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
		Key: map[string]*dynamodb.AttributeValue{
			"token": {
				S: aws.String(token),
			},
		},
	})

	if err != nil {
		logger.Println("DynamoDB Error!", err)
		return "", false
	}

	item := TokenLookupItem{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)

	// Validate the given token with one from the database
	// and return user email if the tokens match ...
	if item.Email != "" {
		return item.Email, true

	}
	return "", false
}

func QueryUserNotes(email string) []UserNote {
	dynamoDBClient := GetDynamoDBClient()

	if email == "" {
		return nil
	}

	//emon
	tableName := "Movies"
	minRating := 4.0
	year := 2013
	// Create the Expression to fill the input struct with.
	// Get all movies in that year; we'll pull out those with a higher rating later
	filt := expression.Name("Year").Equal(expression.Value(year))

	// Or we could get by ratings and pull out those with the right year later
	//    filt := expression.Name("info.rating").GreaterThan(expression.Value(min_rating))

	// Get back the title, year, and rating
	proj := expression.NamesList(expression.Name("Title"), expression.Name("Year"), expression.Name("Rating"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	if err != nil {
		log.Fatalf("Got error building expression: %s", err)
	}

	//end

	// User the following date format for "now"
	dateNow := time.Now().Format(time.RFC3339)

	userNotes := []UserNote{}

	output, err := dynamoDBClient.Query(&dynamodb.QueryInput{
		TableName: aws.String("user-notes"),

		Limit: aws.Int64(10),
		KeyConditions: map[string]*dynamodb.Condition{
			"user": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(email),
					},
				},
			},
			"date": {},
		},
		ScanIndexForward: aws.Bool(false),
	})

	if err != nil {
		logger.Println("error getting notes", err.Error())
		return nil
	}
	if err := dynamodbattribute.UnmarshalListOfMaps(output.Items, &userNotes); err != nil {
		return nil
	}
	return userNotes
}

func AuthenticateUser(headers map[string]string) (string, error) {
	// You can get Authentication header in the following way:
	authenticationHeader := headers["Authentication"]
	if authenticationHeader == "" || !strings.HasPrefix("Bearer", authenticationHeader) {
		return "", fmt.Errorf("malformed")
	}
	// Validate the Authentication header and retrieve token

	// TODO: if authentication header is malformed or missing then return 400
	token := strings.TrimLeft(authenticationHeader, "Bearer")
	logger.Println("token ", token)
	if strings.TrimSpace(token) == "" {
		return "", fmt.Errorf("malformed")
	}
	email, ok := GetAuthenticatedUserEmail(token)
	if !ok {
		return "", fmt.Errorf("token not matched")
	}

	return email, nil
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := request.Headers

	email, err := AuthenticateUser(headers)

	if err != nil {
		if err.Error() == "malformed" {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       err.Error(),
			}, err
		} else {
			return events.APIGatewayProxyResponse{
				StatusCode: 403,
				Body:       err.Error(),
			}, err
		}
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
