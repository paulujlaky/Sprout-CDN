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

	Private bool `json:"Private"`

	Authorized []string `json:"Authorized"`

	SubDirLength  int `json:"SubDirLength"`
	SubFileLength int `json:"SubFileLength"`
}

// General

func NewDir(Path string, OriginalPath string, Authorized []string, Private bool) (*Dir, error, error) {

	NormalizedPath := strings.TrimPrefix(OriginalPath, string(filepath.Separator))

	DirInstance := &Dir{

		UID: Functions.RandomString(16),

		Path:           Path,
		NormalizedPath: NormalizedPath,

		Name: filepath.Base(Path),

		Private: Private,

		Authorized: Authorized,
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

	return filepath.Join(filepath.Dir(DirPath), fmt.Sprintf("%s.dirinfo", filepath.Base(DirPath)))

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

// Public

func (AssociatedDir *Dir) Create(RequestingUser string) error {

	if UserAuthed, Exists := AssociatedDir.checkPreconditions(RequestingUser, AssociatedDir.Path); !UserAuthed || !Exists {

		return errors.New("Preconditions failed")

	}

	return DirUtil.Create(AssociatedDir.Path)

}

func (AssociatedDir *Dir) Delete(RequestingUser string) (error, error) {

	if UserAuthed, Exists := AssociatedDir.checkPreconditions(RequestingUser, AssociatedDir.Path); !UserAuthed || !Exists {

		return errors.New("Preconditions failed"), nil

	}

	return os.RemoveAll(AssociatedDir.Path), AssociatedDir.deleteInfo()

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

		<div class="Container HorizontalFlex InlineFile Dir" UID="%s">

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
