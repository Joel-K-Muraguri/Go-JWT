package routes

import (
	"github.com/Joel-K-Muraguri/go-jwt/controllers"
	"github.com/Joel-K-Muraguri/go-jwt/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine){

	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users", controllers.GetUsers())
	incomingRoutes.GET("user/id", controllers.GetUserById())

}