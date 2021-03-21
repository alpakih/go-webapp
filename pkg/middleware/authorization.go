package middleware

import (
	"github.com/alpakih/go-webapp/pkg/helper"
	"github.com/alpakih/go-webapp/pkg/sessions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"strings"
)

type authorizationConfig struct {
	session *sessions.Manager
	Skipper middleware.Skipper
}

func NewAuthorizationMiddleware(session *sessions.Manager) *authorizationConfig {
	return &authorizationConfig{
		session: session,
		Skipper: func(context echo.Context) bool {
			return false
		},
	}
}

func (m *authorizationConfig) AuthorizationMiddleware(roles []string) echo.MiddlewareFunc {
	return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) error {
			if m.Skipper(context) {
				return handlerFunc(context)
			}

			//check role
			var userAuth sessions.UserInfo

			session, err := m.session.Get(context, sessions.IDSession)
			if err != nil {
				return context.Render(http.StatusOK, "403.html",
					echo.Map{"title": "Forbidden", "message": "you must login to access this resource"})
			}
			if userSession, ok := session.(sessions.UserInfo); !ok {
				return context.Render(http.StatusOK, "403.html",
					echo.Map{"title": "Forbidden", "message": "you must login to access this resource"})
			} else {
				userAuth = userSession
			}

			//pass for super admin
			if strings.EqualFold(userAuth.RoleSlug, "super-admin") {
				return handlerFunc(context)
			}

			//check role user
			if !helper.ItemExists(roles, userAuth.RoleSlug) {
				return context.Render(http.StatusOK, "403.html",
					echo.Map{"title": "Forbidden", "message": "this user role don't have permission to access this resource"})
			}

			//check user permission for this url
			if len(userAuth.Permission) != 0 {
				for _, v := range userAuth.Permission {
					if strings.EqualFold(context.Path(), v) {
						return handlerFunc(context)
					}
				}

			} else {
				return context.Render(http.StatusOK, "403.html",
					echo.Map{"title": "Forbidden", "message": "permission is not set for this user"})
			}
			return context.Render(http.StatusOK, "403.html",
				echo.Map{"title": "Forbidden", "message": "you dont have permission to access this resource"})
		}
	}
}
