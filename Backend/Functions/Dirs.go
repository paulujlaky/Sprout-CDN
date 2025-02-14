package Functions

import (
	"os"
)

type DirUtil struct{}

type DirContentItem struct {
	Name  string
	IsDir bool
}

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

func (d *DirUtil) GetContents(Path string) ([]DirContentItem, error) {

	Info, Err := os.ReadDir(Path)

	if Err != nil {

		return nil, Err

	}

	var Contents []DirContentItem

	for _, DirEntry := range Info {

		Contents = append(Contents, DirContentItem{

			Name:  DirEntry.Name(),
			IsDir: DirEntry.IsDir(),

		})

	}

	return Contents, nil

}
