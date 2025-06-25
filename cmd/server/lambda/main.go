package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"

	"github.com/peter-novosel/js-script-service/internal/config"
	"github.com/peter-novosel/js-script-service/internal/db"
	"github.com/peter-novosel/js-script-service/internal/logger"
	"github.com/peter-novosel/js-script-service/internal/router"
)

var adapter *ginadapter.ChiLambda

func init() {
	// Initialize logger
	log := logger.Init()

	// Load config and DB
	cfg := config.Load()
	if err := db.Init(cfg); err != nil {
		log.Fatalf("DB init failed in Lambda: %v", err)
	}

	// Set up router and wrap with Lambda adapter
	adapter = ginadapter.New(router.Setup(cfg))
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return adapter.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
