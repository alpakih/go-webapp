package auth

import (
	"errors"
	"github.com/alpakih/go-webapp/internal"
	"github.com/alpakih/go-webapp/internal/user"
	"github.com/alpakih/go-webapp/pkg/database"
	"github.com/alpakih/go-webapp/pkg/helper"
	"github.com/alpakih/go-webapp/pkg/sessions"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

type Handler struct {
	internal.BaseHandler
}

func NewAuthController(session *sessions.Manager) Handler {
	return Handler{
		BaseHandler: internal.BaseHandler{
			Title: "Authentication | Login",
			Menu:  "Authentication",
			PageHeader: []map[string]interface{}{
				{
					"menu": "Authentication",
					"link": "/admin/auth/login",
				},
			},
			Session: session,
		},
	}
}

func (r *Handler) LoginPage(ctx echo.Context) error {
	pageHeader := map[string]interface{}{
		"menu": "Auth",
		"link": "/admin/auth/login",
	}
	return r.Render(ctx, "auth/login", append(r.PageHeader, pageHeader), nil)
}
func (r *Handler) Login(ctx echo.Context) error {
	var loginDto LoginDto
	if err := ctx.Bind(&loginDto); err != nil {
		r.Session.SetFlashMessage(ctx, "error binding data", "error", nil)
		return ctx.Redirect(302, "/admin/auth/login")
	}
	if err := ctx.Validate(&loginDto); err != nil {
		r.Session.SetFlashMessage(ctx, "validation Error", "error", nil)
		return ctx.Redirect(302, "/admin/auth/login")
	}
	var auth user.User
	if err := database.Conn().Preload("Role.Permission").First(&auth, "email =?", loginDto.Email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Session.SetFlashMessage(ctx, "user with email "+loginDto.Email+" is not found", "error", nil)
			return ctx.Redirect(302, "/admin/auth/login")
		}
		r.Session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.Redirect(302, "/admin/auth/login")
	}

	if !helper.CheckPasswordHash(loginDto.Password, auth.Password) {
		r.Session.SetFlashMessage(ctx, "wrong email or password", "error", nil)
		return ctx.Redirect(302, "/admin/auth/login")
	}
	var permission []string
	for _, v := range auth.Role.Permission {
		permission = append(permission, v.Url)
	}
	if err := r.Session.Set(ctx, sessions.IDSession, r.Session.
		SetUserInfo(
			auth.ID, auth.Username, auth.Email, auth.Role.Slug, auth.ImageUrl, permission,
		)); err != nil {
		r.Session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.Redirect(302, "/admin/auth/login")
	}
	r.Session.SetFlashMessage(ctx, "login success", "success", echo.Map{"Email": auth.Email})
	return ctx.Redirect(302, "/admin/home")
}

func (r *Handler) Logout(ctx echo.Context) error {
	err := r.Session.Delete(ctx, sessions.IDSession)
	if err != nil {
		r.Session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.Redirect(302, "/admin/home")
	}
	r.Session.SetFlashMessage(ctx, "logout success", "success", nil)
	return ctx.Redirect(http.StatusFound, "/admin/auth/login")
}
