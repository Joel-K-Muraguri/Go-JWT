package helper

import (
	"context"
	"log"
	"os/user"
	"time"

	"github.com/Joel-K-Muraguri/go-jwt/database"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

type SignedDetails struct{
	Email string
	First_Name string
	Last_Name string
	User_type string
	jwt.StandardClaims
}

func GenerateTokens(email string, first_name string, last_name string, user_type string, uuid string)(signedToken string, signedRefreshToken string, err error){
	claims := &SignedDetails{
		Email: email,
		First_Name: first_name,
		Last_Name: last_name,
		User_type: user_type,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt : time.Now().Local().Add(time.Hour + time.Duration(24)).Unix(),
		},
	}

}


func UpdateTokens(signedToken string, signedRefreshToken string, userId string){
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second) 

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{"token" , signedToken} )
	updateObj = append(updateObj, bson.E{"Refresh Token", signedRefreshToken})

	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"Updated at", Updated_at})

	upsert := true
	filter := bson.M{"userId" : userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := userCollection.UpdateOne(
		ctx, 
		filter,
		bson.D{
			{"$set", updateObj},

		},
		&opt,
	)

	defer cancel()

	if err != nil{
		log.Panic(err)
		return
	}
	return

}