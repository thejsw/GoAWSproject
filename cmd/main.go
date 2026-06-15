package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"go-openai-server/config"
	"go-openai-server/handler"
	"go-openai-server/repository"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ginLambda *ginadapter.GinLambda
var router *gin.Engine
var initOnce sync.Once
var initErr error

func setupRouter(quizRepo *repository.QuizRepository) *gin.Engine {

	r := gin.Default()

	handler.SetQuizRepository(quizRepo)

	r.Use(cors.Default())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	r.POST("/generate", handler.Generate)

	return r
}

func initialize() error {

	if err := config.LoadEnv(); err != nil {
		return err
	}

	pool, err := pgxpool.New(
		context.Background(),
		os.Getenv("DATABASE_URL"),
	)
	if err != nil {
		return err
	}

	quizRepo := repository.NewQuizRepository(pool)
	router = setupRouter(quizRepo)
	ginLambda = ginadapter.New(router)

	return nil
}

func ensureInitialized() error {
	initOnce.Do(func() {
		initErr = initialize()
	})

	return initErr
}

func LambdaHandler(
	ctx context.Context,
	req events.APIGatewayProxyRequest,
) (
	events.APIGatewayProxyResponse,
	error,
) {

	if err := ensureInitialized(); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf(`{"error":%q}`, err.Error()),
		}, nil
	}

	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	if err := ensureInitialized(); err != nil {
		log.Fatal(err)
	}

	if os.Getenv("AWS_LAMBDA_RUNTIME_API") != "" {

		log.Println("Running on AWS Lambda")

		lambda.Start(LambdaHandler)

	} else {

		log.Println("Running locally on :8080")

		router.Run(":8080")
	}
}
