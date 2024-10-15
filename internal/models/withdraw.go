package models

type GetWithdraw struct {
	OrderNumber string  `json:"order"`
	Sum         float32 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}
