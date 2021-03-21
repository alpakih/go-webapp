package home

import (
	"github.com/alpakih/go-webapp/internal"
	"github.com/alpakih/go-webapp/pkg/sessions"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	internal.BaseHandler
}

func NewHomeController(session *sessions.Manager) Handler {
	return Handler{
		BaseHandler: internal.BaseHandler{
			Title: "Home",
			Menu:  "Home",
			PageHeader: []map[string]interface{}{
				{
					"menu": "Home",
					"link": "/admin/home",
				},
			},
			Session: session,
		},
	}
}

func (c *Handler) Index(ctx echo.Context) error {
	return c.Render(ctx, "home/index", c.PageHeader, nil)
}
