package pkg

type Expression struct {
	ID         string           `json:"id"`
	Expression string           `json:"expression"`
	Status     ExpressionStatus `json:"status"`
	Result     float64          `json:"result"`
}
