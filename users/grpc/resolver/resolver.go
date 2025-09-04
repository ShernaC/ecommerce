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
