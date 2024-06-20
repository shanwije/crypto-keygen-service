package handlers

import (
	"net/http"
	"strconv"

	"crypto-keygen-service/internal/services"
	"crypto-keygen-service/internal/util/errors"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

type KeyGenRequest struct {
	UserID  int    `uri:"userId" validate:"required,gt=0"`
	Network string `uri:"network" validate:"required"`
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

	// Validate userId manually to handle non-integer values gracefully
	userIdStr := c.Param("userId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil || userId <= 0 {
		log.WithError(err).Error("Invalid userId parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrInvalidUserID.Message})
		return
	}
	req.UserID = userId

	// Validations for network and userId
	req.Network = c.Param("network")
	if req.Network == "" {
		log.Error("Network parameter is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrNetworkRequired.Message})
		return
	}

	if err := validate.Struct(req); err != nil {
		log.WithError(err).Error("Validation error")
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.FormatValidationError(err)})
		return
	}

	keyPairAndAddress, err := h.keyService.GetKeysAndAddress(req.UserID, req.Network)
	if err != nil {
		handleServiceError(c, err, req.UserID, req.Network)
		return
	}

	log.WithFields(log.Fields{
		"user_id": req.UserID,
		"network": req.Network,
	}).Info("Successfully acquired keys")

	c.JSON(
		http.StatusOK,
		KeyGenResponse{
			Address:    keyPairAndAddress.Address,
			PublicKey:  keyPairAndAddress.PublicKey,
			PrivateKey: keyPairAndAddress.PrivateKey,
		},
	)
}

func handleServiceError(c *gin.Context, err error, userID int, network string) {
	if apiErr, ok := err.(*errors.KeyGenError); ok {
		log.WithFields(log.Fields{
			"user_id": userID,
			"network": network,
		}).WithError(apiErr).Error("API error")
		c.JSON(apiErr.Code, gin.H{"error": apiErr.Message})
	} else {
		log.WithFields(log.Fields{
			"user_id": userID,
			"network": network,
		}).WithError(err).Error("Internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrInternalServerError.Message})
	}
}
