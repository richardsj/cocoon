// Code generated by protoc-gen-go.
// source: server.proto
// DO NOT EDIT!

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	server.proto

It has these top-level messages:
	Ok
	InvokeParam
	InvokeResponse
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

type Ok struct {
	Status int32 `protobuf:"varint,1,opt,name=status" json:"status,omitempty"`
}

func (m *Ok) Reset()                    { *m = Ok{} }
func (m *Ok) String() string            { return proto1.CompactTextString(m) }
func (*Ok) ProtoMessage()               {}
func (*Ok) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Ok) GetStatus() int32 {
	if m != nil {
		return m.Status
	}
	return 0
}

type InvokeParam struct {
	ID       string   `protobuf:"bytes,1,opt,name=ID" json:"ID,omitempty"`
	Function string   `protobuf:"bytes,2,opt,name=function" json:"function,omitempty"`
	Params   []string `protobuf:"bytes,3,rep,name=params" json:"params,omitempty"`
}

func (m *InvokeParam) Reset()                    { *m = InvokeParam{} }
func (m *InvokeParam) String() string            { return proto1.CompactTextString(m) }
func (*InvokeParam) ProtoMessage()               {}
func (*InvokeParam) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *InvokeParam) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *InvokeParam) GetFunction() string {
	if m != nil {
		return m.Function
	}
	return ""
}

func (m *InvokeParam) GetParams() []string {
	if m != nil {
		return m.Params
	}
	return nil
}

type InvokeResponse struct {
	ID     string `protobuf:"bytes,1,opt,name=ID" json:"ID,omitempty"`
	Status int32  `protobuf:"varint,2,opt,name=status" json:"status,omitempty"`
	Body   []byte `protobuf:"bytes,3,opt,name=body,proto3" json:"body,omitempty"`
}

func (m *InvokeResponse) Reset()                    { *m = InvokeResponse{} }
func (m *InvokeResponse) String() string            { return proto1.CompactTextString(m) }
func (*InvokeResponse) ProtoMessage()               {}
func (*InvokeResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *InvokeResponse) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *InvokeResponse) GetStatus() int32 {
	if m != nil {
		return m.Status
	}
	return 0
}

func (m *InvokeResponse) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

func init() {
	proto1.RegisterType((*Ok)(nil), "proto.Ok")
	proto1.RegisterType((*InvokeParam)(nil), "proto.InvokeParam")
	proto1.RegisterType((*InvokeResponse)(nil), "proto.InvokeResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Stub service

type StubClient interface {
	HealthCheck(ctx context.Context, in *Ok, opts ...grpc.CallOption) (*Ok, error)
	Invoke(ctx context.Context, in *InvokeParam, opts ...grpc.CallOption) (*InvokeResponse, error)
}

type stubClient struct {
	cc *grpc.ClientConn
}

func NewStubClient(cc *grpc.ClientConn) StubClient {
	return &stubClient{cc}
}

func (c *stubClient) HealthCheck(ctx context.Context, in *Ok, opts ...grpc.CallOption) (*Ok, error) {
	out := new(Ok)
	err := grpc.Invoke(ctx, "/proto.Stub/HealthCheck", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stubClient) Invoke(ctx context.Context, in *InvokeParam, opts ...grpc.CallOption) (*InvokeResponse, error) {
	out := new(InvokeResponse)
	err := grpc.Invoke(ctx, "/proto.Stub/Invoke", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Stub service

type StubServer interface {
	HealthCheck(context.Context, *Ok) (*Ok, error)
	Invoke(context.Context, *InvokeParam) (*InvokeResponse, error)
}

func RegisterStubServer(s *grpc.Server, srv StubServer) {
	s.RegisterService(&_Stub_serviceDesc, srv)
}

func _Stub_HealthCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Ok)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StubServer).HealthCheck(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Stub/HealthCheck",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StubServer).HealthCheck(ctx, req.(*Ok))
	}
	return interceptor(ctx, in, info, handler)
}

func _Stub_Invoke_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InvokeParam)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StubServer).Invoke(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Stub/Invoke",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StubServer).Invoke(ctx, req.(*InvokeParam))
	}
	return interceptor(ctx, in, info, handler)
}

var _Stub_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Stub",
	HandlerType: (*StubServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HealthCheck",
			Handler:    _Stub_HealthCheck_Handler,
		},
		{
			MethodName: "Invoke",
			Handler:    _Stub_Invoke_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "server.proto",
}

func init() { proto1.RegisterFile("server.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 222 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x64, 0x8f, 0xc1, 0x4e, 0x84, 0x30,
	0x18, 0x84, 0x43, 0xd9, 0x25, 0xf2, 0xef, 0x66, 0x0f, 0x7f, 0xa2, 0x69, 0x88, 0x07, 0x82, 0x17,
	0x4e, 0x7b, 0x90, 0x47, 0x90, 0x83, 0x24, 0x26, 0x68, 0x7d, 0x01, 0x0b, 0xd6, 0x60, 0xaa, 0x2d,
	0x69, 0x0b, 0x89, 0x6f, 0x6f, 0x28, 0x28, 0x18, 0x4f, 0x9d, 0xf9, 0xd2, 0x4e, 0x67, 0xe0, 0x68,
	0x85, 0x19, 0x85, 0x39, 0xf7, 0x46, 0x3b, 0x8d, 0x7b, 0x7f, 0x64, 0xd7, 0x40, 0x6a, 0x89, 0x57,
	0x10, 0x59, 0xc7, 0xdd, 0x60, 0x69, 0x90, 0x06, 0xf9, 0x9e, 0x2d, 0x2e, 0x7b, 0x82, 0x43, 0xa5,
	0x46, 0x2d, 0xc5, 0x23, 0x37, 0xfc, 0x13, 0x4f, 0x40, 0xaa, 0xd2, 0x5f, 0x89, 0x19, 0xa9, 0x4a,
	0x4c, 0xe0, 0xe2, 0x6d, 0x50, 0xad, 0x7b, 0xd7, 0x8a, 0x12, 0x4f, 0x7f, 0xfd, 0x14, 0xd9, 0x4f,
	0x8f, 0x2c, 0x0d, 0xd3, 0x30, 0x8f, 0xd9, 0xe2, 0xb2, 0x07, 0x38, 0xcd, 0x91, 0x4c, 0xd8, 0x5e,
	0x2b, 0x2b, 0xfe, 0xa5, 0xae, 0x65, 0xc8, 0xb6, 0x0c, 0x22, 0xec, 0x1a, 0xfd, 0xfa, 0x45, 0xc3,
	0x34, 0xc8, 0x8f, 0xcc, 0xeb, 0xdb, 0x17, 0xd8, 0x3d, 0xbb, 0xa1, 0xc1, 0x1b, 0x38, 0xdc, 0x0b,
	0xfe, 0xe1, 0xba, 0xbb, 0x4e, 0xb4, 0x12, 0xe3, 0x79, 0xe4, 0xb9, 0x96, 0xc9, 0x2a, 0xb1, 0x80,
	0x68, 0xfe, 0x1a, 0x71, 0x81, 0x9b, 0x71, 0xc9, 0xe5, 0x1f, 0xf6, 0xd3, 0xae, 0x89, 0x3c, 0x2d,
	0xbe, 0x03, 0x00, 0x00, 0xff, 0xff, 0x19, 0x05, 0xe0, 0xbc, 0x3e, 0x01, 0x00, 0x00,
}
