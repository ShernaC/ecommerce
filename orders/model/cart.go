package model

import "time"

type Cart struct {
	ID        int         `json:"id" gorm:"type:int;primaryKey;autoIncrement"`
	UserID    int         `json:"user_id" gorm:"type:int;unique;not null"`
	CreatedAt time.Time   `json:"created_at" gorm:"type:timestamp;not null"`
	UpdatedAt *time.Time  `json:"updated_at" gorm:"type:timestamp;null"`
	Items     []*CartItem `json:"items" gorm:"-"`
}

type CartItem struct {
	ID        int       `json:"id" gorm:"type:int;primaryKey;autoIncrement"`
	CartID    int       `json:"cart_id" gorm:"type:int;unique;not null"`
	ProductID int       `json:"product_id" gorm:"type:int;unique;not null"`
	Quantity  int       `json:"quantity" gorm:"type:int;not null"`
	Price     float64   `json:"price" gorm:"type:decimal(10,2);not null"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;not null"`
}

type CartItemInput struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type NewCartItem struct {
	CartID    int     `json:"cart_id"`
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type CartResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    Cart   `json:"data"`
}
