package Functions

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/skip2/go-qrcode"
)

// HTTP Requests

type RequestOptions struct {
	Headers map[string]string
}

var HTTPClient = &http.Client{}

func MakeHTTPRequest(Method string, URL string, Options RequestOptions) (*http.Response, error) {

	// Make the request

	RequestInstance, Err := http.NewRequest(Method, URL, nil)

	if Err != nil {

		return nil, Err

	}

	for HeaderKey, HeaderValue := range Options.Headers {

		RequestInstance.Header.Set(HeaderKey, HeaderValue)

	}

	RequestResponse, RequestErr := HTTPClient.Do(RequestInstance)

	if RequestErr != nil {

		return nil, RequestErr // Return the error

	}

	return RequestResponse, nil

}

func GetHTTPRequestJSONResponse(Response *http.Response) (map[string]interface{}, error) {

	var ResponseData map[string]interface{}

	Decoder := json.NewDecoder(Response.Body)

	if Err := Decoder.Decode(&ResponseData); Err != nil {

		return nil, Err

	}

	return ResponseData, nil

}

// QR Code

func GenerateQRCode(URL string) ([]byte, error) {

	// Generate the QR code

	QRCode, Err := qrcode.Encode(URL, qrcode.Medium, 256)

	if Err != nil {

		return nil, Err

	}

	return QRCode, nil

}

// Random

func RandomString(Length int) string {

	var Charset string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var RandomStr string

	for i := 0; i < Length; i++ {

		RandomStr += string(Charset[rand.Intn(len(Charset))])

	}

	return RandomStr

}

// Strings

func CleanEscapedString(InputStr string) string {

	Replacer := strings.NewReplacer("\n", "", "\t", "")

	return Replacer.Replace(InputStr)

}

func PluralizeString(Count int, Str string) string {

	if Count == 1 {

		return Str

	} else {

		return Str + "s"

	}

}

// Paths

func SanitizePath(Path string) string {

	// Check base path

	if !strings.Contains(filepath.Dir(Path), "Store") {

		Path = strings.ReplaceAll(Path, "..", "")

	}

	// Remove any trailing slashes

	if strings.HasSuffix(Path, "/") {

		Path = Path[:len(Path)-1]

	}

	return Path

}

func AdjustPathToStore(Path string) string {

	// Make path work for CDN dir

	return filepath.Join("../Store", Path)

}

func RemovePathFromStore(Path string) string {

	return strings.ReplaceAll(Path, "../Store", "")

}

func RemovePathDoubleSlashes(Path string) string {

	return strings.ReplaceAll(Path, fmt.Sprintf("%c%c", filepath.Separator, filepath.Separator), string(filepath.Separator))

}

func ReplacePathWithForwardSlashes(Path string) string {

	return strings.ReplaceAll(Path, string(filepath.Separator), "/")

}

func RemoveStoreFromPath(Path string) string {

	return strings.ReplaceAll(strings.ReplaceAll(Path, "../Store/", ""), "../Store", "")

}

func NormalizePath(Path string) string {

	return RemoveStoreFromPath(RemovePathDoubleSlashes(ReplacePathWithForwardSlashes(Path)))

}
