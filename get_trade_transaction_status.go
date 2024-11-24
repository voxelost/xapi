package xapi

type tradeTransactionStatusInput struct {
	OrderID int `json:"order"`
}

type TradeStatus int

var (
	TradeStatusError    TradeStatus = 1
	TradeStatusPending  TradeStatus = 2
	TradeStatusAccepted TradeStatus = 3
	TradeStatusRejected TradeStatus = 4
)

type TradeTransactionStatus struct {
	Ask           float64     `json:"ask"`
	Bid           float64     `json:"bid"`
	CustomComment string      `json:"customComment"`
	Message       *string     `json:"message"`
	OrderID       int         `json:"order"`
	RequestStatus TradeStatus `json:"requestStatus"`
}

func (c *client) GetTradeTransactionStatus(orderID int) (TradeTransactionStatus, error) {
	return getSync[tradeTransactionStatusInput, TradeTransactionStatus](c, "tradeTransactionStatus", tradeTransactionStatusInput{
		OrderID: orderID,
	})
}
