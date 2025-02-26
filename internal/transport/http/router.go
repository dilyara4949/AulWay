package http

import (
	"aulway/internal/handler/auth"
	"aulway/internal/handler/bus"
	"aulway/internal/handler/healthz"
	"aulway/internal/handler/route"
	"aulway/internal/handler/user"
	busRepostory "aulway/internal/repository/bus"
	routeRepostory "aulway/internal/repository/route"
	userRepository "aulway/internal/repository/user"
	"aulway/internal/service"
	"aulway/internal/utils/config"
	fbAuth "firebase.google.com/go/auth"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"

	_ "aulway/docs"
)

type Router struct {
	c  config.Config
	db *gorm.DB
	fb *fbAuth.Client
}

func NewRouter(c config.Config, db *gorm.DB, fbAuth *fbAuth.Client) *Router {
	return &Router{
		c:  c,
		db: db,
		fb: fbAuth,
	}
}

func (r *Router) Build() *echo.Echo {
	userRepo := userRepository.NewRepository(r.db)
	userService := service.NewUserService(userRepo)

	busRepo := busRepostory.New(r.db)
	busService := service.NewBusService(busRepo)

	routeRepo := routeRepostory.New(r.db)
	routeService := service.NewRouteService(routeRepo)

	// Middleware
	timeoutWithConfig := echoMiddleware.TimeoutWithConfig(
		echoMiddleware.TimeoutConfig{
			Skipper:      echoMiddleware.DefaultSkipper,
			ErrorMessage: "timeout error",
			Timeout:      r.c.HeaderTimeout,
		})

	e := echo.New()
	e.HideBanner = true
	e.Debug = true
	e.Use(echoMiddleware.Recover(), timeoutWithConfig, echoMiddleware.Logger())

	e.GET("/health", healthz.CheckHealth())
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	services := e.Group("")
	services.POST("/signin", auth.FirebaseSignIn(userService, r.fb))

	//firebaseAuthMiddleware := middleware.FirebaseAuthMiddleware(r.fb)
	protected := services.Group("")

	protected.PUT("/user/:userId", user.UpdateUserHandler(userService))

	protected.POST("/bus", bus.CreateBusHandler(busService, r.c))
	protected.GET("/bus/:busId", bus.GetBusHandler(busService, r.c))

	protected.POST("/route", route.CreateRouteHandler(routeService, busService, r.c))
	protected.GET("/route/:routeId", route.GetRouteHandler(routeService, r.c))
	protected.PUT("/route/:routeId", route.UpdateRouteHandler(routeService, r.c))
	protected.DELETE("/route/:routeId", route.DeleteRouteHandler(routeService, r.c))
	protected.GET("/route", route.GetRoutesListHandler(routeService, r.c))

	return e
}
