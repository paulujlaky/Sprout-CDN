package Routes

import (
	"elucid503/SproutCDN/Functions"
	"elucid503/SproutCDN/Models"
	"elucid503/SproutCDN/Types"

	"github.com/gin-gonic/gin"
)

func RenameDir(GinContext *gin.Context) {

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
	NewName, NewNameExists := Body["Name"].(string)

	if !PathExists || !NewNameExists {

		GinContext.JSON(400, Types.Response{

			Message: "Invalid request",
		})

	}

	// Attempt to rename dir

	DirToRename, ErrLoadingDir := Models.LoadDirFromDotInfo(Functions.AdjustPathToStore(Path))

	if ErrLoadingDir != nil {

		GinContext.JSON(500, Types.Response{

			Message: "Could not load/find dir",
		})

		return

	}

	// Rename dir

	RenameError := DirToRename.Rename(User.UID, NewName)

	if RenameError != nil {

		GinContext.JSON(500, Types.Response{

			Message: "Could not rename dir",
		})

		return

	}

	GinContext.JSON(200, Types.Response{

		Message: "Dir renamed",

		JSON: map[string]any{

			"UID": DirToRename.UID,
		},
	})

}
