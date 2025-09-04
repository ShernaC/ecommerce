package model

import (
	"encoding/json"
	"time"
)

type Order struct {
	ID              int             `json:"id" gorm:"type:int;primaryKey;autoIncrement"`
	UserID          int             `json:"user_id" gorm:"type:int;unique;not null"`
	Status          string          `json:"status" gorm:"type:varchar(50);not null"`
	TotalAmount     float64         `json:"total_amount" gorm:"type:decimal(10,2);not null;"`
	ShippingAddress json.RawMessage `json:"shipping_address" gorm:"type:jsonb;not null"`
	PaymentMethod   json.RawMessage `json:"payment_method" gorm:"type:jsonb;not null"`
	CreatedAt       time.Time       `json:"created_at" gorm:"type:timestamp;not null"`
	UpdatedAt       *time.Time      `json:"updated_at" gorm:"type:timestamp;null"`
	DeletedAt       *time.Time      `json:"deleted_at" gorm:"type:timestamp;null"`
	Items           []OrderItem     `json:"items" gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
}

type OrderItem struct {
	ID                     int             `json:"id" gorm:"type:int;primaryKey;autoIncrement"`
	OrderID                int             `json:"order_id" gorm:"type:int;unique;not null"`
	ProductID              int             `json:"product_id" gorm:"type:int;unique;not null"`
	Quantity               int             `json:"quantity" gorm:"type:int;unique;not null"`
	PriceAtPurchase        float64         `json:"price_at_purchase" gorm:"type:decimal(10,2);not null;"`
	ProductSnapshotDetails json.RawMessage `json:"product_snapshot_details" gorm:"type:jsonb;not null"`
	CreatedAt              time.Time       `json:"created_at" gorm:"type:timestamp;not null"`
	UpdatedAt              *time.Time      `json:"updated_at" gorm:"type:timestamp;null"`
	DeletedAt              *time.Time      `json:"deleted_at" gorm:"type:timestamp;null"`
}
