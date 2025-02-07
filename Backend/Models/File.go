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

func (AssociatedFile *File) checkPreconditions(RequestingUser string) (bool, bool) {

	return slices.Contains(AssociatedFile.Authorized, RequestingUser), FileUtil.FileExists(AssociatedFile.Path)

}

// Public

func (AssociatedFile *File) Write(RequestingUser string, Data []byte) error {

	if Exists, Authorized := AssociatedFile.checkPreconditions(RequestingUser); !Exists || !Authorized {

		return errors.New("Preconditions failed.")

	}

	return os.WriteFile(AssociatedFile.Path, Data, os.ModePerm)

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
