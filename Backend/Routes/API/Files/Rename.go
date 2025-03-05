package Routes

import (
	"elucid503/SproutCDN/Functions"
	"elucid503/SproutCDN/Models"
	"elucid503/SproutCDN/Types"

	"github.com/gin-gonic/gin"
)

func RenameFile(GinContext *gin.Context) {

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

		return

	}

	// Attempt to rename file

	FileToRename, ErrLoadingFile := Models.LoadFileFromDotInfo(Functions.AdjustPathToStore(Path))

	if ErrLoadingFile != nil {

		GinContext.JSON(500, Types.Response{

			Message: "Could not load/find file",
		})

		return

	}

	// Move the file

	if Err := FileToRename.Rename(User.UID, NewName); Err != nil {

		GinContext.JSON(500, Types.Response{

			Message: Err.Error(),
		})

		return

	}

	GinContext.JSON(200, Types.Response{

		Message: "File renamed",

		JSON: map[string]any{

			"UID": FileToRename.UID,
		},
	})

}
