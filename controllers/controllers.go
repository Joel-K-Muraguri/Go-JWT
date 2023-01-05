package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Joel-K-Muraguri/go-jwt/database"
	"github.com/Joel-K-Muraguri/go-jwt/helper"
	"github.com/Joel-K-Muraguri/go-jwt/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "User")
var validate = validator.New()


func HashPassword(){

}

func VerifyPassword(){

}

func SignUp() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		






	}

}


func LogIn() gin.HandlerFunc{
	
}


func GetUsers() gin.HandlerFunc{


}

func GetUserById() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		userId := ctx.Params("id")

		if err := helper.MatchUserTypeById(ctx, userId); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"id":userId }).Decode(&user)
		defer cancel()

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, user)




	}


}