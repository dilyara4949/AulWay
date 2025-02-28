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
	middleware "aulway/internal/transport/middlware"
	"aulway/internal/utils/config"
	fbAuth "firebase.google.com/go/auth"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"

	_ "aulway/docs"
)

const (
	AdminRole = "admin"
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

	e.POST("/auth/firebase-signin", auth.FirebaseSignIn(userService, r.fb))

	e.PUT("/users/:userId", user.UpdateUserHandler(userService))
	e.GET("/users/:userId", user.GetUserByIdHandler(userService))
	e.GET("/users", user.GetUsersList(userService))

	e.POST("/buses", bus.CreateBusHandler(busService, r.c))
	e.GET("/buses/:busId", bus.GetBusHandler(busService, r.c))

	e.POST("/routes", route.CreateRouteHandler(routeService, busService, r.c))
	e.GET("/routes/:routeId", route.GetRouteHandler(routeService, r.c))
	e.PUT("/routes/:routeId", route.UpdateRouteHandler(routeService, r.c))
	e.DELETE("/routes/:routeId", route.DeleteRouteHandler(routeService, r.c))
	e.GET("/routes", route.GetRoutesListHandler(routeService, r.c))

	// ------- user APIs
	//
	publicProtected := e.Group("/api", middleware.FirebaseAuthMiddleware(r.fb))
	//
	publicProtected.PUT("/users/:userId", user.UpdateUserHandler(userService))
	publicProtected.GET("/users/:userId", user.GetUserByIdHandler(userService))
	//
	//publicProtected.GET("/buses/:busId", bus.GetBusHandler(busService, r.c))
	//
	//publicProtected.GET("/routes/:routeId", route.GetRouteHandler(routeService, r.c))
	//publicProtected.GET("/routes", route.GetRoutesListHandler(routeService, r.c))
	//
	//// ----- admin APIs
	//
	//adminProtected := e.Group("/admin", middleware.AccessCheckMiddleware(AdminRole))
	//
	//adminProtected.PUT("/users/:userId", user.UpdateUserHandler(userService))
	//adminProtected.GET("/users/userId", user.GetUserByIdHandler(userService))
	//adminProtected.GET("/users", user.GetUsersList(userService))
	//
	//adminProtected.POST("/buses", bus.CreateBusHandler(busService, r.c))
	//adminProtected.GET("/buses/:busId", bus.GetBusHandler(busService, r.c))
	//
	//adminProtected.POST("/routes", route.CreateRouteHandler(routeService, busService, r.c))
	//adminProtected.GET("/routes/:routeId", route.GetRouteHandler(routeService, r.c))
	//adminProtected.PUT("/routes/:routeId", route.UpdateRouteHandler(routeService, r.c))
	//adminProtected.DELETE("/routes/:routeId", route.DeleteRouteHandler(routeService, r.c))
	//adminProtected.GET("/routes", route.GetRoutesListHandler(routeService, r.c))

	return e
}

// ----- test APIs

//e.PUT("/users/:userId", user.UpdateUserHandler(userService))
//e.GET("/users/userId", user.GetUserByIdHandler(userService))
//adminProtected.GET("/users", user.GetUsersList(userService))
//
//e.POST("/buses", bus.CreateBusHandler(busService, r.c))
//e.GET("/buses/:busId", bus.GetBusHandler(busService, r.c))
//
//e.POST("/routes", route.CreateRouteHandler(routeService, busService, r.c))
//e.GET("/routes/:routeId", route.GetRouteHandler(routeService, r.c))
//e.PUT("/routes/:routeId", route.UpdateRouteHandler(routeService, r.c))
//e.DELETE("/routes/:routeId", route.DeleteRouteHandler(routeService, r.c))
//e.GET("/routes", route.GetRoutesListHandler(routeService, r.c))
