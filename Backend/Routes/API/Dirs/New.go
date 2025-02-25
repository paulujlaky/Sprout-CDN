package Routes

import (
	"elucid503/SproutCDN/Functions"
	"elucid503/SproutCDN/Models"
	"elucid503/SproutCDN/Types"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func CreateDir(GinContext *gin.Context) {

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
	Private, _ := Body["Private"].(bool)

	if Name == "" || !NameExists || !PathExists {

		GinContext.JSON(400, Types.Response{

			Message: "Invalid request",
		})

		return

	}

	OriginalPath := Path
	Path = Functions.AdjustPathToStore(Functions.SanitizePath(Path)) // Prevent user from being able to go up and set CDN root

	DirUtil := Functions.DirUtil{}

	// Check if dir exists

	NewDirPath := filepath.Join(Path, Name)

	AlreadyExists := DirUtil.Exists(NewDirPath)

	if AlreadyExists {

		GinContext.JSON(400, Types.Response{

			Message: "Directory already exists",
		})

		return

	}

	// NewDir creates the directory and writes the .info file

	ModelRepresentation, InfoWriteError, CreateError := Models.NewDir(NewDirPath, filepath.Join(OriginalPath, Name), []string{User.UID}, Private)

	if InfoWriteError != nil || CreateError != nil {

		GinContext.JSON(500, Types.Response{

			Message: "Directory create error occurred",
		})

		return

	}

	GinContext.JSON(200, Types.Response{

		Message: "Directory created",

		JSON: ModelRepresentation,
		HTML: ModelRepresentation.ToHTML(),
	})

}
