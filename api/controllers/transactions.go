package controllers

import (
	"budget/api/environment"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AccountTransactions(e *environment.Environment, c *gin.Context) {
	accountId := c.Param("accountId")

	transactions, err := e.Repositories.Transactions.GetTransactionsForAccount(c.Request.Context(), accountId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	c.JSON(http.StatusOK, transactions)
}
