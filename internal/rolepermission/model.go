package rolepermission

type RolePermission struct {
	RoleID       string    `gorm:"type:varchar(60);column:role_id"`
	PermissionID string    `gorm:"type:varchar(60);column:permission_id"`
}

func (c *RolePermission) TableName() string {
	return "role_permissions"
}
