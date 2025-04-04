package routes

import (
	"github.com/dayalubajpai/jwtlearninggo/controllers"
	middleware "github.com/dayalubajpai/jwtlearninggo/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(importRoutes *gin.Engine) {
	importRoutes.Use(middleware.Authenticate())
	importRoutes.GET("users", controllers.GetUsers())
	importRoutes.GET("users/:user_id", controllers.GetUser())
}
