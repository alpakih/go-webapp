package permission

import (
	"errors"
	"fmt"
	"github.com/alpakih/go-webapp/internal"
	"github.com/alpakih/go-webapp/internal/permission/model"
	"github.com/alpakih/go-webapp/pkg/database"
	"github.com/alpakih/go-webapp/pkg/sessions"
	"github.com/alpakih/go-webapp/pkg/validation"
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type Handler struct {
	internal.BaseHandler
}

func NewPermissionController(session *sessions.Manager) Handler {
	return Handler{
		BaseHandler: internal.BaseHandler{
			Title: "User Management | Permission",
			Menu:  "User Management",
			PageHeader: []map[string]interface{}{
				{
					"menu": "User Management",
					"link": "/admin/permissions/list",
				},
			},
			Session: session,
		},
	}
}

func (r *Handler) Index(ctx echo.Context) error {
	pageHeader := map[string]interface{}{
		"menu": "Permission",
		"link": "/admin/permissions/list",
	}
	return r.Render(ctx, "/permission/index", append(r.PageHeader, pageHeader), nil)
}

func (r *Handler) Add(ctx echo.Context) error {

	pageHeader := []map[string]interface{}{{
		"menu": "Permission",
		"link": "/admin/permissions/list",
	},
		{
			"menu": "Add",
			"link": "/admin/permissions/add",
		},
	}

	return r.Render(ctx, "/permission/add", append(r.PageHeader, pageHeader...), nil)
}

func (r *Handler) Datatable(ctx echo.Context) error {
	var qResult []model.ResultDatatablePermission

	draw := ctx.Request().URL.Query().Get("draw")
	search := ctx.Request().URL.Query().Get("search[value]")
	start := ctx.Request().URL.Query().Get("start")
	length := ctx.Request().URL.Query().Get("length")
	order := ctx.Request().URL.Query().Get("order[0][column]")
	orderName := ctx.Request().URL.Query().Get("columns[" + order + "][name]")
	orderAscDesc := ctx.Request().URL.Query().Get("order[0][dir]")

	var recordTotal int64
	var recordFiltered int64
	if err := database.Conn().Table("permissions").Count(&recordTotal).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	sql := `SELECT
				id,
				group_permission,
				feature,
				url,
				description,
				COUNT (*) OVER () AS filter_row_count
			FROM
				permissions`

	if search != "" {
		sql += fmt.Sprintf(` WHERE group LIKE '%s'`, `%`+search+`%`)
		sql += fmt.Sprintf(` OR feature LIKE '%s'`, `%`+search+`%`)
		sql += fmt.Sprintf(` OR url LIKE '%s'`, `%`+search+`%`)
	}

	if orderName != "" && orderAscDesc != "" {
		sql += ` ORDER BY ` + orderName + ` ` + orderAscDesc
	} else {
		sql += ` ORDER BY created_at DESC`
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
					<a href="/admin/permissions/edit/` + v.ID + `" style="text-decoration: none;font-weight: 400; color: #333;"  data-toggle="tooltip" data-placement="right" title="Edit"><i class="fa fa-edit"></i>Edit </a>
				  </li>`
		action += `<li style="list-style-type: none;display: inline">
					<a href="/admin/permissions/view/` + v.ID + `" style="text-decoration: none;font-weight: 400; color: #333;"  data-toggle="tooltip" data-placement="right" title="View"><i class="fa fa-search"></i>View</a>
				  </li>`
		action += `<li style="list-style-type: none;display: inline">
					<a href="JavaScript:void(0);" onclick="Delete('` + v.ID + `')" style="text-decoration: none;font-weight: 400; color: #333;" data-toggle="tooltip" data-placement="right" title="Delete"><i class="fa fa-trash" style="color: #ff4d65"></i>Delete</a>
				  </li>`
		listOfData[k] = map[string]interface{}{
			"id":          v.ID,
			"group":       v.Group,
			"feature":     v.Feature,
			"url":         v.Url,
			"description": v.Description,
			"action":      action,
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

func (r *Handler) Store(ctx echo.Context) error {
	var request StoreDto
	if err := ctx.Bind(&request); err != nil {
		r.Session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.Redirect(302, "/admin/permissions/add")
	}
	if err := ctx.Validate(&request); err != nil {
		var validationErrors []validation.ErrorValidation
		if errs, ok := err.(validator.ValidationErrors); ok {
			validationErrors = validation.WrapValidationErrors(errs)
		}
		r.Session.SetFlashMessage(ctx, err.Error(), "error", validationErrors)
		return ctx.Redirect(302, "/admin/permissions/add")
	}

	if err := database.Conn().Create(&Permission{
		Group:       request.Group,
		Feature:     request.Feature,
		Url:         request.Url,
		Description: request.Description,
	}).Error; err != nil {
		r.Session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.Redirect(302, "/admin/permissions/add")
	}
	r.Session.SetFlashMessage(ctx, "save data success", "success", nil)
	return ctx.Redirect(http.StatusFound, "/admin/permissions/list")
}

func (r *Handler) Edit(ctx echo.Context) error {
	var permission Permission
	id := ctx.Param("id")

	pageHeader := []map[string]interface{}{{
		"menu": "Permission",
		"link": "/admin/permissions/list",
	},
		{
			"menu": "Edit",
			"link": "/admin/permissions/edit/" + id,
		},
	}
	err := database.Conn().First(&permission, "id =?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Session.SetFlashMessage(ctx, err.Error(), "error", nil)
			return ctx.Redirect(302, "/admin/permissions/list")
		}
		r.Session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.Redirect(302, "/admin/permissions/list")
	}
	return r.Render(ctx, "/permission/edit", append(r.PageHeader, pageHeader...), permission)
}

func (r *Handler) Update(ctx echo.Context) error {
	id := ctx.Param("id")
	var permission Permission
	var request UpdateDto
	if err := ctx.Bind(&request); err != nil {
		r.Session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.Redirect(302, "/admin/permissions/edit/"+id)
	}
	if err := ctx.Validate(&request); err != nil {
		var validationErrors []validation.ErrorValidation
		if errs, ok := err.(validator.ValidationErrors); ok {
			validationErrors = validation.WrapValidationErrors(errs)
		}
		r.Session.SetFlashMessage(ctx, err.Error(), "error", validationErrors)
		return ctx.Redirect(302, "/admin/permissions/edit/"+id)
	}

	err := database.Conn().First(&permission, "id =?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Session.SetFlashMessage(ctx, err.Error(), "error", nil)
			return ctx.Redirect(302, "/admin/permissions/list")
		}
		r.Session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.Redirect(302, "/admin/permissions/list")
	}

	permission.Group = request.Group
	permission.Feature = request.Feature
	permission.Url = request.Url
	permission.Description = request.Description
	if err := database.Conn().Save(&permission).Error; err != nil {
		r.Session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.Redirect(302, "/admin/permissions/list")
	}
	r.Session.SetFlashMessage(ctx, "update data success", "success", nil)
	return ctx.Redirect(http.StatusFound, "/admin/permissions/list")
}

func (r *Handler) View(ctx echo.Context) error {
	var permission Permission
	id := ctx.Param("id")

	pageHeader := []map[string]interface{}{{
		"menu": "Permission",
		"link": "/admin/permissions/list",
	},
		{
			"menu": "View",
			"link": "/admin/permissions/view/" + id,
		},
	}
	err := database.Conn().First(&permission, "id =?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Session.SetFlashMessage(ctx, err.Error(), "error", nil)
			return ctx.Redirect(302, "/admin/permissions/list")
		}
		r.Session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.Redirect(302, "/admin/permissions/list")
	}
	return r.Render(ctx, "/permission/view", append(r.PageHeader, pageHeader...), permission)
}

func (r *Handler) ListSelect2(ctx echo.Context) error {
	search := ctx.QueryParam("term")

	var data []Permission
	if search != "" {
		if err := database.Conn().Find(&data, "feature LIKE ? OR group_permission LIKE ?", "%"+search+"%","%"+search+"%").Error; err != nil {
			return ctx.JSON(500, echo.Map{"message": "error get permissions"})
		}
		result := make([]map[string]interface{}, len(data))
		for k, v := range data {
			result[k] = map[string]interface{}{
				"id":   v.ID,
				"text": v.Feature,
			}
		}
		return ctx.JSON(200, result)
	}

	if err := database.Conn().Find(&data).Error; err != nil {
		return ctx.JSON(500, echo.Map{"message": "error get permissions"})
	}
	result := make([]map[string]interface{}, len(data))
	for k, v := range data {
		result[k] = map[string]interface{}{
			"id":   v.ID,
			"text": v.Feature,
		}
	}

	return ctx.JSON(200, result)
}

func (r *Handler) Delete(ctx echo.Context) error {
	var permission Permission
	id := ctx.Param("id")

	if err := database.Conn().First(&permission, "id=?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.JSON(400, echo.Map{"message": "permission not found"})
		}
		return ctx.JSON(500, echo.Map{"message": "error when trying delete data"})
	}

	if err := database.Conn().Delete(&Permission{}, "id=?", id).Error; err != nil {
		return ctx.JSON(500, echo.Map{"message": "error when trying delete data"})
	}

	return ctx.JSON(200, echo.Map{"message": "delete data has been success"})
}
