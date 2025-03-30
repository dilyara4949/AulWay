package http

import (
	"aulway/internal/handler/auth"
	"aulway/internal/handler/bus"
	favorite "aulway/internal/handler/favorites"
	"aulway/internal/handler/healthz"
	"aulway/internal/handler/page"
	"aulway/internal/handler/route"
	"aulway/internal/handler/ticket"
	"aulway/internal/handler/user"
	busRepostory "aulway/internal/repository/bus"
	favRepository "aulway/internal/repository/favorite"
	pageRepository "aulway/internal/repository/page"
	paymentRepostory "aulway/internal/repository/payment"
	routeRepostory "aulway/internal/repository/route"
	ticketRepository "aulway/internal/repository/ticket"
	userRepository "aulway/internal/repository/user"
	"aulway/internal/service"
	middleware "aulway/internal/transport/middlware"
	"aulway/internal/utils/config"
	"aulway/internal/utils/logger"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"

	_ "aulway/docs"
)

const (
	AdminRole = "admin"
)

type Router struct {
	c     config.Config
	db    *gorm.DB
	redis *redis.Client
}

func NewRouter(c config.Config, db *gorm.DB, redis *redis.Client) *Router {
	return &Router{
		c:     c,
		db:    db,
		redis: redis,
	}
}

func (r *Router) Build() *echo.Echo {
	userRepo := userRepository.NewRepository(r.db)
	userService := service.NewUserService(userRepo)

	authService := service.NewAuthService(userRepo, r.redis, r.c.SMTP)

	busRepo := busRepostory.New(r.db)
	busService := service.NewBusService(busRepo)

	routeRepo := routeRepostory.New(r.db)
	routeService := service.NewRouteService(routeRepo)

	paymentRepo := paymentRepostory.New(r.db)
	paymentService := service.NewFPaymentProcessor()

	ticketRepo := ticketRepository.New(r.db)
	ticketService := service.NewTicketService(ticketRepo, paymentRepo, routeRepo, paymentService, busRepo)

	pageRepo := pageRepository.New(r.db)
	pageService := service.NewPageService(pageRepo)

	favRepo := favRepository.New(r.db)
	favService := service.NewFavoriteService(favRepo)

	timeoutWithConfig := echoMiddleware.TimeoutWithConfig(
		echoMiddleware.TimeoutConfig{
			Skipper:      echoMiddleware.DefaultSkipper,
			ErrorMessage: "timeout error",
			Timeout:      r.c.HeaderTimeout,
		})

	e := echo.New()
	e.HideBanner = true
	e.Debug = true

	aulLogger := logger.New()
	e.Use(aulLogger.LogRequest)

	e.Use(echoMiddleware.Recover(), timeoutWithConfig)

	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:8080", "https://localhost:5173"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	e.GET("/health", healthz.CheckHealth())
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.POST("/auth/signup", auth.SignupHandler(r.redis, r.c, userService))
	e.POST("/auth/signup/verify", auth.VerifyEmailHandler(r.redis, userService, authService, r.c))
	e.POST("/auth/signin", auth.SigninHandler(authService, userService, r.c))
	e.POST("/auth/forgot-password", auth.ForgotPasswordHandler(authService))
	e.POST("/auth/forgot-password/verify", auth.VerifyForgotPasswordHandler(authService))

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
	adminProtected.DELETE("/buses/:busId", bus.DeleteBusHandler(busService))

	adminProtected.GET("/all-routes", route.GetAllRoutesListHandler(routeService, r.c))
	adminProtected.POST("/routes", route.CreateRouteHandler(routeService, busService, r.c))
	publicProtected.GET("/routes/:routeId", route.GetRouteHandler(routeService, busService, r.c))
	adminProtected.PUT("/routes/:routeId", route.UpdateRouteHandler(routeService, r.c))
	adminProtected.DELETE("/routes/:routeId", route.DeleteRouteHandler(routeService, r.c))
	publicProtected.GET("/routes", route.GetRoutesListHandler(routeService, r.c))

	publicProtected.POST("/tickets/:routeId", ticket.BuyTicketHandler(ticketService, r.c))
	publicProtected.GET("/tickets/users/:userId", ticket.GetUserTicketsHandler(ticketService))
	publicProtected.GET("/tickets/users/:userId/:ticketId", ticket.GetTicketDetailsHandler(ticketService))
	adminProtected.GET("/tickets", ticket.GetTicketsSortByHandler(ticketService))

	adminProtected.PUT("/pages/:title", page.UpdatePageHandler(pageService))
	publicProtected.GET("/pages/:title", page.GetPageHandler(pageService))

	publicProtected.POST("/users/:userId/favorites", favorite.AddFavoriteHandler(favService))
	publicProtected.DELETE("/users/:userId/favorites/:routeId", favorite.RemoveFavoriteHandler(favService))
	publicProtected.GET("/users/:userId/favorites", favorite.GetFavoritesHandler(favService))

	return e
}
