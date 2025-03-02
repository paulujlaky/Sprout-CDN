package Models

import (
	"elucid503/SproutCDN/Functions"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type Dir struct {
	UID string `json:"UID"`

	Path           string `json:"Path"`
	NormalizedPath string `json:"NormalizedPath"` // Path without the store prefix

	Name string `json:"Name"`
	URL  string `json:"URL"`

	Private bool `json:"Private"`

	Authorized []string `json:"Authorized"`

	SubDirLength  int `json:"SubDirLength"`
	SubFileLength int `json:"SubFileLength"`
}

// General

func NewDir(Path string, OriginalPath string, Authorized []string, Private bool) (*Dir, error, error) {

	DirInstance := &Dir{

		UID: Functions.RandomString(16),

		Path:           Path,
		NormalizedPath: Functions.NormalizePath(OriginalPath),

		Name: filepath.Base(Path),
		URL:  Domain + "/Dash/" + Functions.NormalizePath(OriginalPath),

		Private: Private,

		Authorized: Authorized,
	}

	CreateError := DirUtil.Create(Path)
	InfoWriteError := DirInstance.writeInfo()

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

	return filepath.Join(DirPath, fmt.Sprintf("%s.dirinfo", filepath.Base(DirPath)))

}

// Private

func (AssociatedDir *Dir) checkPreconditions(RequestingUser string, DirPath string) (bool, bool) {

	PassesPrivateCheck := slices.Contains(AssociatedDir.Authorized, RequestingUser)
	PassesExistCheck := DirUtil.Exists(filepath.Dir(DirPath))

	if AssociatedDir.Private {

		return PassesPrivateCheck, PassesExistCheck

	} else {

		return true, PassesExistCheck

	}

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

func (AssociatedDir *Dir) getContentInfo() (int, int) {

	// Get contents

	Contents, _ := DirUtil.GetContents(AssociatedDir.Path)

	// Count files and directories

	SubDirLength := 0
	SubFileLength := 0

	for _, Content := range Contents {

		// skip .fileinfo and .dirinfo files

		if strings.HasSuffix(Content.Name, ".fileinfo") || strings.HasSuffix(Content.Name, ".dirinfo") {

			continue

		}

		ContentPath := filepath.Join(AssociatedDir.Path, Content.Name)

		Info, _ := FileUtil.GetInfo(ContentPath)

		if Info.IsDir() {

			SubDirLength++

		} else {

			SubFileLength++

		}

	}

	return SubDirLength, SubFileLength

}

func (AssociatedDir *Dir) moveContents(NewPath string) {

	// TODO: Improve error handling here, maybe

	// Get ALL contents of AssociatedDir and move them to NewPath

	InitialContents, _ := DirUtil.GetContents(AssociatedDir.Path)

	for _, Content := range InitialContents {

		// Skip .fileinfo and .dirinfo files

		if strings.HasSuffix(Content.Name, ".fileinfo") || strings.HasSuffix(Content.Name, ".dirinfo") {

			continue

		}

		ContentPath := filepath.Join(AssociatedDir.Path, Content.Name)

		NewContentPath := filepath.Join(NewPath, Content.Name)

		if Content.IsDir {

			// Move the directory

			OriginalDir, _ := LoadDirFromDotInfo(ContentPath)

			NewDir(NewContentPath, Functions.NormalizePath(NewContentPath), OriginalDir.Authorized, OriginalDir.Private)

			fmt.Println("Moving dir to: ", NewContentPath)

			OriginalDir.moveContents(NewContentPath) // Performs a recursive move

		} else {

			// Move the file

			OriginalFile, _ := LoadFileFromDotInfo(ContentPath)

			NewFileInstance := NewFile(OriginalFile.Name, OriginalFile.Size, OriginalFile.Private, OriginalFile.Authorized, NewPath)

			NewFileInstance.writeInfo()

			OriginalContents, _ := OriginalFile.Read(OriginalFile.Authorized[0])

			NewFileInstance.Write(OriginalFile.Authorized[0], OriginalContents)

		}

	}

}

// Public

func (AssociatedDir *Dir) Create(RequestingUser string) error {

	if UserAuthed, Exists := AssociatedDir.checkPreconditions(RequestingUser, AssociatedDir.Path); !UserAuthed || !Exists {

		return errors.New("Preconditions failed")

	}

	return DirUtil.Create(AssociatedDir.Path)

}

func (AssociatedDir *Dir) Move(RequestingUser string, NewPath string) error {

	if UserAuthed, Exists := AssociatedDir.checkPreconditions(RequestingUser, AssociatedDir.Path); !UserAuthed || !Exists {

		return errors.New("Preconditions failed")

	}

	// Get the new path

	JoinedNewPath := filepath.Join(NewPath, AssociatedDir.Name)

	NewPath = Functions.AdjustPathToStore(JoinedNewPath)

	// Extra Preconditions: Check if there is a dir with the same name in the new path

	if DirUtil.Exists(NewPath) {

		return errors.New("Directory with the same name already exists in the new path")

	}

	// Create the new directory

	CreateError := DirUtil.Create(NewPath)

	if CreateError != nil {

		return errors.New("Failed to create new directory")

	}

	// Copy everything

	AssociatedDir.moveContents(NewPath)

	// Delete the old dir

	MoveError := os.RemoveAll(AssociatedDir.Path)

	// Delete the old info

	AssociatedDir.deleteInfo()

	// Update the path in the struct

	AssociatedDir.Path = NewPath
	AssociatedDir.NormalizedPath = Functions.NormalizePath(JoinedNewPath)
	WriteInfoError := AssociatedDir.writeInfo()

	if MoveError != nil || WriteInfoError != nil {

		return errors.New("Failed to move directory")

	}

	return nil

}

func (AssociatedDir *Dir) Delete(RequestingUser string) error {

	if UserAuthed, Exists := AssociatedDir.checkPreconditions(RequestingUser, AssociatedDir.Path); !UserAuthed || !Exists {

		return errors.New("Preconditions failed")

	}

	return os.RemoveAll(AssociatedDir.Path)

}

func (AssociatedDir *Dir) GetContents(RequestingUser string) ([]Functions.DirContentItem, error) {

	if UserAuthed, Exists := AssociatedDir.checkPreconditions(RequestingUser, AssociatedDir.Path); !UserAuthed || !Exists {

		return nil, errors.New("Preconditions failed")

	}

	return DirUtil.GetContents(AssociatedDir.Path)

}

func (AssociatedDir *Dir) ToHTML() string {

	CurrentSubDirs, CurrentSubFiles := AssociatedDir.getContentInfo()

	PrivateIndicatorVisibility := "none"

	if AssociatedDir.Private {

		PrivateIndicatorVisibility = "block"

	}

	return Functions.CleanEscapedString(fmt.Sprintf(`

		<div class="Container HorizontalFlex InlineFile Dir" UID="%s" draggable="true">

			<div class="InlineFileContent Left"> 

				<div class="InlineFileIcon">

					<ion-icon name="folder-outline"></ion-icon>

				</div>

				<div class="Container Transparent VerticalFlex InlineFileDetails">

					<div class="InlineFileName">%s</div>

					<ul class="InlineFileStats">

						<li class="InlineFileStat SubFiles">%d %s</li>

						<li class="InlineFileStat SubDirs">%d %s</li>

						<li class="InlineFileStat PrivateIndicator" style="display:%s">Private</li>

					</ul>
				
				</div>

			</div>

			<div class="InlineFileActions Right">

				<ion-icon name="ellipsis-horizontal"></ion-icon>

			</div>

		</div>

	`, AssociatedDir.UID, AssociatedDir.Name, CurrentSubFiles, Functions.PluralizeString(CurrentSubFiles, "File"), CurrentSubDirs, Functions.PluralizeString(CurrentSubDirs, "Folder"), PrivateIndicatorVisibility))

}
