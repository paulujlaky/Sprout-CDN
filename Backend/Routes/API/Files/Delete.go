package Routes

import (
	"elucid503/SproutCDN/Models"
	"elucid503/SproutCDN/Types"

	"github.com/gin-gonic/gin"
)

func DeleteFile(GinContext *gin.Context) {

	Account, Exists := GinContext.Get("Account")

	if !Exists {

		GinContext.JSON(401, Types.Response{
			Message: "Could not get your account",
		})

		return

	}

	User := Account.(*Models.SproutAccount)

	// Get request body

	Body := map[string]any{}

	GinContext.BindJSON(&Body)

	// Parse any data

	Path, PathExists := Body["Path"].(string)

	if Path == "" || !PathExists {

		GinContext.JSON(400, Types.Response{
			Message: "Invalid request",
		})

		return

	}

	// Attempt to delete file

	FileToDelete, ErrLoadingFile := Models.LoadFileFromDotInfo(Path)

	if ErrLoadingFile != nil {

		GinContext.JSON(500, Types.Response{

			Message: "Could not load/find file",
		})

		return

	}

	DeleteFileError, DeleteInfoError := FileToDelete.Delete(User.UID) // user must have permission. will check and return precondition failed if not

	if DeleteFileError != nil || DeleteInfoError != nil {

		GinContext.JSON(500, Types.Response{

			Message: "Could not delete file",
		})

		return

	}

	GinContext.JSON(200, Types.Response{

		Message: "File deleted",
		JSON: map[string]any{
			"Success": true,
		},
	})

}
