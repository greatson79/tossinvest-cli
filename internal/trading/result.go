package trading

type MutationResult struct {
	Kind                  string   `json:"kind"`
	Status                string   `json:"status"`
	OrderID               string   `json:"order_id,omitempty"`
	OriginalOrderID       string   `json:"original_order_id,omitempty"`
	CurrentOrderID        string   `json:"current_order_id,omitempty"`
	Symbol                string   `json:"symbol,omitempty"`
	Market                string   `json:"market,omitempty"`
	Quantity              float64  `json:"quantity,omitempty"`
	FilledQuantity        float64  `json:"filled_quantity,omitempty"`
	Price                 float64  `json:"price,omitempty"`
	AverageExecutionPrice float64  `json:"average_execution_price,omitempty"`
	OrderDate             string   `json:"order_date,omitempty"`
	Warnings              []string `json:"warnings,omitempty"`
}
