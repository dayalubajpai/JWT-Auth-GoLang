package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func CheckUserType(c *gin.Context, role string) (err error) {
	userType := c.GetString("user_type")
	err = nil
	if userType != role {
		err = errors.New("Unauthorized to access the resources")
	}
	return err
}

func MatchUserTypeToUID(ctx *gin.Context, userId string) (err error) {

	var userType = ctx.GetString("user_type")
	var uid = ctx.GetString("uid")

	err = nil

	if userType == "USER" && uid != userId {
		err = errors.New("Unauthorized to access the resources")
	}

	err = CheckUserType(ctx, userType)

	return err
}
