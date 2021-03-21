package user

import (
	"github.com/alpakih/go-webapp/internal/role"
	"github.com/alpakih/go-webapp/pkg/helper"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        string    `gorm:"type:varchar(60);column:id;primary_key:true"`
	Username  string    `gorm:"type:varchar(100);column:username"`
	Email     string    `gorm:"type:varchar(50);column:email;unique"`
	Password  string    `gorm:"type:varchar(200);column:password"`
	ImageUrl  string    `gorm:"type:varchar(255);column:image_url"`
	RoleID    string    `gorm:"type:varchar(60);column:role_id;"`
	Role      role.Role `gorm:"foreignkey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (c User) TableName() string {
	return "users"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (c *User) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New().String()

	return
}

func Seeder(roles map[string]interface{}) []User {
	password, _ := helper.HashPassword("123123")
	var users = []User{
		{Username: "admin", Email: "admin@admin.com", Password: password, RoleID: roles["super-admin"].(string)},
		{Username: "manager", Email: "admin@manager.com", Password: password, RoleID: roles["manager"].(string)},
	}

	return users
}
