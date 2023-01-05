package helper

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func MatchUserTypeById(ctx *gin.Context, userId string)(err error){

	userType := ctx.GetString("user_type")
	uid := ctx.GetString("uid")
	err = nil

	if userType == "USER" && uid != userId {
		err = errors.New("unauthorized to access this resource")
		return err
	}
	err = CheckUserType(ctx, userType)
	return err


}


func CheckUserType(ctx *gin.Context, role string)(err error){
	userType := ctx.GetString("user_type")
	err = nil
	if userType != role {
		err = errors.New("unauthorized to access this resource")
		return err
	}
	return err

}