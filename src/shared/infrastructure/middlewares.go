package infrastructure

import (
	"fmt"

	shared_errors "github.com/UPB-Code-Labs/main-api/src/shared/domain/errors"
	"github.com/gin-gonic/gin"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors[0]

			switch e := err.Err.(type) {
			case shared_errors.DomainError:
				c.JSON(e.StatusCode(), gin.H{
					"message": e.Error(),
				})
			default:
				c.JSON(500, gin.H{
					"message": "There was an error processing your request",
				})
			}
		}
	}
}

func WithAuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("session")
		if err != nil {
			c.Error(shared_errors.UnauthorizedError{
				Message: "No session was provided",
			})
			c.Abort()
			return
		}

		claims, err := GetJwtTokenHandler().ValidateToken(cookie)
		if err != nil {
			c.Error(shared_errors.UnauthorizedError{
				Message: "Invalid session",
			})
			c.Abort()
			return
		}

		// Set session data in the chain context
		c.Set("session_uuid", claims.UUID)
		c.Set("session_role", claims.Role)
		c.Next()
	}
}

func WithAuthorizationMiddleware(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionRole, _ := c.Get("session_role")

		if sessionRole != role {
			c.Error(shared_errors.NotEnoughPermissionsError{
				Message: fmt.Sprintf("%s role is required", role),
			})
			c.Abort()
		}

		c.Next()
	}
}
