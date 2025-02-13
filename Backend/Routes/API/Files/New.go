package Routes

import (
	"elucid503/SproutCDN/Models"
	"elucid503/SproutCDN/Types"

	"github.com/gin-gonic/gin"
)

func NewFile(GinContext *gin.Context) {

	Account, Exists := GinContext.Get("Account")

	if !Exists {

		GinContext.JSON(400, Types.Response{
			Message: "Could not get your account",
		})

		return

	}

	User := Account.(*Models.SproutAccount)

	// Get form data

	FileHeader, FileHeaderError := GinContext.FormFile("File")
	Path, PathExists := GinContext.GetPostForm("Path")
	Private, PrivateExists := GinContext.GetPostForm("Private")

	// Parse any data

	var PrivateBool bool = PrivateExists && Private == "true"

	if Path == "" || !PathExists || FileHeaderError != nil {

		GinContext.JSON(400, Types.Response{
			Message: "Invalid request",
		})

		return

	}

	// Save file

	var Bytes []byte = make([]byte, FileHeader.Size)

	File, FileError := FileHeader.Open()

	if FileError != nil {

		GinContext.JSON(400, Types.Response{
			Message: "Could not open file",
		})

		return

	}

	_, ReadError := File.Read(Bytes) // File.Read will put the file data into Bytes

	if ReadError != nil {

		GinContext.JSON(400, Types.Response{
			Message: "Could not read file",
		})

		return

	}

	// Create file struct and write

	NewFile := Models.NewFile(FileHeader.Filename, FileHeader.Size, PrivateBool, []string{User.UID}, Path)

	// Write the file

	WriteError := NewFile.Write(User.UID, Bytes)

	if WriteError != nil {

		GinContext.JSON(400, Types.Response{
			Message: WriteError.Error(),
		})

		return

	}

	GinContext.JSON(200, Types.Response{

		Message: "File created",

		JSON: NewFile,
		HTML: NewFile.ToHTML(),
	})

}
