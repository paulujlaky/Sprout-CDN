package Routes

import (
	"elucid503/SproutCDN/Functions"
	"elucid503/SproutCDN/Models"
	"elucid503/SproutCDN/Types"

	"github.com/gin-gonic/gin"
)

func UpdateAccess(GinContext *gin.Context) {

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

	Private, PrivateExists := Body["Private"].(bool)
	// var _ []string = Body["Collaborators"].([]string) // Collaborators. Not to be used yet, but will be in the future

	if !PathExists {

		GinContext.JSON(400, Types.Response{

			Message: "Invalid request",
		})

		return

	}

	// Attempt to rename file

	FileToUpdate, ErrLoadingFile := Models.LoadFileFromDotInfo(Functions.AdjustPathToStore(Path))

	if ErrLoadingFile != nil {

		GinContext.JSON(500, Types.Response{

			Message: "Could not load/find file",
		})

		return

	}

	if !PrivateExists {

		Private = FileToUpdate.Private // Keep existing value

	}

	// Update the Dir

	ErrorUpdatingAccess := FileToUpdate.UpdateAccess(User.UID, Private) // TODO: Eventually add Collaborators

	if ErrorUpdatingAccess != nil {

		GinContext.JSON(500, Types.Response{

			Message: ErrorUpdatingAccess.Error(),
		})

	}

	GinContext.JSON(200, Types.Response{

		Message: "File updated",

		JSON: map[string]any{

			"UID": FileToUpdate.UID,
		},
	})

}
