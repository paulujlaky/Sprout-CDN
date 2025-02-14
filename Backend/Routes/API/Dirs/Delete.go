package Routes

import (
	"elucid503/SproutCDN/Functions"
	"elucid503/SproutCDN/Models"
	"elucid503/SproutCDN/Types"
	"path/filepath"

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

	Name, NameExists := Body["Name"].(string)
	Path, PathExists := Body["Path"].(string)

	if Name == "" || Path == "" || !NameExists || !PathExists {

		GinContext.JSON(401, Types.Response{

			Message: "Invalid request",
		})

		return

	}

	Path = Functions.SanitizePath(Path) // Prevent user from being able to go up

	// Check if dir exists. must already exist to delete

	DirToDeletePath := filepath.Join(Path, Name)

	// Get Dir

	DirToDelete, ErrLoadingDir := Models.LoadDirFromDotInfo(DirToDeletePath)

	if ErrLoadingDir != nil {

		GinContext.JSON(500, Types.Response{

			Message: "Could not load/find directory",
		})

		return

	}

	DeleteFileError, DeleteInfoError := DirToDelete.Delete(User.UID)

	if DeleteFileError != nil || DeleteInfoError != nil {

		GinContext.JSON(500, Types.Response{

			Message: DeleteFileError.Error() + "; " + DeleteInfoError.Error(),
		})

		return

	}

	GinContext.JSON(200, Types.Response{

		Message: "Directory deleted",
	})

}
