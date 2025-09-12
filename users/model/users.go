package model

import "time"

type User struct {
	ID            int        `json:"id" gorm:"type:int;primaryKey;autoIncrement"`
	Name          string     `json:"name" gorm:"type:varchar(100);not null"`
	Password      string     `json:"password" gorm:"type:varchar(100);not null"`
	Email         string     `json:"email" gorm:"type:varchar(100);unique;not null"`
	Phone         string     `json:"phone" gorm:"type:varchar(20);not null"`
	Address       *string    `json:"address" gorm:"type:varchar(255);null"`
	RememberToken *string    `json:"remember_token" gorm:"type:varchar(100);null"`
	CreatedAt     time.Time  `json:"created_at" gorm:"type:timestamp;not null"`
	UpdatedAt     *time.Time `json:"updated_at" gorm:"type:timestamp;null"`
	DeletedAt     *time.Time `json:"deleted_at" gorm:"type:timestamp;null"`
}

type UserData struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address   *string   `json:"address"`
	CreatedAt time.Time `json:"created_at"`
}

type NewUser struct {
	Name            string `json:"name"`
	Password        string `json:"password"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Address         string `json:"address"`
	ConfirmPassword string `json:"confirm_password"`
}

type UserResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Data    UserData `json:"user_data"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	Success bool                     `json:"success"`
	Message string                   `json:"message"`
	Data    []*UserLoginResponseNode `json:"data"`
}

type UserLoginResponseNode struct {
	TokenType string   `json:"token_type"`
	Token     string   `json:"token"`
	UserData  UserData `json:"user"`
}
