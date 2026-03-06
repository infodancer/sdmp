package sdmp

import (
	"context"

	pb "github.com/infodancer/sdmp/proto/sdmp/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// FetchServer implements the SDMP FetchService gRPC interface.
type FetchServer struct {
	pb.UnimplementedFetchServiceServer
}

func (s *FetchServer) FetchMessage(_ context.Context, req *pb.FetchMessageRequest) (*pb.FetchMessageResponse, error) {
	if len(req.GetMessageId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "message_id is required")
	}
	return nil, status.Error(codes.Unimplemented, "FetchMessage not yet implemented")
}

func (s *FetchServer) BatchFetchMessages(req *pb.BatchFetchMessagesRequest, stream pb.FetchService_BatchFetchMessagesServer) error {
	if len(req.GetMessageIds()) == 0 {
		return status.Error(codes.InvalidArgument, "at least one message_id is required")
	}
	return status.Error(codes.Unimplemented, "BatchFetchMessages not yet implemented")
}
