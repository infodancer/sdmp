package sdmp

import (
	"context"

	pb "github.com/infodancer/sdmp/proto/sdmp/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// KeyServer implements the SDMP KeyService gRPC interface.
type KeyServer struct {
	pb.UnimplementedKeyServiceServer
}

func (s *KeyServer) GetDomainKey(_ context.Context, req *pb.GetDomainKeyRequest) (*pb.DomainKeyResponse, error) {
	if req.GetDomain() == "" {
		return nil, status.Error(codes.InvalidArgument, "domain is required")
	}
	return nil, status.Error(codes.Unimplemented, "GetDomainKey not yet implemented")
}

func (s *KeyServer) GetUserKey(_ context.Context, req *pb.GetUserKeyRequest) (*pb.UserKeyResponse, error) {
	if req.GetAddress() == "" {
		return nil, status.Error(codes.InvalidArgument, "address is required")
	}
	return nil, status.Error(codes.Unimplemented, "GetUserKey not yet implemented")
}
