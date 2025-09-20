package resolver

import (
	"context"
	"orders/service"
	"utils/orders"
)

type Server struct {
	orders.UnimplementedOrderServer
}

func (s *Server) CreateCart(ctx context.Context, req *orders.CreateCartRequest) (*orders.CreateCartResponse, error) {
	userID := req.UserId

	tx := service.GetTransaction()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback(r)
			panic(r)
		}
	}()

	success, err := tx.CreateCart(ctx, int(userID))
	if err != nil {
		return nil, err
	}

	tx.Commit()

	return &orders.CreateCartResponse{
		Success: success,
	}, nil

}
