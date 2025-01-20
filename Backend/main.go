package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"elucid503/SproutCDN/Routes"
	"elucid503/SproutCDN/Types"
)

var ConfigFilePath string = "./Assets/Config.json"

var RouteHandlers = map[string]func(*gin.Context){

	"/": Routes.Index,
}

func LoadRoutes(Router *gin.Engine, RootDir string) error {

	return filepath.Walk(RootDir, func(Path string, Info os.FileInfo, Err error) error {

		if Err != nil {
			return Err
		}

		// Skip directories

		if Info.IsDir() {
			return nil
		}

		RelativePath, _ := filepath.Rel(RootDir, Path)
		Route := "/" + strings.TrimSuffix(strings.ReplaceAll(RelativePath, string(os.PathSeparator), "/"), ".go")

		// Replace certain names (index) with the root route

		if Route == "/index" {
			Route = "/"
		}

		// Registering the handler

		if Handler, Exists := RouteHandlers[Route]; Exists {

			Router.GET(Route, Handler) // TODO: Add POST, PUT, DELETE, etc.

		}

		return nil

	})

}

func LoadConfig() (*Types.Config, error) {

	// Read the file

	File, Err := os.Open(ConfigFilePath)

	if Err != nil {

		return nil, Err

	}

	defer File.Close()

	// Decode the file

	Decoder := json.NewDecoder(File)

	Config := &Types.Config{} // Get the pointer to the Config struct

	if Err := Decoder.Decode(Config); Err != nil {

		return nil, Err

	}

	return Config, nil

}

func main() {

	Config, Error := LoadConfig()

	if Error != nil {

		panic(Error) // Can't continue

	}

	fmt.Println("Loaded Config...")

	fmt.Printf("Version %s", Config.Versions.Server)

	gin.SetMode(gin.ReleaseMode)

	GinRouter := gin.Default()

	RoutesDir := "./Routes"

	if Err := LoadRoutes(GinRouter, RoutesDir); Err != nil {

		panic(Err)

	}

	fmt.Printf("Loaded %d routes...\n", len(RouteHandlers))

	fmt.Printf("Listening on port %d...\n", Config.Server.Port)

	GinRouter.Run(fmt.Sprintf(":%d", Config.Server.Port))

}
