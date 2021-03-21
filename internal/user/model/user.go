package model

import "time"

type ResultDatatableUser struct {
	ID             string    `gorm:"type:varchar(60);column:id;primary_key:true"`
	Username       string    `gorm:"type:varchar(100);column:username"`
	Email          string    `gorm:"type:varchar(50);column:email;unique"`
	ImageUrl       string    `gorm:"type:varchar(255);column:image_url"`
	RoleName       string    `gorm:"type:varchar(60);column:role_name;"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	FilterRowCount *int64    `gorm:"column:filter_row_count"`
}
