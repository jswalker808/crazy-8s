package main

import (
	"context"
	"crazy-8s/api"
	"crazy-8s/repository"
	"crazy-8s/service"
	"crazy-8s/transport"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type APIGatewayRequest = events.APIGatewayWebsocketProxyRequest
type APIGatewayResponse = events.APIGatewayProxyResponse

var router *api.Router

func handler(request APIGatewayRequest) (APIGatewayResponse, error) {

	log.Println("In the handler func")

	switch routeKey := request.RequestContext.RouteKey; routeKey {
		case "$connect": 
			return APIGatewayResponse{
				StatusCode: 200,
				Body: "Successfully connected",
			}, nil
		case "$disconnect":
			return APIGatewayResponse{
				StatusCode: 200,
				Body: "Successfully disconnected",
			}, nil
		default:
			return handleGamePlay(request)
	}
}

func handleGamePlay(apiGatewayRequest APIGatewayRequest) (APIGatewayResponse, error) {
	log.Printf("Request: %v", apiGatewayRequest.Body)

	baseRequest, err := transport.NewBaseRequest(apiGatewayRequest.Body)
	if (err != nil) {
		return APIGatewayResponse {
			StatusCode: 400,
			Body: err.Error(),
		}, err
	}

	gameRequest, err := transport.NewGameRequest(baseRequest)
	if (err != nil) {
		return APIGatewayResponse {
			StatusCode: 400,
			Body: err.Error(),
		}, err
	}
	
	if err := router.HandleRequest(baseRequest.Action, gameRequest); err != nil {
		return APIGatewayResponse {
			StatusCode: 400,
		}, err
	}

	return APIGatewayResponse {
		StatusCode: 200,
	}, nil
}

func main() {
	awsConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-1"))
	if err != nil {
		log.Panic("Unable to load aws configuration")
	}
	dynamoDbClient := dynamodb.NewFromConfig(awsConfig)
	gameRepository := repository.NewGameRepository(dynamoDbClient)
	gameService := service.NewGameService(gameRepository)
	router = api.NewRouter(gameService)
	lambda.Start(handler)
}
