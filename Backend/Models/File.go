package Models

import (
	"elucid503/SproutCDN/Functions"

	"errors"
	"os"
	"slices"
)

var FileUtil = Functions.FileUtil{}

type File struct {
	Authorized []string `json:"Authorized"`

	Path string `json:"Path"`
	Type string `json:"Type"`
}

// Private

func (AssociatedFile *File) resolvePath() (string, error) {

	return FileUtil.ResolvePath(AssociatedFile.Path)

}

func (AssociatedFile *File) checkPreconditions(UserAction string) (bool, bool) {

	return slices.Contains(AssociatedFile.Authorized, UserAction), FileUtil.FileExists(AssociatedFile.Path)

}

// Public

func (AssociatedFile *File) Write(UserID string, Data []byte) error {

	if Exists, Authorized := AssociatedFile.checkPreconditions(UserID); !Exists || !Authorized {

		return errors.New("Preconditions failed.")

	}

	return os.WriteFile(AssociatedFile.Path, Data, os.ModePerm)

}

func (AssociatedFile *File) Read(UserID string) ([]byte, error) {

	if Exists, Authorized := AssociatedFile.checkPreconditions(UserID); !Exists || !Authorized {

		return nil, errors.New("Preconditions failed.")

	}

	return os.ReadFile(AssociatedFile.Path)

}
