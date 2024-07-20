package models

type BudgetDefinition struct {
	Id            string `json:"id"`
	UserId        string `json:"user_id"`
	BudgetId      string `json:"budget_id"`
	MaxAllocation int64  `json:"max_allocation"`
}

func NewBudgetDefinition(budgetId string) BudgetDefinition {
	def := BudgetDefinition{}
	def.SetBudgetId(budgetId)
	return def
}

func NewUserBudgetDefinition(userId string, budgetId string) BudgetDefinition {
	def := NewBudgetDefinition(budgetId)
	def.SetUserId(userId)
	return def
}

func (bd *BudgetDefinition) SetUserId(userId string) *BudgetDefinition {
	bd.UserId = userId
	return bd
}

func (bd *BudgetDefinition) SetBudgetId(budgetId string) *BudgetDefinition {
	bd.BudgetId = budgetId
	return bd
}

func (bd *BudgetDefinition) SetMaxAllocation(maxAllocation int64) *BudgetDefinition {
	bd.MaxAllocation = maxAllocation
	return bd
}
