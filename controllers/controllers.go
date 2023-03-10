package controllers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
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


		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_ID  = user.ID.Hex()
		token, refresh_token, _ := helper.GenerateTokens(*user.Email, *user.First_Name, *user.Last_Name, *user.User_type,*&user.User_ID)
		user.Token = &token
		user.Refresh_token = &refresh_token

		resultInsertionNUmber, insertErr := userCollection.InsertOne(c, user)
		if insertErr != nil {
			msg := fmt.Sprintf("Error: User not created")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return		
		}
		defer cancel()
		ctx.JSON(http.StatusOK, resultInsertionNUmber)
	}
}


func LogIn() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User


		if err := ctx.BindJSON(&user); err != nil{
			ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()} )
		}

		err := userCollection.FindOne(c, bson.M{"error": user.Email}).Decode(&foundUser)
		defer cancel()
		if err!= nil{
			ctx.JSON(http.StatusInternalServerError, gin.H{"error" : "email or password is incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if passwordIsValid != true{
			ctx.JSON(http.StatusInternalServerError, gin.H{"Error": msg})
			return
		}

		if foundUser.Email == nil {
			ctx.JSON(http.StatusInternalServerError, gin.H {"Error": "user not found"})
		}

		token, refresh_token := helper.GenerateTokens(*foundUser.Email, *foundUser.First_Name, *foundUser.Last_Name, *foundUser.User_type, foundUser.User_ID)
		helper.UpdateTokens(token, refresh_token, foundUser.User_ID)
		err = userCollection.FindOne(c, bson.M{"userId": foundUser.User_ID}).Decode(&foundUser)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"Error" :  err.Error()})
			return
			
		}

		ctx.JSON(http.StatusOK, foundUser)
	}
}


func GetUsers() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		if err := helper.CheckUserType(ctx, "ADMIN"); err != nil  {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
			
		}

		var c,cancel = context.WithTimeout(context.Background(), 100*time.Second)
		recordPerPage, err := strconv.Atoi(ctx.Query("recordPerPage"))

		if err != nil || recordPerPage < 1{
			recordPerPage = 10
		}
		page, err := strconv.Atoi(ctx.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}

		startIndex := (page -  1) * recordPerPage
		startIndex, err = strconv.Atoi(ctx.Query("startIndex"))
		
		matchStage := bson.D{{"$match", bson.D{{}}}}
		groupStage := bson.D{{"$group", bson.D{
			{"_id", bson.D{{"_id", "null"}}}, 
			{"total_count", bson.D{{"$sum", 1}}}, 
			{"data", bson.D{{"$push", "$$ROOT"}}}}}}
		projectStage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},}},

	}

	result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage, groupStage, projectStage,
	})
	defer cancel()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"User": "error occurred when listing users"})

	}

	var allUsers []bson.M
	if err = result.All(ctx, &allUsers);( err!= nil) {
		log.Fatal(err)
		
	}
	ctx.JSON(http.StatusOK, allUsers[0])


}
}

func GetUserById() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		userId := ctx.Param("id")

		if err := helper.MatchUserTypeById(ctx, userId); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return
		}

		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		err := userCollection.FindOne(c, bson.M{"id":userId }).Decode(&user)
		defer cancel()

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, user)
	}
}