package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	uc "github.com/egymgmbh/mobile-update-check/pb"
	"google.golang.org/grpc"
)

func main() {
	serverAddress := flag.String("server-address", "localhost:8080", "Address of update check server to query.")
	flag.Parse()

	if len(flag.Args()) < 1 {
		log.Fatalf("missing version")
	}

	conn, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	ucClient := uc.NewUpdateCheckClient(conn)

	action, err := ucClient.VersionCheck(context.Background(), &uc.UpdateVersion{
		Version: flag.Args()[0],
	})
	if err != nil {
		log.Fatalf("version check: %v", err)
	}
	fmt.Printf("%v\n", action)

}
