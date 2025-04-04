package main

import (
	"fmt"
	"os"

	routes "github.com/dayalubajpai/jwtlearninggo/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.GET("/api-1", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"sucess": "Access granted for API-1"})
	})

	router.GET("/api-2", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"sucess": "Access granted for API-2"})
	})

	fmt.Printf("Server is running on port %s\n", port)
	router.Run(":" + port)
}
