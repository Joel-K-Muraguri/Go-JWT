package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/Joel-K-Muraguri/go-jwt/controllers"
)

func UserRoutes(incomingRoutes *gin.Engine){

	incomingRoutes.POST("/users/login", controllers.LogIn())
	incomingRoutes.POST("/users/signup", controllers.SignUp())



}