package Routes

import (
	"elucid503/SproutCDN/Functions"
	"elucid503/SproutCDN/Models"
	"elucid503/SproutCDN/Types"

	"github.com/gin-gonic/gin"
)

func MoveDir(GinContext *gin.Context) {

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

	// Parse any data

	OldPath, OldPathExists := Body["OldPath"].(string)
	NewPath, NewPathExists := Body["NewPath"].(string)

	if !OldPathExists || !NewPathExists {

		GinContext.JSON(400, Types.Response{

			Message: "Invalid request",
		})

		return

	}

	// Attempt to move dir

	DirToMove, ErrorLoadingDir := Models.LoadDirFromDotInfo(Functions.AdjustPathToStore(OldPath))

	if ErrorLoadingDir != nil {

		GinContext.JSON(500, Types.Response{

			Message: "Could not load/find dir",
		})

		return

	}

	// Move the file

	if Err := DirToMove.Move(User.Username, NewPath); Err != nil {

		GinContext.JSON(500, Types.Response{

			Message: Err.Error(),
		})

		return

	}

	GinContext.JSON(200, Types.Response{

		Message: "Directory moved",

		JSON: map[string]any{

			"UID": DirToMove.UID,
		},
	})

}
