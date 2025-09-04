package service

import (
	"context"
	grpcclient "products/grpc_client"
	"utils/user"
)

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
