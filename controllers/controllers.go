package controllers

import (
	"github.com/Joel-K-Muraguri/go-jwt/database"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "User")
var validate = validator.New()


func HashPassword(){

}

func VerifyPassword(){

}

func SignUp() gin.HandlerFunc{

}


func LogIn() gin.HandlerFunc{
	
}


func GetUsers() gin.HandlerFunc{


}

func GetUserById() gin.HandlerFunc{


}