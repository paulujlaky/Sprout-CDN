package Routes

import (
	"elucid503/SproutCDN/Models"
	"elucid503/SproutCDN/Types"

	"github.com/gin-gonic/gin"
)

func Authorize(GinContext *gin.Context) {

	// Get request body

	Body := gin.H{}

	GinContext.BindJSON(&Body)

	// Get token

	Token, TokenExists := Body["Token"].(string)

	if Token == "" || !TokenExists {

		GinContext.JSON(401, Types.Response{

			Message: "Unauthorized",
		})

		return

	}

	// Get account

	Account, Err := Models.GetSproutAccountByToken(Token)

	if Err != nil {

		GinContext.JSON(401, Types.Response{

			Message: "Could not get your account",
		})

		return

	}

	GinContext.SetCookie("Sprout-JWT", Token, 60*60*24*7, "/", "", false, true)

	GinContext.JSON(200, Account)

}
