package grpcclient

import (
	"context"
	"utils/product"
)

func GetProductDetails(ctx context.Context, req *product.GetProductDetailsRequest) (*product.GetProductDetailsResponse, error) {
	productConn, conn := product.Connect(product.ConnectionOption{})
	defer conn.Close()

	productDetails, err := productConn.GetProductDetails(ctx, req)
	if err != nil {
		panic(err)
	}

	return productDetails, nil
}

func UpdateStock(ctx context.Context, req *product.UpdateStockRequest) (*product.UpdateStockResponse, error) {
	productConn, conn := product.Connect(product.ConnectionOption{})
	defer conn.Close()

	updateStock, err := productConn.UpdateStock(ctx, req)
	if err != nil {
		panic(err)
	}

	return updateStock, nil
}
