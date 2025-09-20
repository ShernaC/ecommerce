package service

import (
	"context"
	"fmt"
	grpcclient "orders/grpc_client"
	"utils/product"
)

type ProductDetail struct {
	ID          int
	SellerID    int
	Name        string
	Description string
	Price       float64
	Stock       int
	SKU         string
	ShopName    string
}

func (s *Service) GetProductDetails(ctx context.Context, id int) (*ProductDetail, error) {
	if id <= 0 {
		return nil, fmt.Errorf("product id is invalid")
	}

	product, err := grpcclient.GetProductDetails(ctx, &product.GetProductDetailsRequest{Id: int64(id)})
	if err != nil {
		return nil, err
	}

	productDetails := ProductDetail{
		ID:          int(product.Id),
		SellerID:    int(product.SellerId),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       int(product.Stock),
		SKU:         product.Sku,
		ShopName:    product.ShopName,
	}

	return &productDetails, nil
}

func (s *Service) UpdateStock(ctx context.Context, id int, qty int) (bool, error) {
	if id <= 0 || qty <= 0 {
		return false, fmt.Errorf("invalid input to update stock")
	}

	stockUpdated, err := grpcclient.UpdateStock(ctx, &product.UpdateStockRequest{Id: int64(id), QtyBought: int64(qty)})
	if err != nil {
		return false, err
	}

	return stockUpdated.Success, nil
}
