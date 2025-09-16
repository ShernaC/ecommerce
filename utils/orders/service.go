package orders

import (
	"os"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ConnectionOption struct {
	Host *string
}

func Connect(connectionOption ConnectionOption) (OrderClient, *grpc.ClientConn) {
	defaultGRPCAddress := os.Getenv("ORDER_GRPC")

	grpcAddress := connectionOption.Host
	if grpcAddress == nil {
		grpcAddress = &defaultGRPCAddress
	}

	conn, err := grpc.Dial(*grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	return NewOrderClient(conn), conn
}
