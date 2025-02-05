package main

import (
	"elucid503/SproutCDN/Middleware"
	"elucid503/SproutCDN/Routes"

	"github.com/gin-gonic/gin"
)

type RouteInfo struct {
	Method  string
	Handler func(*gin.Context)

	RateLimitConfig *Middleware.RouteRateLimitConfig
}

func GetRoutes() map[string]RouteInfo {

	return map[string]RouteInfo{

		"/": {

			"GET", Routes.Index,

			&Middleware.RouteRateLimitConfig{

				MaxRequestsAllowed: 10,
				TimeWindow:         60,
			},
		},
	}

}
