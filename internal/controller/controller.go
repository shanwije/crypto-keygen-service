package controller

import (
	"net/http"

	"crypto-keygen-service/internal/service"
	"crypto-keygen-service/internal/util/errors"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

type GenerateRequest struct {
	UserID  int    `uri:"userId" validate:"required"`
	Network string `uri:"network" validate:"required,oneof=bitcoin ethereum"`
}

type KeyController struct {
	keyService *service.KeyService
}

func NewKeyController(keyService *service.KeyService) *KeyController {
	validate = validator.New()
	return &KeyController{keyService: keyService}
}

func (kc *KeyController) RegisterRoutes(router *gin.Engine) {
	router.GET("/generate/:userId/:network", kc.handleGenerateKeyPair)
}

func (kc *KeyController) handleGenerateKeyPair(c *gin.Context) {
	var req GenerateRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	address, publicKey, privateKey, err := kc.keyService.GetKeysAndAddress(req.UserID, req.Network)
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
