package dbmanager

type BalanceTransaction struct {
	UserID int     `json:"user_id"`
	Money  float64 `json:"money"`
}

type ReserveTransaction struct {
	UserID    int     `json:"user_id"`
	ServiceID int     `json:"service_id"`
	OrderID   int     `json:"order_id"`
	Money     float64 `json:"money"`
}

type TranslateTransaction struct {
	FromID int     `json:"from_id"`
	ToID   int     `json:"to_id"`
	Money  float64 `json:"money"`
}
