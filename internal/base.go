package internal

import (
	"github.com/alpakih/go-webapp/pkg/sessions"
	"github.com/foolin/goview/supports/echoview-v4"
	"github.com/labstack/echo/v4"
	"net/http"
)

type BaseHandler struct {
	Title      string
	Menu       string
	Session    *sessions.Manager
	PageHeader []map[string]interface{}
}

func (r *BaseHandler) Render(ctx echo.Context, view string, pageHeader []map[string]interface{}, data interface{}) error {
	return echoview.Render(ctx, http.StatusOK, view, echo.Map{
		"title": r.Title, "menu": r.Menu, "pageHeader": pageHeader, "data": data,
		"flashMessage": r.Session.GetFlashMessage(ctx), "ctx": ctx,
	})
}
