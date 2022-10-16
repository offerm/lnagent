package lnagent

import (
	"context"
	"github.com/offerm/lnagent/protobuf"
	"google.golang.org/grpc"
	"log"
	"net"
)

type CoordinatorServer struct {
	protobuf.UnimplementedCoordinatorServer
}

func (c CoordinatorServer) StatusUpdate(ctx context.Context, request *protobuf.StatusUpdateRequest) (*protobuf.StatusUpdateResponse, error) {
	//TODO implement me
	return &protobuf.StatusUpdateResponse{}, nil

	//print("fff")
	//panic("implement me")
}

func (c CoordinatorServer) Tasks(server protobuf.Coordinator_TasksServer) error {
	go func() {
		server.Send(&protobuf.Task{
			Type: &protobuf.Task_InitType{
				InitType: &protobuf.Task_Init{
					To: &protobuf.Payment{ //to B
						AmtMsat:  1000000,
						FeeMsat:  0,
						TimeLock: 0,
					},
					From: &protobuf.Payment{
						AmtMsat:  1000000,
						FeeMsat:  0,
						TimeLock: 0,
					},
				}}})
	}()
	for {
		err := server.Context().Err()
		if err != nil {
			return err
		}

		resp, recErr := server.Recv()
		if recErr != nil {
			return recErr
		} else {
			if resp != nil {
				break
			}
		}
	}
	return nil
}

func newMockCoordinator() {
	lis, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	protobuf.RegisterCoordinatorServer(grpcServer, CoordinatorServer{})
	grpcServer.Serve(lis)
}
