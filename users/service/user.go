package service

import (
	"context"
	"fmt"
	"strings"
	"time"
	"users/model"
	"users/tools"
	"utils/middleware"

	"github.com/google/uuid"
)

func (s *Service) UserRegister(ctx context.Context, input model.NewUser) (*model.User, error) {
	validInput, err := s.UserOnCreate(ctx, input)
	if !validInput {
		panic(err)
	}
	if err != nil {
		return nil, err
	}

	hashedPw, err := tools.HashAndSalt(input.Password)
	if err != nil {
		panic(err)
	}

	var user = model.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPw,
		Phone:    input.Phone,
		Address:  &input.Address,
	}

	if err := s.DB.Model(&user).Create(&user).Error; err != nil {
		panic(err)
	}

	return &user, nil
}

func (s *Service) UserOnCreate(ctx context.Context, input model.NewUser) (bool, error) {
	if input.Email == "" || input.Name == "" || input.Password == "" || input.Phone == "" || input.ConfirmPassword == "" {
		return false, fmt.Errorf("data cannot be empty")
	}

	input.Email = strings.ToLower(input.Email)
	input.Email = strings.TrimSpace(input.Email)

	validEmail := tools.CheckEmailValidity(input.Email)
	if !validEmail {
		return false, fmt.Errorf("invalid email")
	}

	validPhone := tools.CheckPhoneValidity(input.Phone)
	if !validPhone {
		return false, fmt.Errorf("invalid phone number")
	}

	validEmail, err := s.UserCheckEmailValid(ctx, input.Email)
	if err != nil {
		return false, err
	}
	if !validEmail {
		return false, fmt.Errorf("email has been registered")
	}

	validPhone, err = s.UserCheckPhoneValid(ctx, input.Phone)
	if err != nil {
		return false, err
	}
	if !validPhone {
		return false, fmt.Errorf("invalid phone number")
	}

	if input.Password != input.ConfirmPassword {
		return false, fmt.Errorf("passwords do not match")
	}

	return true, nil

}

func (s *Service) UserLogin(ctx context.Context, input model.UserLogin) (*model.UserLoginResponse, error) {
	var role string

	if input.Email == "" || input.Password == "" {
		panic(fmt.Errorf("email or password cannot be empty"))
	}

	input.Email = strings.ToLower(input.Email)
	input.Email = strings.TrimSpace(input.Email)
	user, _ := s.UserGetByEmail(ctx, input.Email)

	valid, err := tools.CompareHash(user.Password, input.Password)
	if err != nil {
		panic(err)
	}
	if !valid {
		panic(fmt.Errorf("invalid password"))
	}

	// check if is seller
	isSeller, _ := s.SellerCheckIsValid(ctx, user.ID)
	if !isSeller {
		role = "user"
	} else {
		role = "seller"
	}

	jti := uuid.New().String()

	//create refresh token
	refreshToken, err := tools.CreateToken(user.ID, input.Email, role, 24*time.Hour, jti)
	if err != nil {
		panic(err)
	}

	// store jti in db
	if _, err = s.UserUpdateRememberToken(ctx, user.ID, jti); err != nil {
		panic(err)
	}

	// create access token
	accessToken, err := tools.CreateToken(user.ID, input.Email, role, 30*time.Minute, "")
	if err != nil {
		panic(err)
	}

	// return
	return &model.UserLoginResponse{
		Success: true,
		Message: "login successful",
		Data: []*model.UserLoginResponseNode{
			{
				TokenType: "access",
				Token:     accessToken,
				UserData:  *tools.UserToUserData(user),
			},
			{
				TokenType: "refresh",
				Token:     refreshToken,
				UserData:  *tools.UserToUserData(user),
			},
		},
	}, nil

}

func (s *Service) UserGetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	if err := s.DB.Model(&user).Scopes(tools.IsDeletedAtNull).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Service) UserGetByID(ctx context.Context, id int) (*model.User, error) {
	var user model.User

	if err := s.DB.Model(&user).Scopes(tools.IsDeletedAtNull).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Service) UserUpdateRememberToken(ctx context.Context, id int, jti string) (*model.User, error) {
	var (
		user model.User
	)

	if err := s.DB.Model(&user).Scopes(tools.IsDeletedAtNull).Where("id = ?", id).Update("remember_token", jti).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Service) UserCheckEmailValid(ctx context.Context, email string) (bool, error) {
	var count int64

	if err := s.DB.Table("user").Scopes(tools.IsDeletedAtNull).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}

	if int(count) > 0 {
		return false, nil
	}

	return true, nil
}

func (s *Service) UserCheckPhoneValid(ctx context.Context, phone string) (bool, error) {
	var count int64

	if err := s.DB.Table("user").Scopes(tools.IsDeletedAtNull).Where("phone = ?", phone).Count(&count).Error; err != nil {
		return false, err
	}

	if int(count) > 0 {
		return false, nil
	}

	return true, nil
}

func (s *Service) UserCheckTokenValid(ctx context.Context, userID int, jti string) (bool, error) {
	var user model.User

	if err := s.DB.Model(&user).Scopes(tools.IsDeletedAtNull).Where("id = ?", userID).First(&user).Error; err != nil {
		return false, err
	}

	return *user.RememberToken == jti, nil
}

func (s *Service) UserRevokeToken(ctx context.Context) error {
	var (
		ctxData = middleware.AuthContext(ctx)
	)
	result := s.DB.Table("user").Where("id = ?", ctxData.ID).Update("remember_token", "")

	if result.Error != nil {
		return fmt.Errorf("failed to revoke token: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no user found to revoke token")
	}

	return nil
}
