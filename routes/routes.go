package routes

import (
	"github.com/TypicalAM/HackathonSoftswiss2/config"
	"github.com/TypicalAM/HackathonSoftswiss2/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// controller holds all the variables needed for routes to perform their logic
type controller struct {
	db *gorm.DB
}

// New creates a new router with all the routes
func New(db *gorm.DB, cfg *config.Config) (*gin.Engine, error) {
	store := cookie.NewStore([]byte(cfg.CookieSecret))

	// Allow cors
	corsCofig := cors.DefaultConfig()
	corsCofig.AllowOrigins = cfg.TrustedOrigins
	corsCofig.AllowCredentials = true

	// Default middleware
	router := gin.Default()
	router.Use(cors.New(corsCofig))
	router.Use(sessions.Sessions("save_poznan", store))
	router.Use(middleware.Session(db))
	router.Use(middleware.General())

	controller := controller{db: db}

	// Set up the api
	api := router.Group("/api")

	// Non-Authorized routes, we cannot access with an active session
	noAuth := api.Group("/")
	noAuth.Use(middleware.NoAuth())
	noAuth.Use(middleware.Throttle(cfg.RequestsPerMin))
	noAuth.POST("/register", controller.Register)
	noAuth.POST("/login", controller.Login)

	// Authorized routes, we can access with only an active session
	auth := api.Group("/")
	auth.Use(middleware.Auth())
	auth.Use(middleware.Sensitive())
	auth.GET("/profile", controller.Profile)
	auth.POST("/logout", controller.Logout)

	// Product routes
	products := auth.Group("/products")
	products.Use(middleware.Auth())
	products.Use(middleware.Throttle(cfg.RequestsPerMin))
	products.GET("/check/:EAN", controller.CheckEAN)
	products.POST("/throw", controller.ThrowAway)
	products.GET("/stats", controller.GlobalInfo)

	return router, nil
}
