package Models

import (
	"elucid503/SproutCDN/Functions"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

type Dir struct {
	Path string `json:"Path"`
	Name string `json:"Name"`

	Authorized []string `json:"Authorized"`

	SubDirLength  int `json:"SubDirLength"`
	SubFileLength int `json:"SubFileLength"`
}

// General

func NewDir(Path string, Authorized []string) (*Dir, error, error) {

	// Get contents

	Contents, _ := DirUtil.GetContents(Path)

	// Count files and directories

	SubDirLength := 0
	SubFileLength := 0

	for _, Content := range Contents {

		ContentPath := filepath.Join(Path, Content)

		Info, _ := FileUtil.GetInfo(ContentPath)

		if Info.IsDir() {

			SubDirLength++

		} else {

			SubFileLength++

		}

	}

	DirInstance := &Dir{

		Path: Path,
		Name: filepath.Base(Path),

		Authorized: Authorized,

		SubDirLength:  SubDirLength,
		SubFileLength: SubFileLength,
	}

	// Write info

	InfoWriteError := DirInstance.writeInfo()
	CreateError := DirUtil.Create(Path)

	return DirInstance, InfoWriteError, CreateError

}

func LoadDirFromDotInfo(Path string) (*Dir, error) {

	Path = Functions.AdjustPathToStore(Path)

	// Load the file from the .info file

	InfoPath := GetDirInfoPath(Path)

	InfoData, ReadError := os.ReadFile(InfoPath)

	if ReadError != nil {

		return nil, ReadError

	}

	var DirData Dir // new instance

	DirData.Path = Path

	UnmarshalError := json.Unmarshal(InfoData, &DirData)

	return &DirData, UnmarshalError

}

func GetDirInfoPath(DirPath string) string {

	return filepath.Join(DirPath, fmt.Sprintf("%s.info", filepath.Dir(DirPath)))

}

// Private

func (AssociatedDir *Dir) checkPreconditions(RequestingUser string, DirPath string) (bool, bool) {

	// Preconditions for most actions (deletion, modification, etc...)

	return slices.Contains(AssociatedDir.Authorized, RequestingUser), DirUtil.Exists(filepath.Dir(DirPath))

}

func (AssociatedDir *Dir) writeInfo() error {

	// Write all info in the struct to a .info file which is later parsed as JSON

	DataToWrite, MarshalError := json.Marshal(AssociatedDir)

	if MarshalError != nil {

		return MarshalError

	}

	return os.WriteFile(GetDirInfoPath(AssociatedDir.Path), DataToWrite, os.ModePerm)

}

func (AssociatedDir *Dir) deleteInfo() error {

	return os.Remove(GetDirInfoPath(AssociatedDir.Path))

}

// Public

func (AssociatedDir *Dir) Create(RequestingUser string) error {

	if UserAuthed, Exists := AssociatedDir.checkPreconditions(RequestingUser, AssociatedDir.Path); UserAuthed && Exists {

		return errors.New("Preconditions failed")

	}

	return DirUtil.Create(AssociatedDir.Path)

}

func (AssociatedDir *Dir) Delete(RequestingUser string) (error, error) {

	if UserAuthed, Exists := AssociatedDir.checkPreconditions(RequestingUser, AssociatedDir.Path); UserAuthed && Exists {

		return errors.New("Preconditions failed"), nil

	}

	return os.RemoveAll(AssociatedDir.Path), AssociatedDir.deleteInfo()

}

func (AssociatedDir *Dir) GetContents(RequestingUser string) ([]string, error) {

	if UserAuthed, Exists := AssociatedDir.checkPreconditions(RequestingUser, AssociatedDir.Path); UserAuthed && Exists {

		return nil, errors.New("Preconditions failed")

	}

	return DirUtil.GetContents(AssociatedDir.Path)

}

func (AssociatedDir *Dir) ToHTML() string {

	return Functions.CleanEscapedString(fmt.Sprintf(`

		<div class="InlineFile Dir">

			<div class="InlineFileIcon">

				<ion-icon name="%s"></ion-icon>

			</div>

			<div class="InlineFileDetails">

				<div class="InlineFileName">%s</div>

				<div class="ul InlineFileStats">

					<li class="InlineFileStat SubFiles">%d %s</li>

					<li class="InlineFileStat SubDirs">%d %s</li>

				</div>

		</div>

	`, "folder-outline", AssociatedDir.Name, AssociatedDir.SubFileLength, Functions.PluralizeString(AssociatedDir.SubFileLength, "File"), AssociatedDir.SubDirLength, Functions.PluralizeString(AssociatedDir.SubDirLength, "Folder")))

}
