package middleware

import (
	"github.com/alpakih/go-webapp/pkg/helper"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"net/http"
	"strings"
)

type routeNotFoundConfig struct {
	Skipper middleware.Skipper
}

func NewRouteNotFoundMiddleware() *routeNotFoundConfig {
	return &routeNotFoundConfig{
		Skipper: func(context echo.Context) bool {
			apiV1URI := "/admin"
			if strings.EqualFold(context.Request().RequestURI, apiV1URI+"/auth/login") {
				return true
			}
			return false
		},
	}
}

func (m *routeNotFoundConfig) RouteNotFoundMiddleware(route []*echo.Route) echo.MiddlewareFunc {
	return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) error {
			if m.Skipper(context) {
				return handlerFunc(context)
			}

			var check []string
			for _, v := range route {
				check = append(check, v.Path)
			}
			log.Info(check)


			if !helper.ItemExists(check, context.Path()) {
				return context.Redirect(http.StatusFound,"/admin/404")
			}

			return handlerFunc(context)
		}
	}
}
