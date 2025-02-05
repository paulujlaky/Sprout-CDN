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

func (f *FileUtil) CreateDirectory(Path string) error {

	return os.MkdirAll(Path, os.ModePerm)

}

func (f *FileUtil) CreateFile(Path string) (*os.File, error) {

	return os.Create(Path)

}

func (f *FileUtil) DeleteFile(Path string) error {

	return os.Remove(Path)

}

func (f *FileUtil) DeleteDirectory(Path string) error {

	return os.RemoveAll(Path)

}

func (f *FileUtil) ReadFile(Path string) ([]byte, error) {

	return os.ReadFile(Path)

}

func (f *FileUtil) WriteFile(Path string, Data []byte) error {

	return os.WriteFile(Path, Data, os.ModePerm)

}
