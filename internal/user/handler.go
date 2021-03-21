package user

import (
	"errors"
	"fmt"
	"github.com/alpakih/go-webapp/internal"
	"github.com/alpakih/go-webapp/internal/user/model"
	"github.com/alpakih/go-webapp/pkg/database"
	"github.com/alpakih/go-webapp/pkg/helper"
	"github.com/alpakih/go-webapp/pkg/helper/filesystem"
	"github.com/alpakih/go-webapp/pkg/sessions"
	"github.com/alpakih/go-webapp/pkg/validation"
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	internal.BaseHandler
}

func NewUserController(session *sessions.Manager) Handler {
	return Handler{
		BaseHandler: internal.BaseHandler{
			Title: "User Management | User",
			Menu:  "User Management",
			PageHeader: []map[string]interface{}{
				{
					"menu": "User Management",
					"link": "/admin/users/list",
				},
			},
			Session: session,
		},
	}
}

func (r *Handler) Index(ctx echo.Context) error {
	pageHeader := map[string]interface{}{
		"menu": "User",
		"link": "/admin/users/list",
	}
	return r.Render(ctx, "/user/index", append(r.PageHeader, pageHeader), nil)
}

func (r *Handler) Add(ctx echo.Context) error {

	pageHeader := []map[string]interface{}{{
		"menu": "User",
		"link": "/admin/users/list",
	},
		{
			"menu": "Add",
			"link": "/admin/users/add",
		},
	}

	return r.Render(ctx, "/user/add", append(r.PageHeader, pageHeader...), nil)
}

func (r *Handler) Datatable(ctx echo.Context) error {

	var qResult []model.ResultDatatableUser
	loc, _ := time.LoadLocation("Asia/Jakarta")
	draw := ctx.Request().URL.Query().Get("draw")
	search := ctx.Request().URL.Query().Get("search[value]")
	start := ctx.Request().URL.Query().Get("start")
	length := ctx.Request().URL.Query().Get("length")
	order := ctx.Request().URL.Query().Get("order[0][column]")
	orderName := ctx.Request().URL.Query().Get("columns[" + order + "][name]")
	orderAscDesc := ctx.Request().URL.Query().Get("order[0][dir]")

	var recordTotal int64
	var recordFiltered int64
	if err := database.Conn().Table("users").Count(&recordTotal).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	sql := `SELECT
				users.id,
				users.username,
				users.email,
				users.image_url,
				roles.role_name,
				users.created_at,
				COUNT (*) OVER () AS filter_row_count
			FROM
				users
				LEFT JOIN roles ON users.role_id = roles.id`

	if search != "" {
		sql += fmt.Sprintf(` WHERE users.email LIKE '%s'`, `%`+search+`%`)
		sql += fmt.Sprintf(` OR users.username LIKE '%s'`, `%`+search+`%`)
		sql += fmt.Sprintf(` OR roles.role_name LIKE '%s'`, `%`+search+`%`)
	}

	if orderName != "" && orderAscDesc != "" {
		sql += ` ORDER BY ` + orderName + ` ` + orderAscDesc
	} else {
		sql += ` ORDER BY users.created_at DESC`
	}

	parseLength, err := strconv.Atoi(length)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	if parseLength != -1 {
		sql += ` OFFSET ` + start + ` ROWS FETCH NEXT ` + length + ` ROWS ONLY`
	}
	if err := database.Conn().Raw(sql).Scan(&qResult).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	var action string
	listOfData := make([]map[string]interface{}, len(qResult))

	for k, v := range qResult {
		action = `<li style="list-style-type: none;display: inline">
					<a href="/admin/users/edit/` + v.ID + `" style="text-decoration: none;font-weight: 400; color: #333;"  data-toggle="tooltip" data-placement="right" title="Edit"><i class="fa fa-edit"></i>Edit </a>
				  </li>`
		action += `<li style="list-style-type: none;display: inline">
					<a href="/admin/users/view/` + v.ID + `" style="text-decoration: none;font-weight: 400; color: #333;"  data-toggle="tooltip" data-placement="right" title="View"><i class="fa fa-search"></i>View</a>
				  </li>`
		action += `<li style="list-style-type: none;display: inline">
					<a href="JavaScript:void(0);" onclick="Delete('` + v.ID + `')" style="text-decoration: none;font-weight: 400; color: #333;" data-toggle="tooltip" data-placement="right" title="Delete"><i class="fa fa-trash" style="color: #ff4d65"></i>Delete</a>
				  </li>`
		listOfData[k] = map[string]interface{}{
			"id":         v.ID,
			"username":   v.Username,
			"email":      v.Email,
			"image_url":  v.ImageUrl,
			"role_name":  v.RoleName,
			"created_at": v.CreatedAt.In(loc).Format("2006-01-02 15:04:05"),
			"action":     action,
		}
	}
	if len(qResult) != 0 {
		recordFiltered = *qResult[0].FilterRowCount
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"draw":            draw,
		"recordsTotal":    recordTotal,
		"recordsFiltered": recordFiltered,
		"data":            listOfData,
	})
}

