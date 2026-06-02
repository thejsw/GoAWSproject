package main

import (
	"context"
	"log"
	"os"

	"go-openai-server/config"
	"go-openai-server/handler"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda
var router *gin.Engine

func setupRouter() *gin.Engine {

	r := gin.Default()

	r.Use(cors.Default())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	r.POST("/generate", handler.Generate)

	return r
}

func init() {

	config.LoadEnv()

	router = setupRouter()

	ginLambda = ginadapter.New(router)
}

func LambdaHandler(
	ctx context.Context,
	req events.APIGatewayProxyRequest,
) (
	events.APIGatewayProxyResponse,
	error,
) {

	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {

	if os.Getenv("AWS_LAMBDA_RUNTIME_API") != "" {

		log.Println("Running on AWS Lambda")

		lambda.Start(LambdaHandler)

	} else {

		log.Println("Running locally on :8080")

		router.Run(":8080")
	}
}
