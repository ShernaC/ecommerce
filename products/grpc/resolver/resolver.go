package resolver

import (
	"context"
	"products/service"
	"utils/product"
)

type Server struct {
	product.UnimplementedProductServer
}

func (s Server) GetProductDetails(ctx context.Context, req *product.GetProductDetailsRequest) (*product.GetProductDetailsResponse, error) {
	productID := req.Id

	productDetail, err := service.GetService().ProductGetByID(ctx, int(productID))
	if err != nil {
		return nil, err
	}

	return &product.GetProductDetailsResponse{
		Id:          int64(productDetail.ID),
		SellerId:    int64(productDetail.SellerID),
		Name:        productDetail.Name,
		Description: productDetail.Description,
		Price:       productDetail.Price,
		Stock:       int64(productDetail.Stock),
	}, nil
}

func (s Server) UpdateStock(ctx context.Context, req *product.UpdateStockRequest) (*product.UpdateStockResponse, error) {
	ID := req.Id
	qty := req.QtyBought

	tx := service.GetTransaction()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback(r)
			panic(r)
		}
	}()

	success, err := tx.ProductUpdateStock(ctx, int(ID), int(qty))
	if err != nil {
		return nil, err
	}

	tx.Commit()

	return &product.UpdateStockResponse{
		Success: success,
	}, nil
}
