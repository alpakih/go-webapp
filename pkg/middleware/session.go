package middleware

import (
	session "github.com/alpakih/go-webapp/pkg/sessions"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"strings"
)

type sessionConf struct {
	Skipper       middleware.Skipper
	AuthKey       string
	EncryptionKey string
}

func NewSessionMiddleware(authKey string, encryptionKey string) *sessionConf {
	return &sessionConf{
		Skipper: func(context echo.Context) bool {
			apiV1URI := "/admin"
			if strings.EqualFold(context.Request().RequestURI, apiV1URI+"/auth/login") {
				return true
			}
			return false
		},
		AuthKey:       authKey,
		EncryptionKey: encryptionKey,
	}
}

func (r *sessionConf) NewCookieStore() *sessions.CookieStore {
	authKey := []byte(r.AuthKey)
	encryptionKey := []byte(r.EncryptionKey)
	s := sessions.NewCookieStore(authKey, encryptionKey)
	s.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	return s
}

func (r *sessionConf) SessionMiddleware(s *session.Manager) echo.MiddlewareFunc {
	return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) error {
			if r.Skipper(context) {
				return handlerFunc(context)
			}
			result, err := s.Get(context, session.IDSession)
			if err != nil {
				return context.Redirect(302, "/admin/auth/login")
			}
			if result == nil {
				return context.Redirect(302, "/admin/auth/login")
			} else {
				return handlerFunc(context)
			}
		}
	}
}
