package grpcclient

import (
	"context"
	"utils/user"
)

func CheckSellerExists(ctx context.Context, id *user.CheckSellerExistsRequest) (*user.CheckSellerExistsResponse, error) {
	userConn, conn := user.Connect(user.ConnectionOption{})
	defer conn.Close()

	sellerExist, err := userConn.CheckSellerExists(ctx, id)
	if err != nil {
		panic(err)
	}

	return sellerExist, nil
}
