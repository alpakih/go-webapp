package model

import "time"

type ResultDatatableRole struct {
	ID             string `gorm:"type:varchar(60);column:id;primary_key:true"`
	Slug           string `gorm:"type:varchar(50);column:slug;unique"`
	RoleName       string `gorm:"type:varchar(50);column:role_name"`
	Description    string `gorm:"type:varchar(100);column:description"`
	FilterRowCount *int64 `gorm:"column:filter_row_count"`
	CreatedAt      time.Time `gorm:"column:created_at"`

}
