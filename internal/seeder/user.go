package seeder

import (
	"github.com/alpakih/go-webapp/internal/user"
	"github.com/alpakih/go-webapp/pkg/database"
)

func UserSeeder(roles map[string]interface{}) {
	db := database.Conn()
	var count int64
	if err := db.Model(&user.User{}).Count(&count).Error; err != nil {
		panic(err)
	}
	if count == 0 {
		entity:=user.Seeder(roles)
		if err := db.Model(&user.User{}).Create(&entity).Error; err != nil {
			panic(err)
		}
	}
}
