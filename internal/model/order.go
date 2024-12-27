package model

import "time"

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusComplete  OrderStatus = "complete"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID         int64       `json:"id" gorm:"primaryKey"`
	UserID     int64       `json:"user_id" gorm:"not null"`
	ProductID  int64       `json:"product_id" gorm:"not null"`
	Quantity   int         `json:"quantity" gorm:"not null"`
	TotalPrice float64     `json:"total_price" gorm:"not null"`
	Status     OrderStatus `json:"status" gorm:"not null"`
	TxHash     string      `json:"tx_hash"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}
