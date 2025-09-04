package tools

import (
	"time"
	"users/model"
)

func TimeToDateTimeStringWithoutT(input time.Time) string {
	return input.Format("2006-01-02")
}

func UserToUserData(user *model.User) *model.UserData {
	return &model.UserData{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
	}
}
