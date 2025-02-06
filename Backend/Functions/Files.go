package Functions

import (
	"os"
	"path/filepath"
)

type FileUtil struct{}

func (f *FileUtil) FileExists(Path string) bool {

	_, Err := os.Stat(Path)

	return !os.IsNotExist(Err)

}

func (f *FileUtil) DirectoryExists(Path string) bool {

	Info, Err := os.Stat(Path)

	return !os.IsNotExist(Err) && Info.IsDir()

}

func (f *FileUtil) ResolvePath(Path string) (string, error) {

	RelativePath, error := filepath.Rel(".", Path)

	if error != nil {

		return "", error

	}

	if RelativePath == "." {

		return "", error

	}

	return RelativePath, nil

}