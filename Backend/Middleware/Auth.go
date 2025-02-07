package Middleware

import (
	"github.com/gin-gonic/gin"
)

func Authorize() gin.HandlerFunc {

	return func(GinContext *gin.Context) {

		// Check if the user is authorized
		// If the user is not authorized, return a 401 Unauthorized response
		// If the user is authorized, continue to the next middleware or the handler function
		GinContext.Next()

	}

}
