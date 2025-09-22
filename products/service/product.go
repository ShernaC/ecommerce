package service

import (
	"context"
	"fmt"
	"products/model"
	"utils/middleware"

	"products/tools"
	"time"

	"gorm.io/gorm"
)

func (s *Service) ProductCreate(ctx context.Context, newProd model.NewProduct) (*model.Product, error) {
	valid, err := s.ProductOnCreate(ctx, newProd)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, fmt.Errorf("error creating product")
	}

	seller, err := s.GetSellerDetails(newProd.SellerID)
	if err != nil {
		return nil, err
	}

	product := model.Product{
		Name:        newProd.Name,
		Description: newProd.Description,
		Price:       newProd.Price,
		Stock:       newProd.Stock,
		SellerID:    newProd.SellerID,
		ShopName:    seller.BusinessName,
	}

	if err := s.DB.Model(&product).Create(&product).Error; err != nil {
		return nil, err
	}

	sku := tools.GenerateSKU(&product)
	product.SKU = &sku
	s.DB.Save(&product)

	return &product, nil
}

func (s *Service) ProductOnCreate(ctx context.Context, newProd model.NewProduct) (bool, error) {
	if newProd.Name == "" || newProd.Description == "" {
		return false, fmt.Errorf("invalid input: fields cannot be empty")
	}

	if newProd.Price < 0 || newProd.Stock < 0 || newProd.SellerID <= 0 {
		return false, fmt.Errorf("invalid input: numerical inputs cannot be negative")
	}

	return true, nil
}

func (s *Service) ProductUpdate(ctx context.Context, prodUpdates model.UpdateProduct) (*model.Product, error) {
	var (
		ctxData = middleware.AuthContext(ctx)
	)

	valid, err := s.ProductCheckBelongToSeller(ctx, prodUpdates.ID, ctxData.ID)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, fmt.Errorf("failed to update product: product does not belong to seller")
	}

	if err := s.DB.Table("product").Scopes(tools.IsDeletedAtNull).Where("id = ?", prodUpdates.ID).Updates(map[string]interface{}{
		"name":        prodUpdates.Name,
		"description": prodUpdates.Description,
		"price":       prodUpdates.Price,
		"stock":       prodUpdates.Stock,
	}).Error; err != nil {
		return nil, err
	}

	return s.ProductGetByID(ctx, prodUpdates.ID)
}

func (s *Service) ProductDelete(ctx context.Context, id int) (string, error) {
	var (
		ctxData = middleware.AuthContext(ctx)
		timeNow = time.Now().Format("2006-01-02 15:04:05")
	)

	if id <= 0 {
		return "failed", fmt.Errorf("id cannot be empty or negative")
	}

	valid, err := s.ProductCheckBelongToSeller(ctx, id, ctxData.ID)
	if err != nil {
		return "", err
	}
	if !valid {
		return "", fmt.Errorf("failed to update product: product does not belong to seller")
	}

	if err := s.DB.Table("product").Where(tools.IsDeletedAtNull).Update("deleted_at = ?", timeNow).Error; err != nil {
		return "", err
	}

	return "success", nil
}

func (s *Service) ProductCheckBelongToSeller(ctx context.Context, id int, sellerID int) (bool, error) {
	if id == 0 || sellerID == 0 {
		return false, fmt.Errorf("id or seller id cannot be empty")
	}

	product, err := s.ProductGetByID(ctx, id)
	if err != nil {
		return false, err
	}

	return product.SellerID == sellerID, nil
}

func (s *Service) ProductGetByID(ctx context.Context, id int) (*model.Product, error) {
	var product *model.Product

	if err := s.DB.Model(&product).Where("id = ?", id).First(&product).Error; err != nil {
		return nil, err
	}

	return product, nil
}

func (s *Service) ProductGetProductsBySellerID(ctx context.Context, id int) ([]*model.Product, error) {
	var products []*model.Product

	if err := s.DB.Model(&products).Where("id = ?", id).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (s *Service) ProductGetWithFilter(ctx context.Context, filter string) ([]*model.Product, error) {
	var products []*model.Product

	if err := s.DB.Model(&products).Where("name LIKE %?%", filter).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (s *Service) ProductUpdateStock(ctx context.Context, id int, qty int) (bool, error) {
	var product *model.Product

	result := s.DB.Model(&product).Where("id = ?", id).Update("stock", gorm.Expr("stock - ?", qty))
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, fmt.Errorf("insufficient stock or stock not found")
	}

	return true, nil
}
