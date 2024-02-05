package infrastructure

import (
	sharedErrors "github.com/UPB-Code-Labs/main-api/src/shared/domain/errors"
	"github.com/gin-gonic/gin"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors[0]

			switch e := err.Err.(type) {
			case sharedErrors.DomainError:
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
			c.Error(sharedErrors.UnauthorizedError{
				Message: "You must be logged in",
			})
			c.Abort()
			return
		}

		claims, err := GetJwtTokenHandler().ValidateToken(cookie)
		if err != nil {
			c.Error(sharedErrors.UnauthorizedError{
				Message: "Your session has expired or is not valid",
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

func WithAuthorizationMiddleware(role []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionRole, _ := c.Get("session_role")

		var isRoleAuthorized bool
		for _, r := range role {
			if sessionRole == r {
				isRoleAuthorized = true
				break
			}
		}

		if !isRoleAuthorized {
			c.Error(sharedErrors.NotEnoughPermissionsError{})
			c.Abort()
		}

		c.Next()
	}
}

func WithServerSentEventsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set headers to allow server sent events
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		c.Writer.Header().Set("X-Accel-Buffering", "no")

		c.Next()
	}
}
