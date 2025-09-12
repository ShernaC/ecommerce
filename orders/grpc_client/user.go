package grpcclient

import (
	"context"
	"utils/user"
)

func GetUserDetails(ctx context.Context, id *user.GetUserDetailsRequest) (*user.GetUserDetailsResponse, error) {
	userConn, conn := user.Connect(user.ConnectionOption{})
	defer conn.Close()

	userDetails, err := userConn.GetUserDetails(ctx, id)
	if err != nil {
		return nil, err
	}

	return userDetails, nil
}
