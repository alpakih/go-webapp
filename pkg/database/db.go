package database

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Initializer func(*gorm.DB)

type DialectInitializer func(dsn string) gorm.Dialector

type dialect struct {
	template    string
	initializer DialectInitializer
}

var (
	mu           sync.Mutex
	dbConnection *gorm.DB
	initializers []Initializer
	models       []interface{}
	dialects = map[string]dialect{}

	optionPlaceholders = map[string]string{
		"{username}": "database.username",
		"{password}": "database.password",
		"{host}":     "database.host",
		"{name}":     "database.name",
		"{options}":  "database.options",
	}
)

func RegisterModel(model interface{}) {
	models = append(models, model)
}

func GetRegisteredModels() []interface{} {
	return append(make([]interface{}, 0, len(models)), models...)
}

// ClearRegisteredModels unregister all models.
func ClearRegisteredModels() {
	models = []interface{}{}
}


func RegisterDialect(name, template string, initializer DialectInitializer) {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := dialects[name]; ok {
		panic(fmt.Sprintf("Dialect %q already exists", name))
	}
	dialects[name] = dialect{template, initializer}
}

func AddInitializer(initializer Initializer) {
	initializers = append(initializers, initializer)
}

// ClearInitializers remove all database connection initializer functions.
func ClearInitializers() {
	initializers = []Initializer{}
}

func GetConnection() *gorm.DB {
	mu.Lock()
	defer mu.Unlock()
	if dbConnection == nil {
		dbConnection = newConnection()
	}
	return dbConnection
}

// Conn alias for GetConnection.
func Conn() *gorm.DB {
	return GetConnection()
}

// Close the database connections if they exist.
func Close() error {
	var err error = nil
	mu.Lock()
	defer mu.Unlock()
	if dbConnection != nil {
		db, _ := dbConnection.DB()
		err = db.Close()
		dbConnection = nil
	}

	return err
}

// Migrate migrates all registered models.
func Migrate() {
	db := GetConnection()
	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			panic(err)
		}
	}
}

func newConnection() *gorm.DB {

	driver := viper.GetString("database.connection")

	logLevel := logger.Silent
	if viper.GetBool("app.debug") {
		logLevel = logger.Info
	}

	dialect, ok := dialects[driver]
	if !ok {
		panic(fmt.Sprintf("DB Connection %q not supported, forgotten import? %s", driver))
	}

	dsn := dialect.buildDSN()

	db, err := gorm.Open(dialect.initializer(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel), PrepareStmt: true,
	})
	if err != nil {
		panic(err)
	}

	sql, err := db.DB()
	if err != nil {
		panic(err)
	}

	sql.SetMaxOpenConns(viper.GetInt("database.maxOpenConnections"))
	sql.SetMaxIdleConns(viper.GetInt("database.maxIdleConnections"))
	sql.SetConnMaxLifetime(time.Duration(viper.GetInt("database.maxLifetime")) * time.Second)

	for _, initializer := range initializers {
		initializer(db)
	}

	return db
}

func (d dialect) buildDSN() string {
	connStr := d.template
	for k, v := range optionPlaceholders {
		connStr = strings.Replace(connStr, k, viper.GetString(v), 1)
	}
	connStr = strings.Replace(connStr, "{port}", strconv.Itoa(viper.GetInt("database.port")), 1)

	return connStr
}
