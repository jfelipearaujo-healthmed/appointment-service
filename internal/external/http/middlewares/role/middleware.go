package role

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Role string

const (
	Doctor  Role = "doctor"
	Patient Role = "patient"
	Any     Role = "any"
)

func IsRole(role string, compareTo Role) bool {
	if compareTo == Any {
		return Role(role) == Doctor || Role(role) == Patient
	}

	return Role(role) == compareTo
}

func GetRoleByName(roleName string) Role {
	switch roleName {
	case "doctor":
		return Doctor
	case "patient":
		return Patient
	default:
		return Any
	}
}

func Middleware(roleAllowed Role) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role := c.Get("role").(string)

			if !IsRole(role, roleAllowed) {
				return echo.NewHTTPError(http.StatusUnauthorized, "You are not authorized to perform this action")
			}

			c.Set("role", role)

			return next(c)
		}
	}
}
