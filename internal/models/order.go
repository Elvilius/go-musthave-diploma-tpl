package models

type OrderStatus string

const (
	NEW        OrderStatus = "NEW"
	REGISTERED OrderStatus = "REGISTERED"
	INVALID    OrderStatus = "INVALID"
	PROCESSING OrderStatus = "PROCESSING"
	PROCESSED  OrderStatus = "PROCESSED"
)

type Order struct {
	Number     string      `json:"number"`
	Status     OrderStatus `json:"status"`
	Accrual    float32     `json:"accrual"`
	UploadedAt string      `json:"uploaded_at"`
	UserID     int         `json:"-"`
}

type ExternalOrder struct {
	Order   string
	Status  string
	Accrual float32
}
