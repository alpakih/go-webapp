package seeder

import (
	"errors"
	"github.com/alpakih/go-webapp/internal/permission"
	"github.com/alpakih/go-webapp/internal/role"
	"github.com/alpakih/go-webapp/internal/user"
	"github.com/alpakih/go-webapp/pkg/database"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func Run() {
	if viper.GetString("app.environment") == "local" {

		log.Info("Truncate seeders...")
		if err := truncateTable(); err != nil {
			panic(err)
		}
		log.Info("Running seeders...")
		roles := RoleSeeder()
		UserSeeder(roles)
	}
}

func truncateTable() error {
	return database.Conn().Transaction(func(tx *gorm.DB) (err error) {
		if err := tx.Unscoped().Delete(&permission.Permission{}, "1 =?", 1).Error; err != nil {
			return errors.New("error when truncating table permissions")
		}
		log.Info("truncating table permissions")
		if err := tx.Unscoped().Delete(&role.Role{}, "1 =?", 1).Error; err != nil {
			return errors.New("error when truncating table roles")
		}
		log.Info("truncating table roles")
		if err := tx.Unscoped().Delete(&user.User{}, "1 =?", 1).Error; err != nil {
			return errors.New("error when truncating table users")
		}
		return nil
	})
}
