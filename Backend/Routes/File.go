package Routes

import (
	"elucid503/SproutCDN/Middleware"
	"elucid503/SproutCDN/Models"
	"fmt"
	"path/filepath"
	"slices"

	"github.com/gin-gonic/gin"
)

// Templates

func GetFile(GinContext *gin.Context) {

	Path, PathParamExists := GinContext.Params.Get("File")

	if !PathParamExists {

		ReturnArbitraryNonExistsError(GinContext, "File")
		return

	}

	// Attempt to read file info

	FileRequested, ErrorRequesting := Models.LoadFileFromDotInfo(Path)

	if ErrorRequesting != nil {

		ReturnArbitraryNonExistsError(GinContext, Path)
		return

	}

	// Get dir above

	DirAbove, ErrorLoadingDirAbove := Models.LoadDirFromDotInfo(filepath.Dir(Path))

	if ErrorLoadingDirAbove != nil {

		ReturnArbitraryNonExistsError(GinContext, Path)

	}

	// Check if the dir

	if FileRequested.Private == true || DirAbove.Private == true {

		// Authorize request (only doing here to avoid higher loading overhead on all requests)

		Account, AuthError := Middleware.AuthorizeFromRequest(GinContext)

		if AuthError != nil {

			ReturnArbitraryNonExistsError(GinContext, Path)
			return

		}

		// Check if account is authorized by seeing if UID of acc is in UID slice of authorized users

		if !slices.Contains(FileRequested.Authorized, Account.UID) {

			ReturnArbitraryNonExistsError(GinContext, Path)
			return

		}

		// All private checks are passed

	}

	// Serve file

	GinContext.File(FileRequested.Path)

}

func ReturnArbitraryNonExistsError(GinContext *gin.Context, Path string) {

	GinContext.Header("Content-Type", "text/html")
	GinContext.String(404, GetFileNotFoundTemplate(Path))

}

func GetFileNotFoundTemplate(Path string) string {

	return fmt.Sprintf(`

	<!DOCTYPE html>
	<html lang="en">
	
	<head>
	
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
	
		<title>Sprout CDN</title>

		<link rel="stylesheet" href="https://use.typekit.net/pge8obf.css">

		<style>
	
			body {
	
				font-family: "din-2014", sans-serif;
				font-weight: 400;
				font-style: normal;

				display: flex;
				justify-content: center;
				align-items: center;
	
				height: 100vh;
				margin: 0;

				color: white;
	
				background-color:rgb(10, 10, 10);
		
			}

			h1 {

				margin-top: 15px;
				margin-bottom: 0;

			}
	
			.FileNotFoundPageContent {
	
				text-align: center;

				padding: 5px 20px;
				margin: 0 25px;

				max-width: 400px;

				background-color:rgb(20, 20, 20);

				border-radius: 7.5px;
				border: 1px solid rgb(50, 50, 50);
	
			}

			.FileNotFoundPageCodeBlock {
			
				background-color:rgb(50, 50, 50);

				padding: 2.5px 5px;

				border-radius: 5px;
				
			}
	
		</style>
	
	</head>
	
	<body>
	
		<div class="FileNotFoundPageContent">
	
			<h1>File Not Found</h1>
			<p>The file <span class="FileNotFoundPageCodeBlock">%s</span> was not found. The path may be incorrect, you may not have access, or it may no longer exist.</p>
		
		</div>
	
	</body>
	</html>`, Path)

}
