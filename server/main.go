package main

import (
	"context"
	"log"
	"net"
	"fmt"

	"google.golang.org/grpc"

	uc "github.com/egymgmbh/mobile-update-check/pb"
	"github.com/Masterminds/semver"
)


type rule struct {
	os uc.OSType
	osVersion string
	osVersionConstraints *semver.Constraints
	product uc.ProductType
	productVersion string
	productVersionConstraints *semver.Constraints
	action uc.ResponseAction
}

type ucs struct {
}

var (
	_ uc.UpdateCheckServiceServer = &ucs{}
 	rules []rule
 )

func init() {
	rules := []rule{
		{
				os: uc.OSType_ANDROID,
				osVersion: ">=8.0.0, <9.0.0",
				product: uc.ProductType_FITAPP,
				productVersion: "<=2.3.0",
				action: uc.ResponseAction_ADVICE,
		},
		{
				os: uc.OSType_IOS,
				osVersion: "9.0.0",
				product: uc.ProductType_TRAINERAPP,
				productVersion: "2.3.0",
				action: uc.ResponseAction_FORCE,
		},
	}

	// compile contraints
	for _, rule := range rules {
		err := rule.CompileConstraints()
		if err != nil {
			log.Fatalf("compile constraints: %v", err)
		}
	}
}

func (r *rule) CompileConstraints() error {
		c, err := semver.NewConstraint(r.osVersion)
		if err != nil {
			return err
		}
		r.osVersionConstraints = c

		c, err = semver.NewConstraint(r.productVersion)
		if err != nil {
			return err
		}
		r.productVersionConstraints = c

		return nil
}

func (r *rule) Apply(req *uc.UpdateCheckRequest, osVersion, productVersion *semver.Version) (uc.ResponseAction, bool) {
	if req.OS == r.os &&
	 req.Product == r.product  &&
	 r.osVersionConstraints.Check(osVersion) &&
	 r.productVersionConstraints.Check(productVersion) {
		return r.action, true
	}
	return uc.ResponseAction_NONE, false
}

func (u *ucs) Query(ctx context.Context, in *uc.UpdateCheckRequest) (*uc.UpdateCheckResponse, error) {
	osVersion, err := semver.NewVersion(in.OSVersion)
	if err != nil {
		return nil, fmt.Errorf("malformed request")
	}
	productVersion, err := semver.NewVersion(in.ProductVersion)
	if err != nil {
		return nil, fmt.Errorf("malformed request")
	}
	for _, rule := range rules {
		if action, ok := rule.Apply(in, osVersion, productVersion); ok {
			return &uc.UpdateCheckResponse{
				Action: action,
			}, nil
		}
	}
	return &uc.UpdateCheckResponse{
		Action: uc.ResponseAction_NONE,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	uc.RegisterUpdateCheckServiceServer(grpcServer, &ucs{})
	grpcServer.Serve(lis)
}
