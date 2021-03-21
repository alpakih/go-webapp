package model

type ResultDatatablePermission struct {
	ID             string `gorm:"type:varchar(60);column:id;primary_key:true"`
	Group          string `gorm:"type:varchar(50);column:group_permission;unique"`
	Feature        string `gorm:"type:varchar(50);column:feature"`
	Url            string `gorm:"type:varchar(255);column:url"`
	Description    string `gorm:"type:varchar(100);column:description"`
	FilterRowCount *int64 `gorm:"column:filter_row_count"`
}
