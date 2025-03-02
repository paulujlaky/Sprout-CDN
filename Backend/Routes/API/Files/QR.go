package Routes

import (
	"elucid503/SproutCDN/Models"
	"elucid503/SproutCDN/Types"
	"slices"

	"github.com/gin-gonic/gin"
)

func GetFileQR(GinContext *gin.Context) {

	// Get request body

	Body := map[string]any{}

	GinContext.BindJSON(&Body)

	Account, Exists := GinContext.Get("Account")

	if !Exists {

		GinContext.JSON(401, Types.Response{
			Message: "Could not get your account",
		})

		return

	}

	User := Account.(*Models.SproutAccount)

	PathToFile, Exists := Body["Path"].(string)

	if PathToFile == "" || !Exists {

		GinContext.JSON(400, Types.Response{

			Message: "Invalid request",
		})

		return

	}

	// Attempt to load file

	File, ErrLoadingFile := Models.LoadFileFromDotInfo(PathToFile)

	if ErrLoadingFile != nil {

		GinContext.JSON(500, Types.Response{

			Message: ErrLoadingFile.Error(),
		})

		return

	}

	if File.Private && !slices.Contains(File.Authorized, User.Username) {

		GinContext.JSON(404, Types.Response{

			Message: "File not found",
		})
	}

	// Generate QR code

	QRCode, ErrGeneratingQR := File.GenerateQR()

	if ErrGeneratingQR != nil {

		GinContext.JSON(500, Types.Response{

			Message: ErrGeneratingQR.Error(),
		})

		return

	}

	GinContext.Data(200, "image/png", QRCode)

}
