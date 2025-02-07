package Middleware

import (
	"elucid503/SproutCDN/Models"

	"github.com/gin-gonic/gin"
)

func Authorize() gin.HandlerFunc {

	return func(GinContext *gin.Context) {

		Token, Error := GinContext.Cookie("Sprout-JWT")

		if Error != nil || Token == "" {

			// Try auth header

			Token = GinContext.GetHeader("Authorization")

		}

		if Token == "" {

			GinContext.JSON(401, gin.H{

				"Message": "Unauthorized",
			})

			GinContext.Abort()

			return

		}

		// Get account

		Account, Err := Models.GetSproutAccountByToken(Token)

		if Err != nil {

			GinContext.JSON(401, gin.H{

				"Message": "Could not get your account",
			})

			GinContext.Abort()

			return

		}

		GinContext.Set("Account", Account)

		GinContext.Next()

	}

}
