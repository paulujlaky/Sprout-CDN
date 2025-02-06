package Routes

import (
	"elucid503/SproutCDN/Types"

	"github.com/gin-gonic/gin"
)

func NewFile(GinContext *gin.Context) {

	GinContext.JSON(200, Types.Response{

		Message: "New File!",
	})

}
