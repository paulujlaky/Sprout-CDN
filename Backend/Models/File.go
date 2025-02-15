package Models

import (
	"fmt"

	"elucid503/SproutCDN/Functions"
	"elucid503/SproutCDN/Types"
	"encoding/json"
	"path/filepath"
	"strings"

	"errors"
	"os"
	"slices"
)

var FileUtil = Functions.FileUtil{} // init FileUtil for use in File methods (package level)
var DirUtil = Functions.DirUtil{}   // init DirUtil for use in File methods

type File struct {
	UID  string `json:"UID"`
	Name string `json:"Name"`

	Size int64 `json:"Size"`

	Private bool `json:"Private"`

	Authorized []string `json:"Authorized"`

	Path string `json:"Path"`
}

// General

func NewFile(Name string, Size int64, Private bool, Authorized []string, Path string) *File {

	return &File{

		UID:  Functions.RandomString(16),
		Name: Name,

		Size: Size,

		Private: Private,

		Authorized: Authorized,

		Path: filepath.Join(Functions.AdjustPathToStore(Path), Name),
	}

}

func LoadFileFromDotInfo(Path string) (*File, error) {

	Path = Functions.AdjustPathToStore(Path)

	// Load the file from the .info file

	InfoPath := getFileInfoPath(Path)

	InfoData, ReadError := os.ReadFile(InfoPath)

	if ReadError != nil {

		return nil, ReadError

	}

	var FileData File // new instance

	FileData.Path = Path

	UnmarshalError := json.Unmarshal(InfoData, &FileData)

	return &FileData, UnmarshalError

}

func getFileInfoPath(OriginalPath string) string {

	// We must remove the existing extension and add .info to the end

	LastPeriodIndex := strings.LastIndex(OriginalPath, ".")

	// Get the path without the extension

	ExtensionlessPath := OriginalPath[:LastPeriodIndex]

	// Add .info to the end

	return ExtensionlessPath + ".fileinfo"

}

// Private

func (AssociatedFile *File) checkPreconditions(RequestingUser string) (bool, bool) {

	PassesPrivateCheck := slices.Contains(AssociatedFile.Authorized, RequestingUser)
	PassesExistCheck := DirUtil.Exists(filepath.Dir(AssociatedFile.Path))

	if AssociatedFile.Private {

		return PassesPrivateCheck, PassesExistCheck

	} else {

		return true, PassesExistCheck

	}

}

func (AssociatedFile *File) getMimeTypeAndIcon() (string, string) {

	return FileUtil.GetMimeTypeAndIcon(filepath.Ext(AssociatedFile.Path))

}

func (AssociatedFile *File) writeInfo() error {

	// Write all info in the struct to a .info file which is later parsed as JSON

	DataToWrite, MarshalError := json.Marshal(AssociatedFile)

	if MarshalError != nil {

		return MarshalError

	}

	return os.WriteFile(getFileInfoPath(AssociatedFile.Path), DataToWrite, os.ModePerm)

}

func (AssociatedFile *File) deleteInfo() error {

	return os.Remove(getFileInfoPath(AssociatedFile.Path))

}

// Public

func (AssociatedFile *File) Write(RequestingUser string, Data []byte) error {

	if Exists, Authorized := AssociatedFile.checkPreconditions(RequestingUser); !Exists || !Authorized {

		return errors.New("Preconditions failed.")

	}

	WriteInfoErr := AssociatedFile.writeInfo()

	WriteFileErr := os.WriteFile(AssociatedFile.Path, Data, os.ModePerm)

	if WriteInfoErr != nil || WriteFileErr != nil {

		return errors.New("Failed to write file.")

	}

	return nil

}

func (AssociatedFile *File) Read(RequestingUser string) ([]byte, error) {

	if Exists, Authorized := AssociatedFile.checkPreconditions(RequestingUser); !Exists || !Authorized {

		return nil, errors.New("Preconditions failed.")

	}

	return os.ReadFile(AssociatedFile.Path)

}

func (AssociatedFile *File) Delete(RequestingUser string) (error, error) {

	if Exists, Authorized := AssociatedFile.checkPreconditions(RequestingUser); !Exists || !Authorized {

		return errors.New("Preconditions failed."), nil

	}

	return os.Remove(AssociatedFile.Path), AssociatedFile.deleteInfo()

}

func (AssociatedFile *File) ToHTML() string {

	MimeType, Icon := AssociatedFile.getMimeTypeAndIcon()

	HumanReadableType := Types.MimeTypeToReadableName[MimeType]
	HumanReadableSize := FileUtil.NormalizeSize(AssociatedFile.Size)

	return Functions.CleanEscapedString(fmt.Sprintf(`

		<div class="InlineFile">

			<div class="InlineFileIcon">

				<ion-icon name="%s"></ion-icon>

			</div>

			<div class="InlineFileDetails">

				<div class="InlineFileName">%s</div>

				<div class="ul InlineFileStats">

					<li class="InlineFileStat Size">%s</li>

					<li class="InlineFileStat Type">%s</li>

				</div>

		</div>

	`, Icon, AssociatedFile.Name, HumanReadableSize, HumanReadableType))

}
