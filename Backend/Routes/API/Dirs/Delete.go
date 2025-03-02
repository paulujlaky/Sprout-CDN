package Routes

import (
	"elucid503/SproutCDN/Functions"
	"elucid503/SproutCDN/Models"
	"elucid503/SproutCDN/Types"

	"github.com/gin-gonic/gin"
)

func DeleteDir(GinContext *gin.Context) {

	// Get authed account

	Account, Exists := GinContext.Get("Account")

	if !Exists {

		GinContext.JSON(400, Types.Response{
			Message: "Could not get your account",
		})

		return

	}

	User := Account.(*Models.SproutAccount)

	// Get request body

	Body := map[string]any{}

	GinContext.BindJSON(&Body)

	// Get config vals

	Path, PathExists := Body["Path"].(string)

	if Path == "" || !PathExists {

		GinContext.JSON(401, Types.Response{

			Message: "Invalid request",
		})

		return

	}

	Path = Functions.SanitizePath(Path) // Prevent user from being able to go up

	// Get Dir

	DirToDelete, ErrLoadingDir := Models.LoadDirFromDotInfo(Path)

	if ErrLoadingDir != nil {

		GinContext.JSON(500, Types.Response{

			Message: "Could not load/find directory",
		})

		return

	}

	DeleteFileError := DirToDelete.Delete(User.UID)

	if DeleteFileError != nil {

		GinContext.JSON(500, Types.Response{

			Message: "Could not delete directory",
		})

		return

	}

	GinContext.JSON(200, Types.Response{

		Message: "Directory deleted",
		JSON: map[string]any{
			"Success": true,
		},
	})

}
