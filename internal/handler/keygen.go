package handler

import (
	"net/http"

	"crypto-keygen-service/internal/services"
	"crypto-keygen-service/internal/util/errors"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

type KeyGenRequest struct {
	UserID  int    `uri:"userId" validate:"required"`
	Network string `uri:"network" validate:"required,oneof=bitcoin ethereum"`
}

type KeyGenHandler struct {
	keyService *services.KeyGenService
}

func NewKeyGenHandler(keyService *services.KeyGenService) *KeyGenHandler {
	validate = validator.New()
	return &KeyGenHandler{keyService: keyService}
}

func (h *KeyGenHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/keygen/:userId/:network", h.handleGenerateKeyPair)
}

func (h *KeyGenHandler) handleGenerateKeyPair(c *gin.Context) {
	var req KeyGenRequest
	if err := c.ShouldBindUri(&req); err != nil {
		log.WithError(err).Error("Invalid request parameters")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	if err := validate.Struct(req); err != nil {
		log.WithError(err).Error("Validation error")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	address, publicKey, privateKey, err := h.keyService.GetKeysAndAddress(req.UserID, req.Network)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			log.WithFields(log.Fields{
				"user_id": req.UserID,
				"network": req.Network,
			}).WithError(apiErr).Error("API error")
			c.JSON(apiErr.Code, gin.H{"error": apiErr.Message})
		} else {
			log.WithFields(log.Fields{
				"user_id": req.UserID,
				"network": req.Network,
			}).WithError(err).Error("Internal server error")
			c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrInternalServerError.Message})
		}
		return
	}

	log.WithFields(log.Fields{
		"user_id": req.UserID,
		"network": req.Network,
	}).Info("Successfully acquired keys")

	c.JSON(http.StatusOK, gin.H{
		"address":     address,
		"public_key":  publicKey,
		"private_key": privateKey,
	})
}
