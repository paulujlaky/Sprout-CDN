package Functions

import (
	"elucid503/SproutCDN/Types"
	"fmt"
	"os"
	"path/filepath"
)

type FileUtil struct{}

func (f *FileUtil) FileExists(Path string) bool {

	_, Err := os.Stat(Path)

	return !os.IsNotExist(Err)

}

func (f *FileUtil) GetInfo(Path string) (os.FileInfo, error) {

	return os.Stat(Path)

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

func (f *FileUtil) GetMimeTypeAndIcon(Extension string) (string, string) {

	MimeType, Exists := Types.ExtensionsToMimeType[Extension]

	if !Exists {

		// Return defaults

		return "application/octet-stream", "document-outline"

	}

	Icon := Types.MimeTypeToIcon[MimeType] // must exist since all MimeType's are in the map

	return MimeType, Icon

}

func (f *FileUtil) NormalizeSize(ByteLength int64) string {

	// Decide whether to use bytes, kilobytes, megabytes or gigabytes

	var Breakpoints = []int64{1024, 1024 * 1024, 1024 * 1024 * 1024} // 1KB, 1MB, 1GB
	var BreakpointNames = []string{"Bytes", "KB", "MB", "GB"}

	var BreakpointIndex int = 0 // May keep being 0, in that case bytes work fine

	for ByteLength > Breakpoints[BreakpointIndex] {

		// Go up by one breakpoint

		ByteLength /= 1024
		BreakpointIndex++

	}

	return fmt.Sprintf("%d %s", ByteLength, BreakpointNames[BreakpointIndex])

}
