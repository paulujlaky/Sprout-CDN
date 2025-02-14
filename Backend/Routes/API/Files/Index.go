package Routes

import (
	"elucid503/SproutCDN/Models"
	"elucid503/SproutCDN/Types"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func AllFiles(GinContext *gin.Context) {

	Account, Exists := GinContext.Get("Account")

	if !Exists {

		GinContext.JSON(400, Types.Response{
			Message: "Could not get your account",
		})

		return

	}

	User := Account.(*Models.SproutAccount)

	PathToDir, Exists := GinContext.Params.Get("Path")

	if PathToDir == "" || !Exists {

		PathToDir = "/"

	}

	// Get all contents

	DirToLoad, ErrorLoadingDir := Models.LoadDirFromDotInfo(PathToDir)

	if ErrorLoadingDir != nil {

		GinContext.JSON(400, Types.Response{

			Message: "Could not load/find directory",
		})

		return

	}

	Contents, ErrLoadingContents := DirToLoad.GetContents(User.UID)

	if ErrLoadingContents != nil {

		GinContext.JSON(400, Types.Response{

			Message: "Could not load/find contents",
		})

		return

	}

	// Get all files

	FileArr := []Models.File{}
	DirArr := []Models.Dir{}

	FileHTMLArr := []string{}
	DirHTMLArr := []string{}

	for _, ContentItem := range Contents {

		CurrentContentItemPath := filepath.Join(DirToLoad.Path, ContentItem.Name)

		if ContentItem.IsDir {

			if LoadedDir, ErrorLoadingDir := Models.LoadDirFromDotInfo(CurrentContentItemPath); ErrorLoadingDir == nil {

				DirArr = append(DirArr, *LoadedDir)
				DirHTMLArr = append(DirHTMLArr, LoadedDir.ToHTML())

			}

		} else {

			if LoadedFile, ErrorLoadingFile := Models.LoadFileFromDotInfo(CurrentContentItemPath); ErrorLoadingFile == nil {

				FileArr = append(FileArr, *LoadedFile)
				FileHTMLArr = append(FileHTMLArr, LoadedFile.ToHTML())

			}

		}

	}

	GinContext.JSON(200, Types.Response{

		Message: "Loaded all files",

		JSON: map[string]interface{}{

			"Files": FileArr,
			"Dirs":  DirArr,
		},

		HTML: strings.Join(append(DirHTMLArr, FileHTMLArr...), "\n"), // All elements. first dirs, then files
	})

}
