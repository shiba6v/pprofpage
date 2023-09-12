package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/shiba6v/pprofpage/app/registry"
)

var echoLambda *echoadapter.EchoLambda

func init() {
	e := echo.New()
	e.Use(middleware.Recover())
	s3cli, err := registry.RegisterS3()
	if err != nil {
		log.Fatalf("cannot initialize s3cli %v", err)
	}
	registry.RegisterServer(e, s3cli)
	echoLambda = echoadapter.New(e)
}

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return echoLambda.ProxyWithContext(ctx, req)
}
