package models

type Budget struct {
	Id     string `json:"id"`
	UserId string `json:"user_id"`
	Name   string `json:"name"`
	Color  string `json:"color"`
}

func NewBudget() Budget {
	return Budget{}
}

func NewUserBudget(userId string) Budget {
	budget := NewBudget()
	budget.SetUserId(userId)
	return budget
}

func (b *Budget) SetUserId(userId string) *Budget {
	b.UserId = userId
	return b
}

func (b *Budget) SetName(name string) *Budget {
	b.Name = name
	return b
}

func (b *Budget) SetColor(color string) *Budget {
	b.Color = color
	return b
}
