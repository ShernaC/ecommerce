package product

import (
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ConnectionOption struct {
	Host *string
}

func Connect(connectionOption ConnectionOption) (ProductClient, *grpc.ClientConn) {
	defaultGRPCAddress := os.Getenv("PRODUCT_GRPC")

	grpcAddress := connectionOption.Host
	if grpcAddress == nil {
		grpcAddress = &defaultGRPCAddress
	}

	conn, err := grpc.Dial(*grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	return NewProductClient(conn), conn
}
