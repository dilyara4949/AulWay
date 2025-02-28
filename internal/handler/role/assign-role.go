package role

//import (
//	"aulway/internal/handler/auth"
//	"aulway/internal/handler/user"
//	fbAuth "firebase.google.com/go/auth"
//	"github.com/labstack/echo/v4"
//	"net/http"
//)
//
//func AssignUserRole(authClient *fbAuth.Client, userService *user.Service) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		var req struct {
//			UID  string `json:"uid"`
//			Role string `json:"role"`
//		}
//		if err := c.Bind(&req); err != nil {
//			return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
//		}
//
//		// Validate role
//		if req.Role != "admin" && req.Role != "user" && req.Role != "manager" {
//			return echo.NewHTTPError(http.StatusBadRequest, "invalid role")
//		}
//
//		// Assign role in Firebase
//		err := auth.AssignRole(authClient, req.UID, req.Role)
//		if err != nil {
//			return echo.NewHTTPError(http.StatusInternalServerError, "failed to assign role")
//		}
//
//		// Update role in database
//		err = userService.UpdateUserRole(req.UID, req.Role)
//		if err != nil {
//			return echo.NewHTTPError(http.StatusInternalServerError, "failed to update user role")
//		}
//
//		return c.JSON(http.StatusOK, map[string]string{"message": "role updated"})
//	}
//}
