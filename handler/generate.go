package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-openai-server/models"
	"go-openai-server/repository"
	"go-openai-server/service"
)

var wordRepo *repository.WordRepository

func SetWordRepository(repo *repository.WordRepository) {
	wordRepo = repo
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

	if wordRepo == nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "word repository is not initialized",
			},
		)

		return
	}

	words, err :=
		service.Generate(
			c.Request.Context(),
			req.Topic,
			wordRepo,
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
		models.GenerateResponse{
			Words: words,
		},
	)
}
