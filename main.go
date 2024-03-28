package main

import (
	"hacktiv8-techrawih-go-product-sale/config"
	"hacktiv8-techrawih-go-product-sale/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the database
	db := config.GetDBConnection()

	// Set up Gin engine
	r := gin.Default()

	// Register all application routes
	router.RegisterAPIService(r, db)

	// Start the server
	r.Run(":3000")
}
