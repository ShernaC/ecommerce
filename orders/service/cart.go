package service

import (
	"context"
	"fmt"
	"orders/middleware"
	"orders/model"

	"gorm.io/gorm"
)

func (s *Service) CreateCart(ctx context.Context, userID int) (bool, error) {
	if userID <= 0 {
		return false, fmt.Errorf("invalid user id")
	}

	var cart = model.Cart{
		UserID: userID,
	}

	if err := s.DB.Model(&cart).Create(&cart).Error; err != nil {
		return false, err
	}

	return true, nil
}

func (s *Service) AddToCart(ctx context.Context, newItem model.CartItemInput) (bool, error) {
	var (
		user = middleware.AuthContext(ctx)
	)

	if newItem.ProductID <= 0 || newItem.Quantity <= 0 {
		return false, fmt.Errorf("invalid product ID or quantity")
	}

	cart, err := s.CartGetDetails(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get cart: %w", err)
	}

	if user.ID == 0 || user.ID != cart.UserID {
		return false, fmt.Errorf("invalid user")
	}

	// gRPC call to get product details
	product, err := s.GetProductDetails(ctx, newItem.ProductID)
	if err != nil {
		return false, fmt.Errorf("failed to get product details: %w", err)
	}

	// validate product
	if newItem.Quantity > int(product.Stock) {
		return false, fmt.Errorf("item is out of stock")
	}

	newItemDetails := model.NewCartItem{
		CartID:    cart.ID,
		ProductID: int(product.ID),
		Quantity:  newItem.Quantity,
		Price:     product.Price,
	}

	item, err := s.CartCreateItem(ctx, newItemDetails)
	if err != nil {
		panic(err)
	}

	cart.Items = append(cart.Items, item)

	return true, nil
}

func (s *Service) CartGetDetails(ctx context.Context) (*model.Cart, error) {
	var (
		cart *model.Cart
		user = middleware.AuthContext(ctx)
	)

	if user.ID == 0 {
		return nil, fmt.Errorf("unauthorised user")
	}

	if err := s.DB.Model(&cart).Where("user_id = ?", user.ID).First(&cart).Error; err != nil {
		return nil, err
	}

	cartItems, err := s.CartGetItemsByCartID(ctx, cart.ID)
	if cartItems == nil {
		cartItems = []*model.CartItem{}
	} else if err != nil {
		return nil, err
	}

	cart.Items = cartItems

	return cart, nil
}

func (s *Service) CartGetItemDetails(ctx context.Context, cartID int, itemID int) (*model.CartItem, error) {
	var (
		cartItem *model.CartItem
		user     = middleware.AuthContext(ctx)
	)

	if user.ID == 0 {
		return nil, fmt.Errorf("unauthorised user")
	}

	if err := s.DB.Model(&cartItem).Where("cart_id = ? AND id = ?", cartID, itemID).Find(&cartItem).Error; err != nil {
		return nil, err
	}

	return cartItem, nil
}

func (s *Service) CartGetItemsByCartID(ctx context.Context, cartID int) ([]*model.CartItem, error) {
	var (
		cartItems []*model.CartItem
		user      = middleware.AuthContext(ctx)
	)

	if user.ID == 0 {
		return nil, fmt.Errorf("unauthorised user")
	}

	err := s.DB.Model(&cartItems).Where("cart_id = ?", cartID).Scan(&cartItems).Error
	if err == gorm.ErrRecordNotFound {
		return []*model.CartItem{}, nil
	} else if err != nil {
		return []*model.CartItem{}, err
	}

	return cartItems, nil
}

func (s *Service) CartGetItemsByIDs(ctx context.Context, cartItemIDs []int) ([]*model.CartItem, error) {
	var (
		cartItems []*model.CartItem
	)

	err := s.DB.Model(&cartItems).Where("id IN (?)", cartItemIDs).Scan(&cartItems).Error
	if err == gorm.ErrRecordNotFound {
		return []*model.CartItem{}, nil
	} else if err != nil {
		return []*model.CartItem{}, err
	}

	return cartItems, nil
}

func (s *Service) CartCreateItem(ctx context.Context, newItem model.NewCartItem) (*model.CartItem, error) {
	// if newItem.CartID <= 0 || newItem.ProductID <= 0 || newItem.Quantity <= 0 || newItem.Price <= 0 {
	// 	return nil, fmt.Errorf("input for item details cannot be empty")
	// }

	item := model.CartItem{
		CartID:    newItem.CartID,
		ProductID: newItem.ProductID,
		Quantity:  newItem.Quantity,
		Price:     newItem.Price,
	}

	if err := s.DB.Model(&item).Create(&item).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

func (s *Service) CartUpdateItem(ctx context.Context, itemDetails model.EditCartItem) (bool, error) {
	if itemDetails.CartID <= 0 || itemDetails.ProductID <= 0 {
		return false, fmt.Errorf("cart ID and item ID cannot be empty")
	}

	product, err := s.GetProductDetails(ctx, itemDetails.ProductID)
	if err != nil {
		return false, err
	}

	price := float64(itemDetails.Quantity) * product.Price

	if err := s.DB.Table("cart_item").Where("id = ?", itemDetails.ID).UpdateColumns(model.CartItem{Quantity: itemDetails.Quantity, Price: price}).Error; err != nil {
		return false, err
	}

	return true, nil

}

func (s *Service) CartRemoveItems(ctx context.Context, cartItemIDs []int) (bool, error) {
	if len(cartItemIDs) <= 0 {
		return false, fmt.Errorf("no selected items")
	}

	if err := s.DB.Table("cart_item").Delete(&model.CartItem{}, cartItemIDs).Error; err != nil {
		return false, err
	}

	return true, nil
}
