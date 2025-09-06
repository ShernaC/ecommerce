package grpcclient

import (
	"context"
	"utils/product"
)

func GetProductDetails(ctx context.Context, id *product.GetProductDetailsRequest) (*product.GetProductDetailsResponse, error) {
	productConn, conn := product.Connect(product.ConnectionOption{})
	defer conn.Close()

	productDetails, err := productConn.GetProductDetails(ctx, id)
	if err != nil {
		panic(err)
	}

	return productDetails, nil
}
