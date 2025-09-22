package service

import (
	"context"
	"fmt"
	grpcclient "orders/grpc_client"
	"utils/user"
)

type UserDetails struct {
	Name    string
	Email   string
	Phone   string
	Address string
}

func (s *Service) GetUserDetails(ctx context.Context, id int) (*UserDetails, error) {
	if id <= 0 {
		return nil, fmt.Errorf("user id is invalid")
	}

	userDetails, err := grpcclient.GetUserDetails(ctx, &user.GetUserDetailsRequest{Id: int64(id)})
	if err != nil {
		return nil, err
	}

	details := UserDetails{
		Name:    userDetails.Name,
		Email:   userDetails.Email,
		Phone:   userDetails.Phone,
		Address: userDetails.Address,
	}

	return &details, nil
}
