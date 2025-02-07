package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	"elucid503/SproutCDN/Assets"

	"elucid503/SproutCDN/Middleware"
	"elucid503/SproutCDN/Types"
)

var ConfigFilePath string = "./Assets/Config.json"

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

var RouteHandlers = Assets.GetRoutes()

func LoadRoutes(Router *gin.Engine) {

	for Route, Handler := range RouteHandlers {

		MiddlewareToInclude := []gin.HandlerFunc{}

		if Handler.Authorized {

			MiddlewareToInclude = append(MiddlewareToInclude, Middleware.Authorize())

		}

		if Handler.RateLimitConfig != nil {

			MiddlewareToInclude = append(MiddlewareToInclude, Middleware.RateLimit(Handler.RateLimitConfig))

		}

		MiddlewareToInclude = append(MiddlewareToInclude, Handler.Handler) // Finally add the real handler

		Router.Handle(Handler.Method, Route, MiddlewareToInclude...)

	}

}

func main() {

	Config, Error := LoadConfig()

	if Error != nil {

		panic(Error) // Can't continue

	}

	fmt.Printf("Sprout CDN Backend; Version %s\n", Config.Versions.Server)

	gin.SetMode(Config.GinMode)

	GinRouter := gin.Default()

	LoadRoutes(GinRouter)

	fmt.Printf("Loaded %d routes...\n", len(RouteHandlers))

	fmt.Printf("Listening on port %d...\n", Config.Server.Port)

	GinRouter.Run(fmt.Sprintf(":%d", Config.Server.Port))

}
