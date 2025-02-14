package Routes

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func AllFiles(GinContext *gin.Context) {

	Dir, Exists := GinContext.Params.Get("Path")

	if Dir == "" || !Exists {

		Dir = "/"

	}

	GinContext.JSON(200, gin.H{"message": fmt.Sprintf("All files in %s", Dir)})

}
