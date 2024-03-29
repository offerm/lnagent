// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.17.1
// source: lncoordinator.proto

package protobuf

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// CoordinatorClient is the client API for Coordinator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CoordinatorClient interface {
	StatusUpdate(ctx context.Context, in *StatusUpdateRequest, opts ...grpc.CallOption) (*StatusUpdateResponse, error)
	Tasks(ctx context.Context, opts ...grpc.CallOption) (Coordinator_TasksClient, error)
}

type coordinatorClient struct {
	cc grpc.ClientConnInterface
}

func NewCoordinatorClient(cc grpc.ClientConnInterface) CoordinatorClient {
	return &coordinatorClient{cc}
}

func (c *coordinatorClient) StatusUpdate(ctx context.Context, in *StatusUpdateRequest, opts ...grpc.CallOption) (*StatusUpdateResponse, error) {
	out := new(StatusUpdateResponse)
	err := c.cc.Invoke(ctx, "/protobuf.coordinator/StatusUpdate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *coordinatorClient) Tasks(ctx context.Context, opts ...grpc.CallOption) (Coordinator_TasksClient, error) {
	stream, err := c.cc.NewStream(ctx, &Coordinator_ServiceDesc.Streams[0], "/protobuf.coordinator/Tasks", opts...)
	if err != nil {
		return nil, err
	}
	x := &coordinatorTasksClient{stream}
	return x, nil
}

type Coordinator_TasksClient interface {
	Send(*TaskResponse) error
	Recv() (*Task, error)
	grpc.ClientStream
}

type coordinatorTasksClient struct {
	grpc.ClientStream
}

func (x *coordinatorTasksClient) Send(m *TaskResponse) error {
	return x.ClientStream.SendMsg(m)
}

func (x *coordinatorTasksClient) Recv() (*Task, error) {
	m := new(Task)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// CoordinatorServer is the server API for Coordinator service.
// All implementations must embed UnimplementedCoordinatorServer
// for forward compatibility
type CoordinatorServer interface {
	StatusUpdate(context.Context, *StatusUpdateRequest) (*StatusUpdateResponse, error)
	Tasks(Coordinator_TasksServer) error
	mustEmbedUnimplementedCoordinatorServer()
}

// UnimplementedCoordinatorServer must be embedded to have forward compatible implementations.
type UnimplementedCoordinatorServer struct {
}

func (UnimplementedCoordinatorServer) StatusUpdate(context.Context, *StatusUpdateRequest) (*StatusUpdateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StatusUpdate not implemented")
}
func (UnimplementedCoordinatorServer) Tasks(Coordinator_TasksServer) error {
	return status.Errorf(codes.Unimplemented, "method Tasks not implemented")
}
func (UnimplementedCoordinatorServer) mustEmbedUnimplementedCoordinatorServer() {}

// UnsafeCoordinatorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CoordinatorServer will
// result in compilation errors.
type UnsafeCoordinatorServer interface {
	mustEmbedUnimplementedCoordinatorServer()
}

func RegisterCoordinatorServer(s grpc.ServiceRegistrar, srv CoordinatorServer) {
	s.RegisterService(&Coordinator_ServiceDesc, srv)
}

func _Coordinator_StatusUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatusUpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CoordinatorServer).StatusUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.coordinator/StatusUpdate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CoordinatorServer).StatusUpdate(ctx, req.(*StatusUpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Coordinator_Tasks_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(CoordinatorServer).Tasks(&coordinatorTasksServer{stream})
}

type Coordinator_TasksServer interface {
	Send(*Task) error
	Recv() (*TaskResponse, error)
	grpc.ServerStream
}

type coordinatorTasksServer struct {
	grpc.ServerStream
}

func (x *coordinatorTasksServer) Send(m *Task) error {
	return x.ServerStream.SendMsg(m)
}

func (x *coordinatorTasksServer) Recv() (*TaskResponse, error) {
	m := new(TaskResponse)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Coordinator_ServiceDesc is the grpc.ServiceDesc for Coordinator service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Coordinator_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protobuf.coordinator",
	HandlerType: (*CoordinatorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StatusUpdate",
			Handler:    _Coordinator_StatusUpdate_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Tasks",
			Handler:       _Coordinator_Tasks_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "lncoordinator.proto",
}
