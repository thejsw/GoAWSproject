package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-openai-server/models"
	"go-openai-server/service"
)

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

	result, err :=
		service.Generate(
			req.Topic,
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
			Result: result,
		},
	)
}
