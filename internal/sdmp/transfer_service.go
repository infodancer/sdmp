package sdmp

import (
	"context"

	pb "github.com/infodancer/sdmp/proto/sdmp/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TransferServer implements the SDMP TransferService gRPC interface.
type TransferServer struct {
	pb.UnimplementedTransferServiceServer
}

func (s *TransferServer) AcceptResponsibility(_ context.Context, req *pb.AcceptResponsibilityRequest) (*pb.AcceptResponsibilityResponse, error) {
	if len(req.GetMessageId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "message_id is required")
	}
	if req.GetRecipientDomain() == "" {
		return nil, status.Error(codes.InvalidArgument, "recipient_domain is required")
	}
	if req.GetGatewayDomain() == "" {
		return nil, status.Error(codes.InvalidArgument, "gateway_domain is required")
	}
	return nil, status.Error(codes.Unimplemented, "AcceptResponsibility not yet implemented")
}
