package Models

import (
	"elucid503/SproutCDN/Functions"
	"encoding/json"
	"strings"

	"errors"
	"os"
	"slices"
)

var FileUtil = Functions.FileUtil{} // init FileUtil for use in File methods

type File struct {
	UID  string `json:"UID"`
	Name string `json:"Name"`

	Authorized []string `json:"Authorized"`

	Path string `json:"Path"`
}

// General

func NewFile(Name string, Authorized []string, Path string) *File {

	return &File{

		UID:  Functions.RandomString(16),
		Name: Name,

		Authorized: Authorized,

		Path: Path,
	}

}

func getInfoPath(OriginalPath string) string {

	// We must remove the existing extension and add .info to the end

	LastPeriodIndex := strings.LastIndex(OriginalPath, ".")

	// Get the path without the extension

	ExtensionlessPath := OriginalPath[:LastPeriodIndex]

	// Add .info to the end

	return ExtensionlessPath + ".info"

}

func LoadFileFromDotInfo(Path string) error {

	// Load the file from the .info file

	InfoPath := getInfoPath(Path)

	InfoData, ReadError := os.ReadFile(InfoPath)

	if ReadError != nil {

		return ReadError

	}

	var FileData File // new instance

	FileData.Path = Path

	return json.Unmarshal(InfoData, &FileData) // loads all the data from the .info file into the struct

}

// Private

func (AssociatedFile *File) resolvePath() (string, error) {

	return FileUtil.ResolvePath(AssociatedFile.Path)

}

func (AssociatedFile *File) checkPreconditions(RequestingUser string) (bool, bool) {

	return slices.Contains(AssociatedFile.Authorized, RequestingUser), FileUtil.FileExists(AssociatedFile.Path)

}

func (AssociatedFile *File) writeInfo() error {

	// Write all info in the struct to a .info file which is later parsed as JSON

	DataToWrite, MarshalError := json.Marshal(AssociatedFile)

	if MarshalError != nil {

		return MarshalError

	}

	return os.WriteFile(getInfoPath(AssociatedFile.Path), DataToWrite, os.ModePerm)

}

// Public

func (AssociatedFile *File) DiskInfo(RequestingUser string) (os.FileInfo, error) {

	if Exists, Authorized := AssociatedFile.checkPreconditions(RequestingUser); !Exists || !Authorized {

		return nil, errors.New("Preconditions failed.")

	}

	return os.Stat(AssociatedFile.Path)

}

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

func (AssociatedFile *File) Delete(RequestingUser string) error {

	if Exists, Authorized := AssociatedFile.checkPreconditions(RequestingUser); !Exists || !Authorized {

		return errors.New("Preconditions failed.")

	}

	return os.Remove(AssociatedFile.Path)

}
