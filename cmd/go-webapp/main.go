package main

import (
	"context"
	"encoding/gob"
	"github.com/alpakih/go-webapp/internal/auth"
	"github.com/alpakih/go-webapp/internal/home"
	"github.com/alpakih/go-webapp/internal/permission"
	"github.com/alpakih/go-webapp/internal/role"
	"github.com/alpakih/go-webapp/internal/user"
	"github.com/alpakih/go-webapp/pkg/database"
	_ "github.com/alpakih/go-webapp/pkg/database/dialect/mssql"
	"github.com/alpakih/go-webapp/pkg/env"
	"github.com/alpakih/go-webapp/pkg/logging"
	middlewareFunc "github.com/alpakih/go-webapp/pkg/middleware"
	"github.com/alpakih/go-webapp/pkg/sessions"
	"github.com/alpakih/go-webapp/pkg/validation"
	"github.com/alpakih/go-webapp/web/views"
	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/echoview-v4"
	gorillaCtx "github.com/gorilla/context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

func main() {

	// Load .env file
	env.LoadEnvironment()

	gob.Register(sessions.UserInfo{})
	gob.Register(sessions.FlashMessage{})
	gob.Register(user.User{})
	gob.Register(role.Role{})
	gob.Register(echo.Map{})
	gob.Register(permission.Permission{})
	gob.Register([]validation.ErrorValidation{})

	//Not Found Handler
	echo.NotFoundHandler = func(c echo.Context) error {
		return c.Render(http.StatusOK,"404.html",echo.Map{"title":"Page Not Found"})
	}

	e := echo.New()
	e.Renderer = views.NewRenderer("./web/views/*.html",true)

	e.HTTPErrorHandler = views.CustomHTTPErrorHandler
	// migrate
	if viper.GetBool("database.autoMigrate") {
		database.RegisterModel(user.User{})
		database.RegisterModel(role.Role{})

		database.Migrate()
	}

	//seeder
	//seeder.Run()

	// Set Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	e.Static("/assets", "web/assets")

	// setup log folder,file and log global
	logFile := logging.SetupLogFileAndFolder("web")

	// set logger middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: io.MultiWriter(logFile, os.Stdout),
	}))

	// middleware session
	sessionMiddleware := middlewareFunc.NewSessionMiddleware(viper.GetString("auth.session.key"),
		viper.GetString("auth.session.encryptionKey"))

	// cookie store
	store := sessionMiddleware.NewCookieStore()

	// setup session
	sessionManager := sessions.NewSessionManager(store)

	authorizationMiddleware := middlewareFunc.NewAuthorizationMiddleware(sessionManager)

	e.Validator = validation.NewValidator()

	e.Use(echo.WrapMiddleware(gorillaCtx.ClearHandler))

	// new middleware view backend
	mv := echoview.NewMiddleware(goview.Config{
		Root:      "web/views/backend",
		Extension: ".html",
		Master:    "layout/master",
		Partials: []string{
			"partials/brand-logo",
			"partials/page-header",
			"partials/right-navbar",
			"partials/left-navbar",
			"partials/sidebar-menu",
		},
		Funcs: template.FuncMap{
			"getCsrfToken": func(ctx echo.Context) string {
				return ctx.Get("csrf_token").(string)
			},
			"getUserAuth": func(ctx echo.Context) (userAuth *sessions.UserInfo) {
				session, err := sessionManager.Get(ctx, sessions.IDSession)
				if err != nil {
					return nil
				}
				if userSession, ok := session.(sessions.UserInfo); !ok {
					return nil
				} else {
					userAuth = &userSession
				}
				return userAuth
			},
		},
		DisableCache: true,
	})


	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, "/admin/auth/login")
	})

	// backend group
	backendGroup := e.Group("/admin", mv, middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "form:csrf",
		ContextKey:  "csrf_token",
		Skipper: func(i echo.Context) bool {
			apiV1URI := "/admin"
			if strings.EqualFold(i.Request().RequestURI, apiV1URI+"/auth/login") {
				return true
			}
			return false
		},
	}), sessionMiddleware.SessionMiddleware(sessionManager))

	// home route
	homeController := home.NewHomeController(sessionManager)
	backendGroup.GET("/home", homeController.Index)

	// auth route
	authController := auth.NewAuthController(sessionManager)
	backendGroup.GET("/auth/login", authController.LoginPage)
	backendGroup.POST("/auth/login", authController.Login)
	backendGroup.POST("/auth/logout", authController.Logout)

	// user route
	userController := user.NewUserController(sessionManager)
	userGroup := backendGroup.Group("/users", authorizationMiddleware.AuthorizationMiddleware([]string{"super-admin"}))
	userGroup.GET("/list", userController.Index)
	userGroup.GET("/add", userController.Add)
	userGroup.GET("/datatable", userController.Datatable)
	userGroup.GET("/edit/:id", userController.Edit)
	userGroup.GET("/view/:id", userController.View)
	userGroup.POST("/store", userController.Store)
	userGroup.POST("/update/:id", userController.Update)
	userGroup.DELETE("/delete/:id", userController.Delete)

	// role route
	roleController := role.NewRoleController(sessionManager)
	roleGroup := backendGroup.Group("/roles", authorizationMiddleware.AuthorizationMiddleware([]string{"super-admin"}))
	roleGroup.GET("/list", roleController.Index)
	roleGroup.GET("/add", roleController.Add)
	roleGroup.GET("/datatable", roleController.Datatable)
	roleGroup.GET("/edit/:id", roleController.Edit)
	roleGroup.GET("/view/:id", roleController.View)
	roleGroup.POST("/store", roleController.Store)
	roleGroup.POST("/update/:id", roleController.Update)
	roleGroup.DELETE("/delete/:id", roleController.Delete)
	roleGroup.GET("/select2", roleController.ListSelect2)

	// permission route
	permissionController := permission.NewPermissionController(sessionManager)
	permissionGroup := backendGroup.Group("/permissions", authorizationMiddleware.AuthorizationMiddleware([]string{"super-admin"}))
	permissionGroup.GET("/list", permissionController.Index)
	permissionGroup.GET("/add", permissionController.Add)
	permissionGroup.GET("/datatable", permissionController.Datatable)
	permissionGroup.GET("/edit/:id", permissionController.Edit)
	permissionGroup.GET("/view/:id", permissionController.View)
	permissionGroup.POST("/store", permissionController.Store)
	permissionGroup.POST("/update/:id", permissionController.Update)
	permissionGroup.DELETE("/delete/:id", permissionController.Delete)
	permissionGroup.GET("/select2", permissionController.ListSelect2)

	// Start server
	go func() {
		if err := e.Start(viper.GetString("server.host") + ":" + viper.GetString("server.port")); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
