// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.2
// source: proto/mensajes.proto

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

// NotificacionClient is the client API for Notificacion service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NotificacionClient interface {
	Notificar(ctx context.Context, in *NotiReq, opts ...grpc.CallOption) (*NotiResp, error)
	Inscribir(ctx context.Context, in *InscritosReq, opts ...grpc.CallOption) (*InscritosResp, error)
}

type notificacionClient struct {
	cc grpc.ClientConnInterface
}

func NewNotificacionClient(cc grpc.ClientConnInterface) NotificacionClient {
	return &notificacionClient{cc}
}

func (c *notificacionClient) Notificar(ctx context.Context, in *NotiReq, opts ...grpc.CallOption) (*NotiResp, error) {
	out := new(NotiResp)
	err := c.cc.Invoke(ctx, "/grpc.Notificacion/Notificar", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificacionClient) Inscribir(ctx context.Context, in *InscritosReq, opts ...grpc.CallOption) (*InscritosResp, error) {
	out := new(InscritosResp)
	err := c.cc.Invoke(ctx, "/grpc.Notificacion/Inscribir", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NotificacionServer is the server API for Notificacion service.
// All implementations must embed UnimplementedNotificacionServer
// for forward compatibility
type NotificacionServer interface {
	Notificar(context.Context, *NotiReq) (*NotiResp, error)
	Inscribir(context.Context, *InscritosReq) (*InscritosResp, error)
	mustEmbedUnimplementedNotificacionServer()
}

// UnimplementedNotificacionServer must be embedded to have forward compatible implementations.
type UnimplementedNotificacionServer struct {
}

func (UnimplementedNotificacionServer) Notificar(context.Context, *NotiReq) (*NotiResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Notificar not implemented")
}
func (UnimplementedNotificacionServer) Inscribir(context.Context, *InscritosReq) (*InscritosResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Inscribir not implemented")
}
func (UnimplementedNotificacionServer) mustEmbedUnimplementedNotificacionServer() {}

// UnsafeNotificacionServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NotificacionServer will
// result in compilation errors.
type UnsafeNotificacionServer interface {
	mustEmbedUnimplementedNotificacionServer()
}

func RegisterNotificacionServer(s grpc.ServiceRegistrar, srv NotificacionServer) {
	s.RegisterService(&Notificacion_ServiceDesc, srv)
}

func _Notificacion_Notificar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NotiReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificacionServer).Notificar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Notificacion/Notificar",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificacionServer).Notificar(ctx, req.(*NotiReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Notificacion_Inscribir_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InscritosReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificacionServer).Inscribir(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Notificacion/Inscribir",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificacionServer).Inscribir(ctx, req.(*InscritosReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Notificacion_ServiceDesc is the grpc.ServiceDesc for Notificacion service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Notificacion_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.Notificacion",
	HandlerType: (*NotificacionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Notificar",
			Handler:    _Notificacion_Notificar_Handler,
		},
		{
			MethodName: "Inscribir",
			Handler:    _Notificacion_Inscribir_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/mensajes.proto",
}
