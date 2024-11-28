package middleware

import (
	"budget/api/environment"
	"budget/api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireLoggedInUser(e *environment.Environment, c *gin.Context) {
	sessionId, idErr := c.Cookie("session_id")
	sessionToken, tokErr := c.Cookie("session_token")
	if tokErr != nil || sessionToken == "" || idErr != nil || sessionId == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "no user"})
		return
	}

	session, user, exists, err := e.Repositories.Sessions.GetSessionWithUser(c.Request.Context(), sessionId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if !session.CheckToken([]byte(sessionToken)) {
		models.ClearSessionCookie(c)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("user", user)
	c.Next()
}

func GetCurrentUser(c *gin.Context) (models.User, bool) {
	maybeUser, ok := c.Get("user")
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return models.User{}, false
	}

	user, ok := maybeUser.(models.User)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Invalid user shape"})
		return models.User{}, false
	}

	return user, true
}
