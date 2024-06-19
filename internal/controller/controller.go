package controller

import (
	"crypto-keygen-service/internal/errors"
	"crypto-keygen-service/internal/service"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

var validate *validator.Validate

type GenerateRequest struct {
	UserID  int    `uri:"userId" validate:"required"`
	Network string `uri:"network" validate:"required,oneof=bitcoin ethereum"`
}

func RegisterRoutes(router *gin.Engine, keyService *service.KeyService) {
	router.GET("/generate/:userId/:network", func(c *gin.Context) {
		handleGenerateKeyPair(c, keyService)
	})
}

func handleGenerateKeyPair(c *gin.Context, keyService *service.KeyService) {
	validate = validator.New()

	var req GenerateRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := req.UserID
	network := req.Network

	address, publicKey, privateKey, err := keyService.GenerateKeyPair(userID, network)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.Code, gin.H{"error": apiErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrInternalServerError.Message})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"address":     address,
		"public_key":  publicKey,
		"private_key": privateKey,
	})
}
