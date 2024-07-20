package controllers

import (
	"budget/api/environment"
	"budget/api/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// We receive a username and password
// - If no user with the username, create user and session, set cookie
// - If a user with the username and the password is correct, generate a new session, set cookie
func Login(e *environment.Environment, c *gin.Context) {
	userName, password, ok := c.Request.BasicAuth()
	if !ok {
		log.Printf("Login requested with invalid Authorization header")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Unsupported Authorization header"})
		return
	}

	session, token, err := e.Services.Users.Login(c.Request.Context(), userName, password)
	if err != nil {
		log.Printf("Error loggin user in. Likely incorrect password. Error: %v", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	session.SetSessionCookie(c, token)
	c.JSON(http.StatusNoContent, gin.H{})
}

func Logout(e *environment.Environment, c *gin.Context) {
	sessionId, idErr := c.Cookie("session_id")
	sessionToken, tokErr := c.Cookie("session_token")
	if tokErr != nil || sessionToken == "" || idErr != nil || sessionId == "" {
		models.ClearSessionCookie(c)
		log.Println("Error parsing, or empty, session token or session id")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err := e.Services.Users.Logout(c.Request.Context(), sessionId, sessionToken)
	if err != nil {
		log.Printf("Error logging user out. %v", err)
		models.ClearSessionCookie(c)
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	models.ClearSessionCookie(c)
	c.JSON(http.StatusNoContent, gin.H{})
}

func Me(c *gin.Context) {
	value, authorized := c.Get("user")
	if !authorized {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user, ok := value.(models.User)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Invalid user value in context"})
		return
	}

	c.JSON(http.StatusOK, user)
}
