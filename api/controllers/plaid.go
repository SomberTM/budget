package controllers

import (
	"budget/api/environment"
	"budget/api/middleware"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateLinkToken(e *environment.Environment, c *gin.Context) {
	user, ok := middleware.GetCurrentUser(c)
	if !ok {
		return
	}

	linkToken, err := e.Services.Plaid.GetLinkTokenForUser(c.Request.Context(), user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"linkToken": linkToken,
	})
}

type ExchangePublicTokenRequest struct {
	PublicToken string `json:"publicToken"`
}

func ExchangePublicToken(e *environment.Environment, c *gin.Context) {
	user, ok := middleware.GetCurrentUser(c)
	if !ok {
		return
	}

	var request ExchangePublicTokenRequest
	if err := c.BindJSON(&request); err != nil {
		log.Println("Error exchanging public token:", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	e.Services.Plaid.ExchangePublicToken(c.Request.Context(), user, request.PublicToken)

	c.JSON(http.StatusOK, gin.H{})
}

func GetAccounts(e *environment.Environment, c *gin.Context) {
	user, ok := middleware.GetCurrentUser(c)
	if !ok {
		return
	}

	accounts, err := e.Services.Plaid.GetUserAccounts(c.Request.Context(), user)
	if err != nil {
		log.Printf("Error getting accounts: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, accounts)
}
