package Middleware

import (
	"elucid503/SproutCDN/Models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var AccountCache map[string]Models.SproutAccount = make(map[string]Models.SproutAccount) // Will hold recent accounts so we don't have to make costly external requests every time

func CheckAccountFromCache(Token string) (Models.SproutAccount, bool) {

	Account, Exists := AccountCache[Token]

	return Account, Exists

}

func AddAccountToCache(Token string, Account Models.SproutAccount) {

	AccountCache[Token] = Account

	RemoveAccountFromCacheAfter(Token, 60*60) // 1 hour

}

func RemoveAccountFromCacheAfter(Token string, Timeout int) {

	go func() {

		// goroutine to run in background after timeout

		<-time.After(time.Second * time.Duration(Timeout)) // write to channel; everything after this will run after Timeout

		delete(AccountCache, Token)

	}()
}

func Authorize() gin.HandlerFunc {

	return func(GinContext *gin.Context) {

		Token, Error := GinContext.Cookie("Sprout-JWT")

		if Error != nil || Token == "" {

			// Try auth header

			Token = GinContext.GetHeader("Authorization")

			Token = strings.Split(Token, " ")[1] // Bearer token

		}

		if _, Exists := CheckAccountFromCache(Token); Exists {

			// Account is in cache

			GinContext.Set("Account", AccountCache[Token])

			GinContext.Next() // Early return
			return

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

		// Was successful, so add to cache

		AddAccountToCache(Token, Account)

		GinContext.Next()

	}

}
