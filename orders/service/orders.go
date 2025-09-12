package service

import (
	"context"
	"fmt"
	grpcclient "orders/grpc_client"
	"orders/middleware"
	"orders/model"
	"products/tools"
	"utils/user"
)

type OrderStatus string

const (
	ORDER_STATUS_PENDING   OrderStatus = "pending"
	ORDER_STATUS_PAID      OrderStatus = "paid"
	ORDER_STATUS_SHIPPED   OrderStatus = "shipped"
	ORDER_STATUS_CANCELLED OrderStatus = "cancelled"
	ORDER_STATUS_COMPLETED OrderStatus = "completed"
)

type PaymentMethod string

const (
	PAYMENT_METHOD_COD  PaymentMethod = "cash_on_delivery"
	PAYMENT_METHOD_CARD PaymentMethod = "credit_card"
)

func (s *Service) CreateOrder(ctx context.Context, cart model.Cart, paymentMethod string) (*model.Order, error) {
	var (
		totalAmount float64 = 0
		ctxData             = middleware.AuthContext(ctx)
	)

	if ctxData == nil || ctxData.ID != cart.UserID {
		return nil, fmt.Errorf("user is unauthorized to create order")
	}

	valid, err := s.OrderOnCreate(ctx, cart, paymentMethod)
	if err != nil || !valid {
		return nil, fmt.Errorf("invalid order")
	}

	for _, item := range cart.Items {
		totalAmount += item.Price * float64(item.Quantity)
	}

	// grpc call
	userDetails, err := grpcclient.GetUserDetails(ctx, &user.GetUserDetailsRequest{Id: int64(cart.UserID)})
	if err != nil {
		return nil, err
	}

	order := model.Order{
		UserID:          ctxData.ID,
		Status:          string(ORDER_STATUS_PENDING),
		TotalAmount:     totalAmount,
		ShippingAddress: userDetails.Address,
		PaymentMethod:   paymentMethod,
	}

	if err := s.DB.Create(&order).Error; err != nil {
		return nil, err
	}

	success, err := s.OrderAddItems(ctx, order, cart.Items)
	if err != nil {
		return nil, err
	} else if !success {
		return nil, fmt.Errorf("failed to add items to order")
	}

	return &order, nil
}

func (s *Service) OrderOnCreate(ctx context.Context, cart model.Cart, paymentMethod string) (bool, error) {
	if paymentMethod != string(PAYMENT_METHOD_COD) && paymentMethod != string(PAYMENT_METHOD_CARD) {
		return false, fmt.Errorf("invalid payment method")
	}

	if len(cart.Items) == 0 {
		return false, fmt.Errorf("cart is empty")
	}

	return true, nil
}

func (s *Service) OrderAddItems(ctx context.Context, order model.Order, items []*model.CartItem) (bool, error) {
	var (
		orderItems []*model.OrderItem
	)

	if len(items) == 0 {
		return false, fmt.Errorf("no items to add")
	}

	orderExist, err := s.CheckOrderExists(ctx, order.ID)
	if err != nil {
		return false, err
	}
	if !orderExist {
		return false, fmt.Errorf("order does not exist")
	}

	for _, item := range items {
		newOrderItem := model.NewOrderItem{
			OrderID:         order.ID,
			ProductID:       item.ProductID,
			Quantity:        item.Quantity,
			PriceAtPurchase: item.Price,
			ProductSnapshot: "",
		}

		orderItem, err := s.CreateOrderItem(ctx, newOrderItem)
		if err != nil {
			return false, err
		}

		orderItems = append(orderItems, orderItem)
	}

	order.Items = orderItems

	return true, nil
}

func (s *Service) CreateOrderItem(ctx context.Context, item model.NewOrderItem) (*model.OrderItem, error) {
	if item.OrderID <= 0 || item.ProductID <= 0 || item.Quantity <= 0 || item.PriceAtPurchase <= 0 || item.ProductSnapshot == "" {
		return nil, fmt.Errorf("data cannot be empty")
	}

	orderItem := model.OrderItem{
		OrderID:         item.OrderID,
		ProductID:       item.ProductID,
		Quantity:        item.Quantity,
		PriceAtPurchase: item.PriceAtPurchase,
		ProductSnapshot: item.ProductSnapshot,
	}

	if err := s.DB.Model(&model.OrderItem{}).Create(&orderItem).Error; err != nil {
		return nil, err
	}

	return &orderItem, nil
}

func (s *Service) CheckOrderExists(ctx context.Context, orderID int) (bool, error) {
	var count int64
	if err := s.DB.Model(&model.Order{}).Where("id = ?", orderID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

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
