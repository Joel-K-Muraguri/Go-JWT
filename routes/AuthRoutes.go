package routes

import (
	"github.com/Joel-K-Muraguri/go-jwt/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine){

	incomingRoutes.POST("/users/login", controllers.LogIn())
	incomingRoutes.POST("/users/signup", controllers.SignUp())
}