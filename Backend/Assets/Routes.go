package Assets

import (
	"elucid503/SproutCDN/Middleware"

	// Route Imports from ./Routes

	APIRoutes "elucid503/SproutCDN/Routes/API"
	FileRoutes "elucid503/SproutCDN/Routes/API/Files"
	MiscRoutes "elucid503/SproutCDN/Routes/API/Misc"

	"github.com/gin-gonic/gin"
)

type RouteInfo struct {
	Method  string
	Handler func(*gin.Context)

	// Middleware Opts

	RateLimitConfig *Middleware.RouteRateLimitConfig
	Authorized      bool
}

func GetRoutes() map[string]RouteInfo {

	return map[string]RouteInfo{

		"/API/": {

			Method:  "GET",
			Handler: APIRoutes.Index,

			RateLimitConfig: nil,
		},

		"/API/Files": {

			Method:  "GET",
			Handler: FileRoutes.AllFiles,

			RateLimitConfig: &Middleware.RouteRateLimitConfig{

				MaxRequestsAllowed: 120,
				TimeWindow:         60,
			},

			Authorized: true,
		},

		"/API/Files/New": {

			Method:  "POST",
			Handler: FileRoutes.NewFile,

			RateLimitConfig: &Middleware.RouteRateLimitConfig{

				MaxRequestsAllowed: 60,
				TimeWindow:         60,
			},

			Authorized: false,
		},

		"/API/Files/:uid/Info": {

			Method:  "GET",
			Handler: FileRoutes.FileInfo,

			RateLimitConfig: &Middleware.RouteRateLimitConfig{

				MaxRequestsAllowed: 240,
				TimeWindow:         60,
			},

			Authorized: true,
		},

		"/API/Files/:uid/Delete": {

			Method:  "DELETE",
			Handler: FileRoutes.DeleteFile,

			RateLimitConfig: &Middleware.RouteRateLimitConfig{

				MaxRequestsAllowed: 120,
				TimeWindow:         60,
			},

			Authorized: true,
		},

		"/API/Misc/Auth": {

			Method:  "POST",
			Handler: MiscRoutes.Authorize,

			RateLimitConfig: &Middleware.RouteRateLimitConfig{

				MaxRequestsAllowed: 90,
				TimeWindow:         60,
			},
		},
	}

}
