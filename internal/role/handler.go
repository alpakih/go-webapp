package role

import (
	"errors"
	"fmt"
	"github.com/alpakih/go-webapp/internal"
	"github.com/alpakih/go-webapp/internal/role/model"
	"github.com/alpakih/go-webapp/internal/rolepermission"
	"github.com/alpakih/go-webapp/pkg/database"
	"github.com/alpakih/go-webapp/pkg/sessions"
	"github.com/alpakih/go-webapp/pkg/validation"
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	internal.BaseHandler
}

func NewRoleController(session *sessions.Manager) Handler {
	return Handler{
		BaseHandler: internal.BaseHandler{
			Title: "User Management | Role",
			Menu:  "User Management",
			PageHeader: []map[string]interface{}{
				{
					"menu": "User Management",
					"link": "/admin/roles/list",
				},
			},
			Session: session,
		},
	}
}

func (r *Handler) Index(ctx echo.Context) error {
	pageHeader := map[string]interface{}{
		"menu": "Role",
		"link": "/admin/roles/list",
	}
	return r.Render(ctx, "/role/index", append(r.PageHeader, pageHeader), nil)
}

func (r *Handler) Add(ctx echo.Context) error {

	pageHeader := []map[string]interface{}{{
		"menu": "Role",
		"link": "/admin/roles/list",
	},
		{
			"menu": "Add",
			"link": "/admin/roles/add",
		},
	}

	return r.Render(ctx, "/role/add", append(r.PageHeader, pageHeader...), nil)
}

