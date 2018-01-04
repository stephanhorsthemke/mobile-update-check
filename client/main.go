package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	pb "github.com/egymgmbh/mobile-update-check/pb"
	"google.golang.org/grpc"
)

func main() {
	serverAddress := flag.String("server-address", "localhost:8080", "Address of update check server to query.")
	os := flag.String("os", pb.OSType_name[0], "Operating system name.")
	osVersion := flag.String("os-version", "1.0.0", "Operating system version.")
	product := flag.String("product", pb.ProductType_name[0], "Product name.")
	productVersion := flag.String("product-version", "1.0.0", "Product version.")
	flag.Parse()

	// input validation
	if _, ok := pb.OSType_value[*os]; !ok {
		log.Fatalf("bad os name: %v", *os)
	}
	if _, ok := pb.ProductType_value[*product]; !ok {
		log.Fatalf("bad product name: %v", *product)
	}

	conn, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	pbClient := pb.NewUpdateCheckServiceClient(conn)

	resp, err := pbClient.Query(context.Background(), &pb.UpdateCheckRequest{
		OS: (pb.OSType) (pb.OSType_value[*os]),
		OSVersion: *osVersion,
		Product: (pb.ProductType) (pb.ProductType_value[*product]),
		ProductVersion: *productVersion,
	})
	if err != nil {
		log.Fatalf("version check: %v", err)
	}
	fmt.Printf("%v\n", resp.Action.String())
}
