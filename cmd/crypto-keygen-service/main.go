package main

import (
	"crypto-keygen-service/internal/controller"
	"crypto-keygen-service/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	keyService := service.NewKeyService()

	router := gin.Default()
	controller.RegisterRoutes(router, keyService)
	router.Run(":8080")
}
