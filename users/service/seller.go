package service

import (
	"context"
	"fmt"
	"users/middleware"
	"users/model"
	"users/tools"

	"gorm.io/gorm"
)

type ApprovalType string

const (
	APPROVAL_TYPE_PENDING  ApprovalType = "pending"
	APPROVAL_TYPE_APPROVED ApprovalType = "approved"
	APPROVAL_TYPE_REJECTED ApprovalType = "rejected"
)

func (s *Service) SellerRegister(ctx context.Context, input model.NewSeller) (*model.Seller, error) {
	var ctxData = middleware.AuthContext(ctx)

	validInput, err := s.SellerOnCreate(ctx, input)
	if !validInput {
		panic(err)
	}
	if err != nil {
		return nil, err
	}

	userDetails, err := s.UserGetByID(ctx, ctxData.ID)
	if err != nil {
		return nil, err
	}

	var seller = model.Seller{
		ID:           userDetails.ID,
		BusinessName: input.BusinessName,
		Address:      input.Address,
		IsApproved:   string(APPROVAL_TYPE_PENDING),
	}

	if err := s.DB.Model(&seller).Create(&seller).Error; err != nil {
		panic(err)
	}

	return &seller, nil
}

func (s *Service) SellerOnCreate(ctx context.Context, input model.NewSeller) (bool, error) {
	if input.BusinessName == "" || input.Address == "" {
		return false, fmt.Errorf("data cannot be empty")
	}

	exists, err := s.SellerCheckExist(ctx, 0)
	if err != nil {
		return false, err
	}
	if exists {
		return false, fmt.Errorf("seller already exist")
	}

	return true, nil
}

func (s *Service) SellerCheckExist(ctx context.Context, id int) (bool, error) {
	var (
		ctxData = middleware.AuthContext(ctx)
		count   int64
	)

	if id == 0 {
		id = ctxData.ID
	}

	if err := s.DB.Table("seller").Scopes(tools.IsDeletedAtNull).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}

	if int(count) > 0 {
		return true, nil
	}

	return false, nil
}

func (s *Service) SellerCheckIsValid(ctx context.Context, id int) (bool, error) {
	seller, err := s.SellerGetByID(ctx, id)
	if seller == nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	if seller.IsApproved == (string)(APPROVAL_TYPE_PENDING) {
		return false, fmt.Errorf("seller account is still pending approval")
	} else if seller.IsApproved == (string)(APPROVAL_TYPE_REJECTED) {
		return false, fmt.Errorf("seller account is rejected")
	}

	return true, nil
}

func (s *Service) SellerUpdateApproval(ctx context.Context, id int, approval ApprovalType) error {
	var (
		seller model.Seller
	)

	if err := s.DB.Model(&seller).Scopes(tools.IsDeletedAtNull).Where("id = ?", id).Update("is_approved", approval).Error; err != nil {
		return err
	}

	return nil
}

func (s *Service) SellerGetByID(ctx context.Context, id int) (*model.Seller, error) {
	var seller model.Seller

	if err := s.DB.Model(&seller).Scopes(tools.IsDeletedAtNull).Where("id = ?", id).First(&seller).Error; err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("seller not found")
	} else if err != nil {
		return nil, err
	}

	return &seller, nil
}
