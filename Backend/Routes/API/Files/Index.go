package Routes

import (
	"elucid503/SproutCDN/Functions"
	"elucid503/SproutCDN/Models"
	"elucid503/SproutCDN/Types"
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

func AllFiles(GinContext *gin.Context) {

	Account, Exists := GinContext.Get("Account")

	if !Exists {

		GinContext.JSON(401, Types.Response{
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

	Directory, ErrorLoadingDir := Models.LoadDirFromDotInfo(PathToDir)

	if ErrorLoadingDir != nil {

		GinContext.JSON(500, Types.Response{

			Message: "Could not load/find directory",
		})

		return

	}

	Contents, ErrLoadingContents := Directory.GetContents(User.UID)

	if ErrLoadingContents != nil {

		GinContext.JSON(500, Types.Response{

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

		CurrentContentItemPath := filepath.Join(Directory.Path, ContentItem.Name)

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

	PrivateIndicatorDisplay := "none"

	if Directory.Private == true {

		PrivateIndicatorDisplay = "block"

	}

	var FileListHTMLHeader string = fmt.Sprintf(`

		<div class="Container Transparent VerticalFlex DashMainHeaderContent"> 
					
			<div class="DashMainHeaderName DashCodeInfill">/%s</div>

				<ul class="DirStats">

					<li class="DashMainHeaderDetailItem SubFiles">%s</li>
					<li class="DashMainHeaderDetailItem SubDirs">%s</li>
					<li class="InlineFileStat PrivateIndicator" style="display:%s">Private</li>
				
				</ul>

		</div>

		<div class="Container Transparent Center DashMainHeaderIcon">

			<ion-icon name="folder-open-outline"></ion-icon>

		</div>
	
	`, strings.Replace(Functions.RemovePathFromStore(Directory.NormalizedPath), "/Store", "", 1), fmt.Sprintf("%d %s", len(FileArr), Functions.PluralizeString(len(FileArr), "File")), fmt.Sprintf("%d %s", len(DirArr), Functions.PluralizeString(len(DirArr), "Folder")), PrivateIndicatorDisplay)

	var FinalHTML string = fmt.Sprintf(`
	
		<div class="Container HorizontalFlex DashMainHeader AlternateDirTarget">

			%s

		</div>

		<div class="DashMainFileList">

			%s

		</div>
	
	`, FileListHTMLHeader, strings.Join(append(DirHTMLArr, FileHTMLArr...), "\n"))

	GinContext.JSON(200, Types.Response{

		Message: "Loaded all files",

		JSON: map[string]interface{}{

			"Parent": Directory,

			"Files": FileArr,
			"Dirs":  DirArr,
		},

		HTML: FinalHTML,
	})

}
