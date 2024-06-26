// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.3
// source: gateway.proto

package gatewaypb

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

const (
	ProfileService_CreateProfile_FullMethodName  = "/proto.ProfileService/CreateProfile"
	ProfileService_ReadProfile_FullMethodName    = "/proto.ProfileService/ReadProfile"
	ProfileService_UpdateProfile_FullMethodName  = "/proto.ProfileService/UpdateProfile"
	ProfileService_GetNextProfile_FullMethodName = "/proto.ProfileService/GetNextProfile"
	ProfileService_Like_FullMethodName           = "/proto.ProfileService/Like"
	ProfileService_GetLikes_FullMethodName       = "/proto.ProfileService/GetLikes"
)

// ProfileServiceClient is the client API for ProfileService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProfileServiceClient interface {
	CreateProfile(ctx context.Context, in *ProfileRequest, opts ...grpc.CallOption) (*ErrorResponse, error)
	ReadProfile(ctx context.Context, in *IdRequest, opts ...grpc.CallOption) (*Profile, error)
	UpdateProfile(ctx context.Context, in *ProfileRequest, opts ...grpc.CallOption) (*ErrorResponse, error)
	GetNextProfile(ctx context.Context, in *IdRequest, opts ...grpc.CallOption) (*Profile, error)
	Like(ctx context.Context, in *TargetRequest, opts ...grpc.CallOption) (*ErrorResponse, error)
	GetLikes(ctx context.Context, in *IdRequest, opts ...grpc.CallOption) (*LikesResponse, error)
}

type profileServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewProfileServiceClient(cc grpc.ClientConnInterface) ProfileServiceClient {
	return &profileServiceClient{cc}
}

func (c *profileServiceClient) CreateProfile(ctx context.Context, in *ProfileRequest, opts ...grpc.CallOption) (*ErrorResponse, error) {
	out := new(ErrorResponse)
	err := c.cc.Invoke(ctx, ProfileService_CreateProfile_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileServiceClient) ReadProfile(ctx context.Context, in *IdRequest, opts ...grpc.CallOption) (*Profile, error) {
	out := new(Profile)
	err := c.cc.Invoke(ctx, ProfileService_ReadProfile_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileServiceClient) UpdateProfile(ctx context.Context, in *ProfileRequest, opts ...grpc.CallOption) (*ErrorResponse, error) {
	out := new(ErrorResponse)
	err := c.cc.Invoke(ctx, ProfileService_UpdateProfile_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileServiceClient) GetNextProfile(ctx context.Context, in *IdRequest, opts ...grpc.CallOption) (*Profile, error) {
	out := new(Profile)
	err := c.cc.Invoke(ctx, ProfileService_GetNextProfile_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileServiceClient) Like(ctx context.Context, in *TargetRequest, opts ...grpc.CallOption) (*ErrorResponse, error) {
	out := new(ErrorResponse)
	err := c.cc.Invoke(ctx, ProfileService_Like_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileServiceClient) GetLikes(ctx context.Context, in *IdRequest, opts ...grpc.CallOption) (*LikesResponse, error) {
	out := new(LikesResponse)
	err := c.cc.Invoke(ctx, ProfileService_GetLikes_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProfileServiceServer is the server API for ProfileService service.
// All implementations must embed UnimplementedProfileServiceServer
// for forward compatibility
type ProfileServiceServer interface {
	CreateProfile(context.Context, *ProfileRequest) (*ErrorResponse, error)
	ReadProfile(context.Context, *IdRequest) (*Profile, error)
	UpdateProfile(context.Context, *ProfileRequest) (*ErrorResponse, error)
	GetNextProfile(context.Context, *IdRequest) (*Profile, error)
	Like(context.Context, *TargetRequest) (*ErrorResponse, error)
	GetLikes(context.Context, *IdRequest) (*LikesResponse, error)
	mustEmbedUnimplementedProfileServiceServer()
}

// UnimplementedProfileServiceServer must be embedded to have forward compatible implementations.
type UnimplementedProfileServiceServer struct {
}

func (UnimplementedProfileServiceServer) CreateProfile(context.Context, *ProfileRequest) (*ErrorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateProfile not implemented")
}
func (UnimplementedProfileServiceServer) ReadProfile(context.Context, *IdRequest) (*Profile, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadProfile not implemented")
}
func (UnimplementedProfileServiceServer) UpdateProfile(context.Context, *ProfileRequest) (*ErrorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateProfile not implemented")
}
func (UnimplementedProfileServiceServer) GetNextProfile(context.Context, *IdRequest) (*Profile, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNextProfile not implemented")
}
func (UnimplementedProfileServiceServer) Like(context.Context, *TargetRequest) (*ErrorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Like not implemented")
}
func (UnimplementedProfileServiceServer) GetLikes(context.Context, *IdRequest) (*LikesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLikes not implemented")
}
func (UnimplementedProfileServiceServer) mustEmbedUnimplementedProfileServiceServer() {}

// UnsafeProfileServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProfileServiceServer will
// result in compilation errors.
type UnsafeProfileServiceServer interface {
	mustEmbedUnimplementedProfileServiceServer()
}

func RegisterProfileServiceServer(s grpc.ServiceRegistrar, srv ProfileServiceServer) {
	s.RegisterService(&ProfileService_ServiceDesc, srv)
}

func _ProfileService_CreateProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProfileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServiceServer).CreateProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProfileService_CreateProfile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServiceServer).CreateProfile(ctx, req.(*ProfileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProfileService_ReadProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServiceServer).ReadProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProfileService_ReadProfile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServiceServer).ReadProfile(ctx, req.(*IdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProfileService_UpdateProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProfileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServiceServer).UpdateProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProfileService_UpdateProfile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServiceServer).UpdateProfile(ctx, req.(*ProfileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProfileService_GetNextProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServiceServer).GetNextProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProfileService_GetNextProfile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServiceServer).GetNextProfile(ctx, req.(*IdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProfileService_Like_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TargetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServiceServer).Like(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProfileService_Like_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServiceServer).Like(ctx, req.(*TargetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProfileService_GetLikes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServiceServer).GetLikes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProfileService_GetLikes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServiceServer).GetLikes(ctx, req.(*IdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ProfileService_ServiceDesc is the grpc.ServiceDesc for ProfileService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ProfileService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.ProfileService",
	HandlerType: (*ProfileServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateProfile",
			Handler:    _ProfileService_CreateProfile_Handler,
		},
		{
			MethodName: "ReadProfile",
			Handler:    _ProfileService_ReadProfile_Handler,
		},
		{
			MethodName: "UpdateProfile",
			Handler:    _ProfileService_UpdateProfile_Handler,
		},
		{
			MethodName: "GetNextProfile",
			Handler:    _ProfileService_GetNextProfile_Handler,
		},
		{
			MethodName: "Like",
			Handler:    _ProfileService_Like_Handler,
		},
		{
			MethodName: "GetLikes",
			Handler:    _ProfileService_GetLikes_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gateway.proto",
}
