package http

import (
	"github.com/UPB-Code-Labs/main-api/src/session/application"
	"github.com/UPB-Code-Labs/main-api/src/session/infrastructure/requests"
	"github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/gin-gonic/gin"
)

type SessionControllers struct {
	UseCases *application.SessionUseCases
}

func (controllers *SessionControllers) HandleLogin(c *gin.Context) {
	// Get the request
	var request requests.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	// Validate request body
	if err := infrastructure.GetValidator().Struct(request); err != nil {
		c.JSON(400, gin.H{
			"message": "Validation error",
			"errors":  err.Error(),
		})
		return
	}

	// Call the use case
	dto := request.ToDTO()
	session, err := controllers.UseCases.Login(*dto)
	if err != nil {
		c.Error(err)
		return
	}

	// Set the cookie
	cookieName := "session"
	cookieSecondsTTL := infrastructure.GetEnvironment().JwtExpirationHours * 60 * 60
	cookieDomain := "" // By default, the cookie is only sent to the domain that set it
	cookiePath := "/"
	cookieSecure := infrastructure.GetEnvironment().IsInProduction
	cookieHttpOnly := true

	c.SetCookie(
		cookieName,
		session.Token,
		cookieSecondsTTL,
		cookiePath,
		cookieDomain,
		cookieSecure,
		cookieHttpOnly,
	)

	// Return the response
	c.JSON(200, gin.H{
		"user": gin.H{
			"uuid":      session.User.UUID,
			"full_name": session.User.FullName,
			"role":      session.User.Role,
		},
	})
}

func (controllers *SessionControllers) HandleLogout(c *gin.Context) {
	// Delete the cookie
	cookieName := "session"
	cookieSecondsTTL := 0 // Expires automatically
	cookieDomain := ""
	cookiePath := "/"
	cookieSecure := infrastructure.GetEnvironment().IsInProduction
	cookieHttpOnly := true

	c.SetCookie(
		cookieName,
		"",
		cookieSecondsTTL,
		cookiePath,
		cookieDomain,
		cookieSecure,
		cookieHttpOnly,
	)

	// Return the response
	c.Status(204)
}

func (controllers *SessionControllers) HandleWhoAmI(c *gin.Context) {
	uuid := c.MustGet("session_uuid").(string)

	session, err := controllers.UseCases.WhoAmI(uuid)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{
		"user": gin.H{
			"uuid":      session.UUID,
			"role":      session.Role,
			"full_name": session.FullName,
		},
	})
}
