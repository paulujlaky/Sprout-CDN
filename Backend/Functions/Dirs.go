package Functions

import (
	"os"
)

type DirUtil struct{}

func (d *DirUtil) Exists(Path string) bool {

	Info, Err := os.Stat(Path)

	return !os.IsNotExist(Err) && Info.IsDir()

}

func (d *DirUtil) Create(Path string) error {

	return os.MkdirAll(Path, os.ModePerm)

}

func (d *DirUtil) Delete(Path string) error {

	return os.RemoveAll(Path)

}

func (d *DirUtil) GetContents(Path string) ([]string, error) {

	Info, Err := os.ReadDir(Path)

	if Err != nil {

		return nil, Err

	}

	var Contents []string

	for _, DirEntry := range Info {

		Contents = append(Contents, DirEntry.Name())

	}

	return Contents, nil

}
