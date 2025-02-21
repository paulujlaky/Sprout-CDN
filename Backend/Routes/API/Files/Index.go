package Routes

import (
	"elucid503/SproutCDN/Models"
	"elucid503/SproutCDN/Types"
	"path/filepath"
	"slices"
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

			Message: ErrLoadingContents.Error(),
		})

		return

	}

	// Get all files

	FileArr := []Models.File{}
	DirArr := []Models.Dir{}

	FileHTMLArr := []string{}
	DirHTMLArr := []string{}

	for _, ContentItem := range Contents {

		// Skip info files

		if strings.HasSuffix(ContentItem.Name, ".fileinfo") || strings.HasSuffix(ContentItem.Name, ".dirinfo") {

			continue

		}

		CurrentContentItemPath := filepath.Join(DirToLoad.Path, ContentItem.Name)

		if ContentItem.IsDir {

			if LoadedDir, ErrorLoadingDir := Models.LoadDirFromDotInfo(CurrentContentItemPath); ErrorLoadingDir == nil {

				// Check permissions

				if LoadedDir.Private && !slices.Contains(LoadedDir.Authorized, User.UID) {

					continue

				}

				DirArr = append(DirArr, *LoadedDir)
				DirHTMLArr = append(DirHTMLArr, LoadedDir.ToHTML())

			}

		} else {

			if LoadedFile, ErrorLoadingFile := Models.LoadFileFromDotInfo(CurrentContentItemPath); ErrorLoadingFile == nil {

				if LoadedFile.Private && !slices.Contains(LoadedFile.Authorized, User.UID) {

					continue

				}

				FileArr = append(FileArr, *LoadedFile)
				FileHTMLArr = append(FileHTMLArr, LoadedFile.ToHTML())

			}

		}

	}

	GinContext.JSON(200, Types.Response{

		Message: "Loaded all files",

		JSON: map[string]interface{}{

			"Parent": DirToLoad,

			"Files": FileArr,
			"Dirs":  DirArr,
		},

		HTML: strings.Join(append(DirHTMLArr, FileHTMLArr...), "\n"), // All elements. first dirs, then files
	})

}
