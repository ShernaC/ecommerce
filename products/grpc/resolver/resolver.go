package resolver

import (
	"context"
	"products/service"
	"utils/product"
)

type Server struct{}

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
