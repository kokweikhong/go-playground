package model

type Finance struct {
	Expenses *Expenses
}

type Expenses struct {
	Category []*ExpensesCategory
}

type ExpensesCategory struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Remarks   string `json:"remarks"`
	CreatedAt string `json:"created_at"`
}

func NewFinance() *Finance {
	return &Finance{
		Expenses: &Expenses{},
	}
}
