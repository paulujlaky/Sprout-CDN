package Routes

import (
	"elucid503/SproutCDN/Types"

	"github.com/gin-gonic/gin"
)

func FileInfo(GinContext *gin.Context) {

	var FileUID string = GinContext.Param("uid")

	GinContext.JSON(200, Types.Response{

		Message: "File Info: " + FileUID + "!",
	})

}
