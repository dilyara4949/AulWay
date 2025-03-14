package http

import (
	"aulway/internal/handler/auth"
	"aulway/internal/handler/bus"
	"aulway/internal/handler/healthz"
	"aulway/internal/handler/page"
	"aulway/internal/handler/route"
	"aulway/internal/handler/ticket"
	"aulway/internal/handler/user"
	busRepostory "aulway/internal/repository/bus"
	pageRepository "aulway/internal/repository/page"
	paymentRepostory "aulway/internal/repository/payment"
	routeRepostory "aulway/internal/repository/route"
	ticketRepository "aulway/internal/repository/ticket"
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

	authService := service.NewAuthService(userRepo, nil)

	busRepo := busRepostory.New(r.db)
	busService := service.NewBusService(busRepo)

	routeRepo := routeRepostory.New(r.db)
	routeService := service.NewRouteService(routeRepo)

	paymentRepo := paymentRepostory.New(r.db)
	paymentService := service.NewFPaymentProcessor()

	ticketRepo := ticketRepository.New(r.db)
	ticketService := service.NewTicketService(ticketRepo, paymentRepo, routeRepo, paymentService)

	pageRepo := pageRepository.New(r.db)
	pageService := service.NewPageService(pageRepo)

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

	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:8080", "https://yourfrontend.com"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	e.GET("/health", healthz.CheckHealth())
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.POST("/auth/signup", auth.SignupHandler(authService, userService, r.c))
	e.POST("/auth/signin", auth.SigninHandler(authService, userService, r.c))

	publicProtected := e.Group("/api", middleware.JWTAuth(r.c.JWTTokenSecret))

	adminProtected := e.Group("/api", middleware.JWTAuth(r.c.JWTTokenSecret), middleware.AccessCheck(AdminRole))

	publicProtected.PUT("/users/:userId", user.UpdateUserHandler(userService))
	publicProtected.GET("/users/:userId", user.GetUserByIdHandler(userService))
	adminProtected.GET("/users", user.GetUsersList(userService))
	publicProtected.DELETE("/users/:userId", user.DeleteUserHandler(userService))
	publicProtected.PUT("/users/:userId/change-password", user.ChangePasswordHandler(userService))

	adminProtected.GET("/buses", bus.GetBusesListHandler(busService, r.c))
	adminProtected.POST("/buses", bus.CreateBusHandler(busService, r.c))
	adminProtected.GET("/buses/:busId", bus.GetBusHandler(busService, r.c))

	adminProtected.GET("/all-routes", route.GetAllRoutesListHandler(routeService, r.c))
	adminProtected.POST("/routes", route.CreateRouteHandler(routeService, busService, r.c))
	publicProtected.GET("/routes/:routeId", route.GetRouteHandler(routeService, r.c))
	adminProtected.PUT("/routes/:routeId", route.UpdateRouteHandler(routeService, r.c))
	adminProtected.DELETE("/routes/:routeId", route.DeleteRouteHandler(routeService, r.c))
	publicProtected.GET("/routes", route.GetRoutesListHandler(routeService, r.c))

	publicProtected.POST("/tickets/:routeId", ticket.BuyTicketHandler(ticketService))
	publicProtected.GET("/tickets/users/:userId", ticket.GetUserTicketsHandler(ticketService))
	publicProtected.GET("/tickets/users/:userId/:ticketId", ticket.GetTicketDetailsHandler(ticketService))

	publicProtected.PUT("/pages/:title", page.UpdatePageHandler(pageService))
	adminProtected.GET("/pages/:title", page.GetPageHandler(pageService))

	return e
}
