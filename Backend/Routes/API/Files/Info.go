package Routes

import (
	"elucid503/SproutCDN/Types"

	"github.com/gin-gonic/gin"
)

func FileInfo(GinContext *gin.Context) {

	var FileID string = GinContext.Param("id")

	GinContext.JSON(200, Types.Response{

		Message: "File Info: " + FileID + "!",
	})

}
