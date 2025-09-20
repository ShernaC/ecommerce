package grpcclient

import (
	"context"
	"utils/orders"
)

func CreateCart(ctx context.Context, req *orders.CreateCartRequest) (*orders.CreateCartResponse, error) {
	orderConn, conn := orders.Connect(orders.ConnectionOption{})
	defer conn.Close()

	cartCreated, err := orderConn.CreateCart(ctx, req)
	if err != nil {
		return nil, err
	}

	return cartCreated, nil
}
