package resolver

import (
	"context"
	"users/service"
	"utils/user"
)

type Server struct {
	user.UnimplementedUserServer
}

func (s *Server) CheckSellerExists(ctx context.Context, request *user.CheckSellerExistsRequest) (*user.CheckSellerExistsResponse, error) {
	sellerID := request.Id

	exists, err := service.GetService().SellerCheckExist(ctx, int(sellerID))
	if err != nil {
		return nil, err
	}

	return &user.CheckSellerExistsResponse{Valid: exists}, nil
}

func (s *Server) GetUserDetails(ctx context.Context, request *user.GetUserDetailsRequest) (*user.GetUserDetailsResponse, error) {
	userID := request.Id

	details, err := service.GetService().UserGetByID(ctx, int(userID))
	if err != nil {
		return nil, err
	}

	resp := &user.GetUserDetailsResponse{
		Name:    details.Name,
		Email:   details.Email,
		Phone:   details.Phone,
		Address: *details.Address,
	}

	return resp, nil
}
