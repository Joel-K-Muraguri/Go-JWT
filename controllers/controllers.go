package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Joel-K-Muraguri/go-jwt/database"
	"github.com/Joel-K-Muraguri/go-jwt/helper"
	"github.com/Joel-K-Muraguri/go-jwt/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "User")
var validate = validator.New()


func HashPassword(password string) string{
	bytes, err := bcrypt.GenerateFromPassword([]byte (password), 14)

	if err != nil {
		log.Panic(err)	
	}
	return string(bytes)

}

func VerifyPassword(userPassword string, loginPassword string)(bool, string){
	err := bcrypt.CompareHashAndPassword([]byte( userPassword), []byte (loginPassword))
	msg := ""
	check := true

	if err != nil {
		msg = fmt.Sprint("Password for this email is wrong")
		check = false
	}

	return check, msg

}

func SignUp() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationError := validate.Struct(user)
		if validationError != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": validationError.Error()})
			return
		}

		count, err := userCollection.CountDocuments(c, bson.M{"email" : user.Email })
		defer cancel()
		if err != nil {
			log.Panic(err)
			ctx.JSON(http.StatusInternalServerError, gin.H { "error": "error occurred when processing your email "})	
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		count, err = userCollection.CountDocuments(c, bson.M{"phone" : user.Phone_Number })
		defer cancel()
		if err != nil {
			log.Panic(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred when processing your phone number "})	
		}

		if count > 0 {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error":"this email or phone number already exists"})
		}


		user.Created_at, _ = time.Parse(time.RFC1123, time.Now().Format(time.RFC1123))
		user.Updated_at, _ = time.Parse(time.RFC1123, time.Now().Format(time.RFC1123))
		user.ID = string(primitive.NewObjectID())
		user.User_ID  = user.ID.Hex()
		token, refresh_token, _ := helper.GenerateTokens()




	
	}
}


func LogIn() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		var c, err = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User



	}
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