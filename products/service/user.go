package service

import (
	"context"
	grpcclient "products/grpc_client"
	"utils/user"
)

type SellerDetails struct {
	BusinessName string
}

func CheckSellerExists(userID int) (bool, error) {
	checkSellerReq, err := grpcclient.CheckSellerExists(context.Background(), &user.CheckSellerExistsRequest{Id: int64(userID)})
	if err != nil {
		return false, err
	}

	if !checkSellerReq.Valid {
		return false, nil
	}

	return true, nil
}

func (s *Service) GetSellerDetails(id int) (*SellerDetails, error) {
	seller, err := grpcclient.GetSellerDetails(context.Background(), &user.GetSellerDetailsRequest{Id: int64(id)})
	if err != nil {
		return nil, err
	}

	sellerDetails := SellerDetails{
		BusinessName: seller.BusinessName,
	}

	return &sellerDetails, nil
}
