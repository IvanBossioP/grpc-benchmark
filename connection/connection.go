package connection

import (
	"crypto/tls"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	pb "grpc-benchmark/protobuf"
	"log"
	"net/url"
)

func InitGeyserClient(address string) pb.GeyserClient {

	uri, err := url.Parse(address)

	if err != nil {
		log.Fatalf("Failed to parse URL %s, error: %v", address, err)
	}

	var transportCredentials grpc.DialOption

	if uri.Scheme == "https" {
		transportCredentials = grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{}))
	} else {
		transportCredentials = grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	conn, err := grpc.NewClient(uri.Host, transportCredentials)

	if err != nil {
		log.Fatalf("Failed to connect to %s server: %v", address, err)
	}

	client := pb.NewGeyserClient(conn)

	return client
}
