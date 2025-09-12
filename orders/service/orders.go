package service

import (
	"context"
	"orders/middleware"
	"orders/model"
	"products/tools"
)

func (s *Service) OrderGetHistoryByUserID(ctx context.Context) ([]*model.Order, error) {
	var (
		orders  []*model.Order
		ctxData = middleware.AuthContext(ctx)
	)

	if err := s.DB.Model(&orders).Scopes(tools.IsDeletedAtNull).Where("user_id = ?", ctxData.ID).Scan(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}
