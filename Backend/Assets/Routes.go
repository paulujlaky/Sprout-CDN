package Assets

import (
	"elucid503/SproutCDN/Middleware"

	// Route Imports from ./Routes

	GlobalRoutes "elucid503/SproutCDN/Routes"
	APIRoutes "elucid503/SproutCDN/Routes/API"
	DirRoutes "elucid503/SproutCDN/Routes/API/Dirs"
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

		"/": {

			Method:  "GET",
			Handler: GlobalRoutes.Dash,

			RateLimitConfig: nil,
		},

		"/Dash/*Route": {

			Method:  "GET",
			Handler: GlobalRoutes.Dash,

			RateLimitConfig: nil,
		},

		"/Files/*File": {

			Method:  "GET",
			Handler: GlobalRoutes.GetFile,

			RateLimitConfig: &Middleware.RouteRateLimitConfig{

				MaxRequestsAllowed: 400,

				TimeWindow: 60,
			},

			Authorized: false,
		},

		"/API/": {

			Method:  "GET",
			Handler: APIRoutes.Index,

			RateLimitConfig: nil,
		},

		"/API/Files/*Path": {

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

			Authorized: true,
		},

		"/API/Files/Delete": {

			Method:  "DELETE",
			Handler: FileRoutes.DeleteFile,

			RateLimitConfig: &Middleware.RouteRateLimitConfig{

				MaxRequestsAllowed: 120,
				TimeWindow:         60,
			},

			Authorized: true,
		},

		"/API/Dirs/New": {

			Method:  "POST",
			Handler: DirRoutes.CreateDir,

			RateLimitConfig: &Middleware.RouteRateLimitConfig{

				MaxRequestsAllowed: 60,
				TimeWindow:         60,
			},

			Authorized: true,
		},

		"/API/Dirs/Delete": {

			Method:  "DELETE",
			Handler: DirRoutes.DeleteDir,

			RateLimitConfig: &Middleware.RouteRateLimitConfig{

				MaxRequestsAllowed: 60,
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
