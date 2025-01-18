package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"elucid503/SproutCDN/Routes"
)

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

			Router.GET(Route, Handler)
			fmt.Printf("Registered route: %s\n", Route)

		}

		return nil

	})

}

func main() {

	gin.SetMode(gin.ReleaseMode)

	GinRouter := gin.Default()

	RoutesDir := "./Routes"

	if Err := LoadRoutes(GinRouter, RoutesDir); Err != nil {

		panic(Err)

	}

	GinRouter.Run(":50300")

}
