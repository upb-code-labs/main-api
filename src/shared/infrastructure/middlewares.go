package infrastructure

import (
	"github.com/UPB-Code-Labs/main-api/src/shared/domain"
	"github.com/gin-gonic/gin"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors[0]

			switch e := err.Err.(type) {
			case domain.DomainError:
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
