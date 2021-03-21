package seeder

import (
	"github.com/alpakih/go-webapp/internal/role"
	"github.com/alpakih/go-webapp/pkg/database"
)

func RoleSeeder() map[string]interface{} {
	db := database.Conn()
	var count int64
	mapRole := map[string]interface{}{}

	if err := db.Model(&role.Role{}).Count(&count).Error; err != nil {
		panic(err)
	}
	if count == 0 {
		entity := role.Seeder()
		if err := db.Model(&role.Role{}).Create(&entity).Error; err != nil {
			panic(err)
		}
		mapRole = map[string]interface{}{}
		for _, v := range entity {
			mapRole[v.Slug] = v.ID
		}
	}
	return mapRole
}
