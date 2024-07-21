package controllers

import (
	"budget/api/environment"
	"budget/api/middleware"
	"budget/api/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TransactionCategories(e *environment.Environment, c *gin.Context) {
	categories, err := e.Repositories.TransactionCategories.GetTransactionCategories(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func MyBudgets(e *environment.Environment, c *gin.Context) {
	user, ok := middleware.GetCurrentUser(c)
	if !ok {
		return
	}

	budgets, err := e.Repositories.Budgeting.GetBudgetsForUser(c.Request.Context(), user.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, budgets)
}

func BudgetBreakdown(e *environment.Environment, c *gin.Context) {
	user, ok := middleware.GetCurrentUser(c)
	if !ok {
		return
	}

	budgetId := c.Param("budgetId")

	breakdown, err := e.Services.Budgeting.GetBudgetBreakdown(c.Request.Context(), budgetId, user.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, breakdown)
}

func CreateBudget(e *environment.Environment, c *gin.Context) {
	user, ok := middleware.GetCurrentUser(c)
	if !ok {
		return
	}

	var request models.Budget
	if err := c.BindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	request.SetUserId(user.Id)

	budget, err := e.Repositories.Budgeting.CreateBudget(c.Request.Context(), request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, budget)
}

type CreateBudgetDefinitionRequest struct {
	models.BudgetDefinition
	TransactionCategoryIds []string `json:"transaction_category_ids"`
}

func CreateBudgetDefinition(e *environment.Environment, c *gin.Context) {
	user, ok := middleware.GetCurrentUser(c)
	if !ok {
		return
	}

	var request CreateBudgetDefinitionRequest
	if err := c.BindJSON(&request); err != nil {
		log.Printf("Error creating budget definition: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	request.SetUserId(user.Id)

	definition, err := e.Repositories.Budgeting.CreateBudgetDefinition(c.Request.Context(), request.BudgetDefinition)
	if err != nil {
		log.Printf("Error assigning budget definition categories: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	err = e.Repositories.Budgeting.AssignCategoriesToBudgetDefinition(c.Request.Context(), definition.Id, request.TransactionCategoryIds)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, definition)
}
