package Functions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/elucid503/Sprout-API-Go/Logs"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
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

type BufferWriteCloser struct {
	buf *bytes.Buffer
}

func (b *BufferWriteCloser) Write(p []byte) (n int, err error) {
	return b.buf.Write(p)
}

func (b *BufferWriteCloser) Close() error {
	// bytes.Buffer doesn't need closing, so this is a no-op
	return nil
}

func GenerateQRCode(URL string) ([]byte, error) {

	QRCodeInstance, ErrorLoadingQr := qrcode.New(URL)

	if ErrorLoadingQr != nil {

		return nil, ErrorLoadingQr

	}

	// QR code options

	Options := []standard.ImageOption{

		standard.WithFgColorRGBHex("#f2f2f2"),
		standard.WithBgColorRGBHex("#161616"),
	}

	// Create buffer to store QR code

	Buffer := bytes.NewBuffer(nil)
	BufferToWrite := &BufferWriteCloser{buf: Buffer}

	// Create writer with buffer and options

	Writer := standard.NewWithWriter(BufferToWrite, Options...)

	// Save QR code to buffer

	if err := QRCodeInstance.Save(Writer); err != nil {

		return nil, err

	}

	return Buffer.Bytes(), nil

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

// Logging 

func Log(Level Logs.LogLevel, Title string, Message string) {
	
	fmt.Printf("[%s] %s: %s\n", Logs.GetLogLevelString(Level), Title, Message)

	// Wrapper just for the service UID not to need to be used multiple times

    Logs.Log("WhUUH4xhsWCKhCs6", Level, Title, Message)

}