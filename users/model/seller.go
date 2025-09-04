package model

import "time"

type Seller struct {
	ID           int        `json:"id" gorm:"type:int;primaryKey;"`
	BusinessName string     `json:"business_name" gorm:"type:varchar(255);not null"`
	Address      string     `json:"address" gorm:"type:varchar(255);not null"`
	IsApproved   string     `json:"is_approved" gorm:"type:varchar(15);default:false"`
	CreatedAt    time.Time  `json:"created_at" gorm:"type:timestamp;not null"`
	UpdatedAt    *time.Time `json:"updated_at" gorm:"type:timestamp;null"`
	DeletedAt    *time.Time `json:"deleted_at" gorm:"type:timestamp;null"`
	User         User       `gorm:"foreignKey:ID;references:ID"`
}

type NewSeller struct {
	BusinessName string `json:"business_name"`
	Address      string `json:"address"`
}

type SellerResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	Data    SellerData `json:"seller_data"`
}

type SellerData struct {
	UserData     UserData  `json:"user_data"`
	BusinessName string    `json:"business_name"`
	Address      string    `json:"address"`
	IsApproved   string    `json:"is_approved"`
	CreatedAt    time.Time `json:"created_at"`
}
