package accrual

type Service interface {
	GetOrderInfo(orderID string) ([]byte, error)
}
