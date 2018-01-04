package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	uc "github.com/egymgmbh/mobile-update-check/pb"
)

type ucs struct {
}

var _ uc.UpdateCheckServer = &ucs{}

func (u *ucs) VersionCheck(ctx context.Context, in *uc.UpdateVersion) (*uc.UpdateAction, error) {
	return &uc.UpdateAction{
		Action: uc.MyAction_ADVICE,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	uc.RegisterUpdateCheckServer(grpcServer, &ucs{})
	grpcServer.Serve(lis)
}
