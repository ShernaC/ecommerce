package model

import (
	"time"
)

type Order struct {
	ID              int          `json:"id" gorm:"type:int;primaryKey;autoIncrement"`
	UserID          int          `json:"user_id" gorm:"type:int;not null"`
	Status          string       `json:"status" gorm:"type:varchar(50);not null"`
	TotalAmount     float64      `json:"total_amount" gorm:"type:decimal(10,2);not null;"`
	ShippingAddress string       `json:"shipping_address" gorm:"type:varchar(255);not null"`
	PaymentMethod   string       `json:"payment_method" gorm:"type:varchar(100);not null"`
	CreatedAt       time.Time    `json:"created_at" gorm:"type:timestamp;not null"`
	UpdatedAt       *time.Time   `json:"updated_at" gorm:"type:timestamp;null"`
	DeletedAt       *time.Time   `json:"deleted_at" gorm:"type:timestamp;null"`
	Items           []*OrderItem `json:"items" gorm:"-"`
}

type OrderItem struct {
	ID              int        `json:"id" gorm:"type:int;primaryKey;autoIncrement"`
	OrderID         int        `json:"order_id" gorm:"type:int;not null"`
	ProductID       int        `json:"product_id" gorm:"type:int;not null"`
	Quantity        int        `json:"quantity" gorm:"type:int;not null"`
	PriceAtPurchase float64    `json:"price_at_purchase" gorm:"type:decimal(10,2);not null;"`
	ProductSnapshot string     `json:"product_snapsho" gorm:"type:string;not null"`
	CreatedAt       time.Time  `json:"created_at" gorm:"type:timestamp;not null"`
	UpdatedAt       *time.Time `json:"updated_at" gorm:"type:timestamp;null"`
	DeletedAt       *time.Time `json:"deleted_at" gorm:"type:timestamp;null"`
}

type NewOrderItem struct {
	OrderID         int         `json:"order_id"`
	ProductID       int         `json:"product_id"`
	Quantity        int         `json:"quantity"`
	PriceAtPurchase float64     `json:"price_at_purchase"`
	ProductSnapshot string      `json:"product_snapshot"`
	ID              int         `json:"id" gorm:"type:int;primaryKey;autoIncrement"`
	UserID          int         `json:"user_id" gorm:"type:int;unique;not null"`
	Status          string      `json:"status" gorm:"type:varchar(50);not null"`
	TotalAmount     float64     `json:"total_amount" gorm:"type:decimal(10,2);not null;"`
	ShippingAddress string      `json:"shipping_address" gorm:"type:varchar(255);not null"`
	PaymentMethod   string      `json:"payment_method" gorm:"type:varchar(255);not null"`
	CreatedAt       time.Time   `json:"created_at" gorm:"type:timestamp;not null"`
	UpdatedAt       *time.Time  `json:"updated_at" gorm:"type:timestamp;null"`
	DeletedAt       *time.Time  `json:"deleted_at" gorm:"type:timestamp;null"`
	Items           []OrderItem `json:"items" gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
}

type OrderResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Data    []*Order `json:"data"`
}

type OrderTracking struct {
	ID          int       `json:"id" gorm:"type:int;primaryKey;autoIncrement"`
	OrderID     int       `json:"order_id" gorm:"type:int;not null"`
	Status      string    `json:"status" gorm:"type:varchar(100);not null"`
	Description string    `json:"description" gorm:"type:varchar(100);not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:timestamp;not null"`
}

type OrderTrackingResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    []*OrderTracking `json:"data"`
}

type ProductSnapshot struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	SellerID        int       `json:"seller_id"`
	ShopName        string    `json:"shop_name"`
	PriceAtPurchase float64   `json:"price_at_purchase"`
	SKU             string    `json:"sku"`
	PrimaryImage    *string   `json:"primary_image"`
	TaxCategory     *string   `json:"tax_category"`
	CapturedAt      time.Time `json:"captured_at"`
}
