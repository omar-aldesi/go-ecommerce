package main

import (
	"ecommerce/app/core"
	v1 "ecommerce/app/endpoints/v1"
	_ "ecommerce/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

// @title Go Ecommerce API
// @version 1.0
// @description Docs and examples for this project api.
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description "JWT token required. Format: Bearer {token}"
func main() {
	// define Gin
	r := gin.Default()

	// Init DB
	if err := core.InitDB(); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
		return
	}
	// define the api schema docs endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Register the routes
	v1.AuthRouter(r)
	v1.OrdersRouter(r)
	v1.ProductsRouter(r)
	v1.CategoriesRouter(r)
	v1.BranchesRouter(r)

	// Start the server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}
