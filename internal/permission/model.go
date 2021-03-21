package permission

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Permission struct {
	ID          string    `gorm:"type:varchar(60);column:id;primary_key:true"`
	Group       string    `gorm:"type:varchar(50);column:group_permission;unique"`
	Feature     string    `gorm:"type:varchar(50);column:feature"`
	Url         string    `gorm:"type:varchar(255);column:url"`
	Description string    `gorm:"type:varchar(100);column:description"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (c *Permission) TableName() string {
	return "permissions"
}

func (c *Permission) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New().String()
	return
}

func Seeder() []Permission {
	var permissions = []Permission{

		{Group: "home", Feature: "Home", Url: "/admin/home", Description: "Home"},

		{Group: "user", Feature: "User Add", Url: "/admin/users/add", Description: "Add users"},
		{Group: "user", Feature: "User List", Url: "/admin/users/list", Description: "List of users"},
		{Group: "user", Feature: "User Edit", Url: "/admin/users/edit/:id", Description: "Edit data user"},
		{Group: "user", Feature: "User Update", Url: "/admin/users/update/:id", Description: "Update data user"},
		{Group: "user", Feature: "User Detail", Url: "/admin/users/detail/:id", Description: "Detail data user"},
		{Group: "user", Feature: "User Delete", Url: "/admin/users/delete/:id", Description: "Delete data user"},
		{Group: "user", Feature: "User Store", Url: "/admin/users/store", Description: "Save data user"},
		{Group: "user", Feature: "User Datatable", Url: "/admin/users/datatable", Description: "List of datatable of users"},

		{Group: "role", Feature: "Role Add", Url: "/admin/roles/add", Description: "Add roles"},
		{Group: "role", Feature: "Role List", Url: "/admin/roles/list", Description: "List of roles"},
		{Group: "role", Feature: "Role Edit", Url: "/admin/roles/edit/:id", Description: "Edit data role"},
		{Group: "role", Feature: "Role Update", Url: "/admin/update/edit/:id", Description: "Update data role"},
		{Group: "role", Feature: "Role Detail", Url: "/admin/roles/detail/:id", Description: "Detail data role"},
		{Group: "role", Feature: "Role Delete", Url: "/admin/roles/delete/:id", Description: "Delete data role"},
		{Group: "role", Feature: "Role Store", Url: "/admin/roles/store", Description: "Save data role"},
		{Group: "role", Feature: "Role Datatable", Url: "/admin/roles/datatable", Description: "List of datatable of roles"},
		{Group: "role", Feature: "Role Select2", Url: "/admin/roles/select2", Description: "List of data role for select2"},

		{Group: "permission", Feature: "Permissions Add", Url: "/admin/permissions/add", Description: "Add permissions"},
		{Group: "permission", Feature: "Permissions List", Url: "/admin/permissions/list", Description: "List of permissions"},
		{Group: "permission", Feature: "Permissions Edit", Url: "/admin/permissions/edit/:id", Description: "Edit data permission"},
		{Group: "permission", Feature: "Permissions Update", Url: "/admin/permissions/update/:id", Description: "Update data permission"},
		{Group: "permission", Feature: "Permissions Detail", Url: "/admin/permissions/detail/:id", Description: "Detail data permission"},
		{Group: "permission", Feature: "Permissions Delete", Url: "/admin/permissions/delete/:id", Description: "Delete data permission"},
		{Group: "permission", Feature: "Permissions Store", Url: "/admin/permissions/store", Description: "Save data permission"},
		{Group: "permission", Feature: "Permissions Datatable", Url: "/admin/permissions/datatable", Description: "List of datatable of permissions"},
		{Group: "permission", Feature: "Permissions Select2", Url: "/admin/permissions/select2", Description: "List of data permissions for select2"},
	}
	return permissions
}