func (r *Handler) Datatable(ctx echo.Context) error {
	var qResult []model.ResultDatatableRole

	loc, _ := time.LoadLocation("Asia/Jakarta")
	draw := ctx.Request().URL.Query().Get("draw")
	search := ctx.Request().URL.Query().Get("search[value]")
	start := ctx.Request().URL.Query().Get("start")
	length := ctx.Request().URL.Query().Get("length")
	order := ctx.Request().URL.Query().Get("order[0][column]")
	orderName := ctx.Request().URL.Query().Get("columns[" + order + "][name]")
	orderAscDesc := ctx.Request().URL.Query().Get("order[0][dir]")

	var recordTotal int64
	var recordFiltered *int64
	if err := database.Conn().Table("roles").Count(&recordTotal).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	sql := `SELECT
				roles.id,
				roles.slug,
				roles.role_name,
				roles.description,
				roles.created_at,
				COUNT (*) OVER () AS filter_row_count
			FROM
				roles`

	if search != "" {
		sql += fmt.Sprintf(` WHERE roles.slug LIKE '%s'`, `%`+search+`%`)
		sql += fmt.Sprintf(` OR roles.role_name LIKE '%s'`, `%`+search+`%`)
	}

	if orderName != "" && orderAscDesc != "" {
		sql += ` ORDER BY ` + orderName + ` ` + orderAscDesc
	} else {
		sql += ` ORDER BY roles.created_at DESC`
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
					<a href="/admin/roles/edit/` + v.ID + `" style="text-decoration: none;font-weight: 400; color: #333;"  data-toggle="tooltip" data-placement="right" title="Edit"><i class="fa fa-edit"></i>Edit </a>
				  </li>`
		action += `<li style="list-style-type: none;display: inline">
					<a href="/admin/roles/view/` + v.ID + `" style="text-decoration: none;font-weight: 400; color: #333;"  data-toggle="tooltip" data-placement="right" title="View"><i class="fa fa-search"></i>View</a>
				  </li>`
		action += `<li style="list-style-type: none;display: inline">
					<a href="JavaScript:void(0);" onclick="Delete('` + v.ID + `')" style="text-decoration: none;font-weight: 400; color: #333;" data-toggle="tooltip" data-placement="right" title="Delete"><i class="fa fa-trash" style="color: #ff4d65"></i>Delete</a>
				  </li>`
		listOfData[k] = map[string]interface{}{
			"id":          v.ID,
			"slug":        v.Slug,
			"role_name":   v.RoleName,
			"description": v.Description,
			"created_at":  v.CreatedAt.In(loc).Format("2006-01-02 15:04:05"),
			"action":      action,
		}
	}
	if len(qResult) != 0 {
		recordFiltered = qResult[0].FilterRowCount
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"draw":            draw,
		"recordsTotal":    recordTotal,
		"recordsFiltered": recordFiltered,
		"data":            listOfData,
	})
}

func (r *Handler) Store(ctx echo.Context) error {
	var request Store
	if err := ctx.Bind(&request); err != nil {
		r.Session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.Redirect(302, "/admin/roles/add")
	}
	if err := ctx.Validate(&request); err != nil {
		var validationErrors []validation.ErrorValidation
		if errs, ok := err.(validator.ValidationErrors); ok {
			validationErrors = validation.WrapValidationErrors(errs)
		}
		r.Session.SetFlashMessage(ctx, err.Error(), "error", validationErrors)
		return ctx.Redirect(302, "/admin/roles/add")
	}
	tx := database.Conn().Begin()
	if err := tx.Error; err != nil {
		r.Session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.Redirect(302, "/admin/roles/add")
	}

	defer func(ctx echo.Context) error {
		if re := recover(); re != nil {
			tx.Rollback()
			r.Session.SetFlashMessage(ctx, fmt.Sprintf("%v", re), "error", nil)
			return ctx.Redirect(302, "/admin/roles/add")
		}
		return nil
	}(ctx)

	entity := Role{
		Slug:        request.Slug,
		RoleName:    request.RoleName,
		Description: request.Description,
	}
	if err := tx.Omit("Permission").Create(&entity).Error; err != nil {
		tx.Rollback()
		r.Session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.Redirect(302, "/admin/roles/add")
	}
	listPermissions := make([]map[string]interface{}, len(request.Permission))
	for k, v := range request.Permission {
		listPermissions[k] = map[string]interface{}{
			"role_id":       entity.ID,
			"permission_id": v,
		}
	}

	if err := tx.Model(&rolepermission.RolePermission{}).Create(&listPermissions).Error; err != nil {
		tx.Rollback()
		r.Session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.Redirect(302, "/admin/roles/add")
	}
	tx.Commit()
	r.Session.SetFlashMessage(ctx, "save data success", "success", nil)
	return ctx.Redirect(http.StatusFound, "/admin/roles/list")
}

func (r *Handler) Edit(ctx echo.Context) error {
	var role Role
	id := ctx.Param("id")

	pageHeader := []map[string]interface{}{{
		"menu": "Role",
		"link": "/admin/roles/list",
	},
		{
			"menu": "Edit",
			"link": "/admin/roles/edit/" + id,
		},
	}
	err := database.Conn().Preload("Permission").First(&role, "roles.id =?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Session.SetFlashMessage(ctx, err.Error(), "error", nil)
			return ctx.Redirect(302, "/admin/roles/list")
		}
		r.Session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.Redirect(302, "/admin/roles/list")
	}
	return r.Render(ctx, "/role/edit", append(r.PageHeader, pageHeader...), role)
}

func (r *Handler) Update(ctx echo.Context) error {
	id := ctx.Param("id")
	var role Role
	var request Update
	if err := ctx.Bind(&request); err != nil {
		r.Session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.Redirect(302, "/admin/roles/edit/"+id)
	}
	if err := ctx.Validate(&request); err != nil {
		var validationErrors []validation.ErrorValidation
		if errs, ok := err.(validator.ValidationErrors); ok {
			validationErrors = validation.WrapValidationErrors(errs)
		}
		r.Session.SetFlashMessage(ctx, err.Error(), "error", validationErrors)
		return ctx.Redirect(302, "/admin/roles/edit/"+id)
	}

	tx := database.Conn()
	err := tx.First(&role, "id =?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Session.SetFlashMessage(ctx, err.Error(), "error", nil)
			return ctx.Redirect(302, "/admin/roles/list")
		}
		r.Session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.Redirect(302, "/admin/roles/list")
	}

	if len(request.Permission) > 0 {
		if err := tx.Unscoped().Delete(&rolepermission.RolePermission{}, "role_id =?", id).Error; err != nil {
			tx.Rollback()
			r.Session.SetFlashMessage(ctx, err.Error(), "error", nil)
			return ctx.Redirect(302, "/admin/roles/list")
		}
		listPermissions := make([]rolepermission.RolePermission, len(request.Permission))
		for k, v := range request.Permission {
			listPermissions[k] = rolepermission.RolePermission{
				RoleID:       id,
				PermissionID: v,
			}
		}
		if err := tx.Model(&rolepermission.RolePermission{}).Create(listPermissions).Error; err != nil {
			tx.Rollback()
			r.Session.SetFlashMessage(ctx, err.Error(), "error", nil)
			return ctx.Redirect(302, "/admin/roles/list")
		}
	}

	role.Slug = request.Slug
	role.Description = request.Description
	role.RoleName = request.RoleName

	if err := tx.Save(&role).Error; err != nil {
		tx.Rollback()
		r.Session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.Redirect(302, "/admin/roles/list")
	}
	tx.Commit()
	r.Session.SetFlashMessage(ctx, "update data success", "success", nil)
	return ctx.Redirect(http.StatusFound, "/admin/roles/list")
}

func (r *Handler) View(ctx echo.Context) error {
	var role Role
	id := ctx.Param("id")

	pageHeader := []map[string]interface{}{{
		"menu": "Role",
		"link": "/admin/roles/list",
	},
		{
			"menu": "View",
			"link": "/admin/roles/view/" + id,
		},
	}
	err := database.Conn().Preload("Permission").First(&role, "roles.id =?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Session.SetFlashMessage(ctx, err.Error(), "error", nil)
			return ctx.Redirect(302, "/admin/roles/list")
		}
		r.Session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.Redirect(302, "/admin/roles/list")
	}
	return r.Render(ctx, "/role/view", append(r.PageHeader, pageHeader...), role)
}

func (r *Handler) ListSelect2(ctx echo.Context) error {
	search := ctx.QueryParam("term")

	var data []Role
	if search != "" {
		if err := database.Conn().Find(&data, "slug LIKE ?", "%"+search+"%").Error; err != nil {
			return ctx.JSON(500, echo.Map{"message": "error get roles"})
		}
		result := make([]map[string]interface{}, len(data))
		for k, v := range data {
			result[k] = map[string]interface{}{
				"id":   v.ID,
				"text": v.RoleName,
			}
		}
		return ctx.JSON(200, result)
	}

	if err := database.Conn().Find(&data).Error; err != nil {
		return ctx.JSON(500, echo.Map{"message": "error get roles"})
	}
	result := make([]map[string]interface{}, len(data))
	for k, v := range data {
		result[k] = map[string]interface{}{
			"id":   v.ID,
			"text": v.RoleName,
		}
	}

	return ctx.JSON(200, result)
}

func (r *Handler) Delete(ctx echo.Context) error {
	var role Role
	id := ctx.Param("id")

	if err := database.Conn().First(&role, "id=?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.JSON(400, echo.Map{"message": "role not found"})
		}
		return ctx.JSON(500, echo.Map{"message": "error when trying delete data"})
	}

	if err := database.Conn().Delete(&Role{}, "id=?", id).Error; err != nil {
		return ctx.JSON(500, echo.Map{"message": "error when trying delete data"})
	}

	return ctx.JSON(200, echo.Map{"message": "delete data has been success"})
}
