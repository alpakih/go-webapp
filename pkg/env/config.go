package env

import (
	"fmt"
	"github.com/spf13/viper"
)


type Environment struct {
	App      app
	Server   server
	Database database
	Auth     auth
}

type app struct {
	Name        string
	Environment string
	Debug       bool
	BcryptCost  int
}

type server struct {
	Host string
	Port int
}
type database struct {
	Connection         string
	Host               string
	Port               int
	Name               string
	Username           string
	Password           string
	Options            string
	MaxOpenConnections int
	MaxIdleConnections int
	MaxLifetime        int
	AutoMigrate        bool
}

type auth struct {
	Jwt jwt
}

type jwt struct {
	Expiry int
	Secret string
}

func LoadEnvironment() {

	var env Environment
	// Set the file name of the configurations file
	viper.SetConfigName("env")

	// Set the path to look for the configurations file
	viper.AddConfigPath(".")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
		panic(err)
	}

	// Set undefined variables
	viper.SetDefault("app.debug", true)
	viper.SetDefault("app.bcryptCost", 10)
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("database.maxOpenConnections", 20)
	viper.SetDefault("database.maxIdleConnections", 20)
	viper.SetDefault("database.maxLifetime", 300)
	viper.SetDefault("database.autoMigrate", false)

	err := viper.Unmarshal(&env)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
		panic(err)
	}

}
