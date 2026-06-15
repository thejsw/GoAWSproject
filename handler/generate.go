package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-openai-server/models"
	"go-openai-server/repository"
	"go-openai-server/service"
)

var quizRepo *repository.QuizRepository

func SetQuizRepository(repo *repository.QuizRepository) {
	quizRepo = repo
}

func Generate(c *gin.Context) {

	var req models.GenerateRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	if quizRepo == nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "quiz repository is not initialized",
			},
		)

		return
	}

	bundle, err :=
		service.GeneratePart5Quiz(
			c.Request.Context(),
			req.Topic,
			quizRepo,
		)

	if err != nil {

		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		bundle,
	)
}
