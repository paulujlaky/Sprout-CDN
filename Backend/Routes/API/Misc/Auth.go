package Routes

import (
	"elucid503/SproutCDN/Types"

	"github.com/gin-gonic/gin"
)

func Authorize(GinContext *gin.Context) {

	GinContext.JSON(200, Types.Response{

		Message: "Auth (maybe)!",
	})

}