func (r *Handler) Delete(ctx echo.Context) error {
	var user User
	id := ctx.Param("id")
	result, err := r.Session.Get(ctx, sessions.IDSession)
	if err != nil {
		return ctx.JSON(500, echo.Map{"message": "cannot get data user from context"})
	}
	userInfo := result.(sessions.UserInfo)

	if id == userInfo.ID {
		return ctx.JSON(400, echo.Map{"message": "cannot delete current active user"})
	}
	if err := database.Conn().First(&user, "id=?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.JSON(400, echo.Map{"message": "user not found"})
		}
		return ctx.JSON(500, echo.Map{"message": "error when trying delete data"})
	}
	dir, _ := os.Getwd()
	fileLocation := filepath.Join(dir, "web"+user.ImageUrl)

	if filesystem.FileExists(fileLocation) {
		filesystem.Delete(fileLocation)
	}
	if err := database.Conn().Delete(&User{}, "id=?", id).Error; err != nil {
		return ctx.JSON(500, echo.Map{"message": "error when trying delete data"})
	}

	return ctx.JSON(200, echo.Map{"message": "delete data has been success"})
}

func (r *Handler) Store(ctx echo.Context) error {
	var userDto StoreDto
	if err := ctx.Bind(&userDto); err != nil {
		r.Session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.Redirect(302, "/admin/users/add")
	}
	if err := ctx.Validate(&userDto); err != nil {
		var validationErrors []validation.ErrorValidation
		if errs, ok := err.(validator.ValidationErrors); ok {
			validationErrors = validation.WrapValidationErrors(errs)
		}
		r.Session.SetFlashMessage(ctx, err.Error(), "error", validationErrors)
		return ctx.Redirect(302, "/admin/users/add")
	}

	var entity User
	dir, _ := os.Getwd()
	file := filesystem.ParseMultipartFiles(ctx.Request(), "image")
	if len(file) != 0 {
		fileName := file[0].Save(filepath.Join(dir, "web/assets/upload/avatars"),
			strings.ReplaceAll(file[0].Header.Filename, " ", ""))
		entity.ImageUrl = "/assets/upload/avatars/" + fileName
	}
	hashPassword, _ := helper.HashPassword(userDto.Password)
	entity.Username = userDto.Username
	entity.Email = userDto.Email
	entity.Password = hashPassword
	entity.RoleID = userDto.Role
	if err := database.Conn().Create(&entity).Error; err != nil {
		r.Session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.Redirect(302, "/admin/users/add")
	}
	r.Session.SetFlashMessage(ctx, "save data success", "success", nil)
	return ctx.Redirect(http.StatusFound, "/admin/users/list")
}

func (r *Handler) Edit(ctx echo.Context) error {
	var user User
	id := ctx.Param("id")
	err := database.Conn().Joins("Role").First(&user, "users.id =?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Session.SetFlashMessage(ctx, err.Error(), "error", nil)
			return ctx.Redirect(302, "/admin/users/list")
		}
		r.Session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.Redirect(302, "/admin/users/list")
	}
	pageHeader := []map[string]interface{}{{
		"menu": "User",
		"link": "/admin/users/list",
	},
		{
			"menu": "Edit",
			"link": "/admin/users/edit/" + id,
		},
	}

	return r.Render(ctx, "/user/edit", append(r.PageHeader, pageHeader...), user)
}

func (r *Handler) Update(ctx echo.Context) error {
	var user User
	var updateDto UpdateDto
	id := ctx.Param("id")
	if err := ctx.Bind(&updateDto); err != nil {
		r.Session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.Redirect(302, "/admin/users/edit/"+id)
	}
	if err := ctx.Validate(&updateDto); err != nil {
		var validationErrors []validation.ErrorValidation
		if errs, ok := err.(validator.ValidationErrors); ok {
			validationErrors = validation.WrapValidationErrors(errs)
		}
		r.Session.SetFlashMessage(ctx, err.Error(), "error", validationErrors)
		return ctx.Redirect(302, "/admin/users/edit/"+id)
	}
	err := database.Conn().First(&user, "id =?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Session.SetFlashMessage(ctx, err.Error(), "error", nil)
			return ctx.Redirect(302, "/admin/users/list")
		}
		r.Session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.Redirect(302, "/admin/users/list")
	}
	file := filesystem.ParseMultipartFiles(ctx.Request(), "image")

	dir, _ := os.Getwd()
	if len(file) != 0 {
		if file[0].Header.Filename != "" {
			if user.ImageUrl != "" {
				fileLocation := filepath.Join(dir, "web"+user.ImageUrl)
				if filesystem.FileExists(fileLocation) {
					filesystem.Delete(fileLocation)
				}
			}
			fileName := file[0].Save(filepath.Join(dir, "web/assets/upload/avatars"),
				strings.ReplaceAll(file[0].Header.Filename, " ", ""))
			user.ImageUrl = "/assets/upload/avatars/" + fileName
		}
	}

	user.Username = updateDto.Username
	user.Email = updateDto.Email
	if updateDto.Password != "" {
		user.Password = updateDto.Password
	}
	user.RoleID = updateDto.Role

	if err := database.Conn().Save(&user).Error; err != nil {
		r.Session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.Redirect(302, "/admin/users/list")
	}
	r.Session.SetFlashMessage(ctx, "update data success", "success", nil)
	return ctx.Redirect(http.StatusFound, "/admin/users/list")
}

func (r *Handler) View(ctx echo.Context) error {
	var user User
	id := ctx.Param("id")

	pageHeader := []map[string]interface{}{{
		"menu": "User",
		"link": "/admin/users/list",
	},
		{
			"menu": "View",
			"link": "/admin/users/view/" + id,
		},
	}
	err := database.Conn().Joins("Role").First(&user, "users.id =?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Session.SetFlashMessage(ctx, err.Error(), "error", nil)
			return ctx.Redirect(302, "/admin/users/list")
		}
		r.Session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.Redirect(302, "/admin/users/list")
	}
	return r.Render(ctx, "/user/view", append(r.PageHeader, pageHeader...), user)
}
