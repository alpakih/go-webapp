package seeder

import (
	"github.com/alpakih/go-webapp/internal/permission"
	"github.com/alpakih/go-webapp/pkg/database"
)

func PermissionSeed() {
	db := database.Conn()
	var count int64
	if err := db.Model(&permission.Permission{}).Count(&count).Error; err != nil {
		panic(err)
	}
	if count == 0 {
		entity:=permission.Seeder()
		if err := db.Model(&permission.Permission{}).Create(&entity).Error; err != nil {
			panic(err)
		}
	}
}
