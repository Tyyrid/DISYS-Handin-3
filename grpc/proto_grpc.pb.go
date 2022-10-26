// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.7
// source: grpc/proto.proto

package proto

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

// TimeAskClient is the client API for TimeAsk service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TimeAskClient interface {
	SendMessage(ctx context.Context, in *ClientPublishMessage, opts ...grpc.CallOption) (*ServerPublishMessageOk, error)
	ConnectToServer(ctx context.Context, in *ClientConnectMessage, opts ...grpc.CallOption) (TimeAsk_ConnectToServerClient, error)
}

type timeAskClient struct {
	cc grpc.ClientConnInterface
}

func NewTimeAskClient(cc grpc.ClientConnInterface) TimeAskClient {
	return &timeAskClient{cc}
}

func (c *timeAskClient) SendMessage(ctx context.Context, in *ClientPublishMessage, opts ...grpc.CallOption) (*ServerPublishMessageOk, error) {
	out := new(ServerPublishMessageOk)
	err := c.cc.Invoke(ctx, "/simpleGuide.TimeAsk/SendMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *timeAskClient) ConnectToServer(ctx context.Context, in *ClientConnectMessage, opts ...grpc.CallOption) (TimeAsk_ConnectToServerClient, error) {
	stream, err := c.cc.NewStream(ctx, &TimeAsk_ServiceDesc.Streams[0], "/simpleGuide.TimeAsk/ConnectToServer", opts...)
	if err != nil {
		return nil, err
	}
	x := &timeAskConnectToServerClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type TimeAsk_ConnectToServerClient interface {
	Recv() (*MessageStreamConnection, error)
	grpc.ClientStream
}

type timeAskConnectToServerClient struct {
	grpc.ClientStream
}

func (x *timeAskConnectToServerClient) Recv() (*MessageStreamConnection, error) {
	m := new(MessageStreamConnection)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TimeAskServer is the server API for TimeAsk service.
// All implementations must embed UnimplementedTimeAskServer
// for forward compatibility
type TimeAskServer interface {
	SendMessage(context.Context, *ClientPublishMessage) (*ServerPublishMessageOk, error)
	ConnectToServer(*ClientConnectMessage, TimeAsk_ConnectToServerServer) error
	mustEmbedUnimplementedTimeAskServer()
}

// UnimplementedTimeAskServer must be embedded to have forward compatible implementations.
type UnimplementedTimeAskServer struct {
}

func (UnimplementedTimeAskServer) SendMessage(context.Context, *ClientPublishMessage) (*ServerPublishMessageOk, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
}
func (UnimplementedTimeAskServer) ConnectToServer(*ClientConnectMessage, TimeAsk_ConnectToServerServer) error {
	return status.Errorf(codes.Unimplemented, "method ConnectToServer not implemented")
}
func (UnimplementedTimeAskServer) mustEmbedUnimplementedTimeAskServer() {}

// UnsafeTimeAskServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TimeAskServer will
// result in compilation errors.
type UnsafeTimeAskServer interface {
	mustEmbedUnimplementedTimeAskServer()
}

func RegisterTimeAskServer(s grpc.ServiceRegistrar, srv TimeAskServer) {
	s.RegisterService(&TimeAsk_ServiceDesc, srv)
}

func _TimeAsk_SendMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClientPublishMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TimeAskServer).SendMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/simpleGuide.TimeAsk/SendMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TimeAskServer).SendMessage(ctx, req.(*ClientPublishMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _TimeAsk_ConnectToServer_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ClientConnectMessage)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TimeAskServer).ConnectToServer(m, &timeAskConnectToServerServer{stream})
}

type TimeAsk_ConnectToServerServer interface {
	Send(*MessageStreamConnection) error
	grpc.ServerStream
}

type timeAskConnectToServerServer struct {
	grpc.ServerStream
}

func (x *timeAskConnectToServerServer) Send(m *MessageStreamConnection) error {
	return x.ServerStream.SendMsg(m)
}

// TimeAsk_ServiceDesc is the grpc.ServiceDesc for TimeAsk service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TimeAsk_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "simpleGuide.TimeAsk",
	HandlerType: (*TimeAskServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendMessage",
			Handler:    _TimeAsk_SendMessage_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ConnectToServer",
			Handler:       _TimeAsk_ConnectToServer_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "grpc/proto.proto",
}
