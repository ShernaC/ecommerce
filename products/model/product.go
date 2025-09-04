package model

import "time"

type Product struct {
	ID          int        `json:"id" gorm:"type:int;primaryKey;"`
	SellerID    int        `json:"seller_id" gorm:"type:int;not null;"`
	Name        string     `json:"name" gorm:"type:varchar(100);not null;"`
	Description string     `json:"description" gorm:"type:text;"`
	Price       float64    `json:"price" gorm:"type:decimal(10,2);not null;"`
	Stock       int        `json:"stock" gorm:"type:int;not null;"`
	CreatedAt   time.Time  `json:"created_at" gorm:"type:timestamp;not null"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"type:timestamp;null"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"type:timestamp;null"`
}

type NewProduct struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	SellerID    int     `json:"-"`
}

type UpdateProduct struct {
	ID          int
	Name        *string
	Description *string
	Price       *float64
	Stock       *int
}
