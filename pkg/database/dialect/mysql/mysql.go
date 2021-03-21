package mysql

import (
	"github.com/alpakih/go-webapp/pkg/database"
	"gorm.io/driver/mysql"
)

func init() {
	database.RegisterDialect("mysql", "{username}:{password}@({host}:{port})/{name}?{options}", mysql.Open)
}
