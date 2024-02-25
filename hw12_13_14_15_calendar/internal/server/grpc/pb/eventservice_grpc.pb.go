// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: eventservice/eventservice.proto

package pb

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	EventsService_Health_FullMethodName            = "/event.EventsService/Health"
	EventsService_CreateEvent_FullMethodName       = "/event.EventsService/CreateEvent"
	EventsService_ModifyEvent_FullMethodName       = "/event.EventsService/ModifyEvent"
	EventsService_RemoveEvent_FullMethodName       = "/event.EventsService/RemoveEvent"
	EventsService_GetEvent_FullMethodName          = "/event.EventsService/GetEvent"
	EventsService_GetEventsForDay_FullMethodName   = "/event.EventsService/GetEventsForDay"
	EventsService_GetEventsForWeek_FullMethodName  = "/event.EventsService/GetEventsForWeek"
	EventsService_GetEventsForMonth_FullMethodName = "/event.EventsService/GetEventsForMonth"
)

// EventsServiceClient is the client API for EventsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EventsServiceClient interface {
	Health(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*empty.Empty, error)
	CreateEvent(ctx context.Context, in *Event, opts ...grpc.CallOption) (*CreateEventResponse, error)
	ModifyEvent(ctx context.Context, in *Event, opts ...grpc.CallOption) (*empty.Empty, error)
	RemoveEvent(ctx context.Context, in *EventIdRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	GetEvent(ctx context.Context, in *EventIdRequest, opts ...grpc.CallOption) (*Event, error)
	GetEventsForDay(ctx context.Context, in *DateRequest, opts ...grpc.CallOption) (*Events, error)
	GetEventsForWeek(ctx context.Context, in *DateRequest, opts ...grpc.CallOption) (*Events, error)
	GetEventsForMonth(ctx context.Context, in *DateRequest, opts ...grpc.CallOption) (*Events, error)
}

type eventsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEventsServiceClient(cc grpc.ClientConnInterface) EventsServiceClient {
	return &eventsServiceClient{cc}
}

func (c *eventsServiceClient) Health(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, EventsService_Health_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventsServiceClient) CreateEvent(ctx context.Context, in *Event, opts ...grpc.CallOption) (*CreateEventResponse, error) {
	out := new(CreateEventResponse)
	err := c.cc.Invoke(ctx, EventsService_CreateEvent_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventsServiceClient) ModifyEvent(ctx context.Context, in *Event, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, EventsService_ModifyEvent_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventsServiceClient) RemoveEvent(ctx context.Context, in *EventIdRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, EventsService_RemoveEvent_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventsServiceClient) GetEvent(ctx context.Context, in *EventIdRequest, opts ...grpc.CallOption) (*Event, error) {
	out := new(Event)
	err := c.cc.Invoke(ctx, EventsService_GetEvent_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventsServiceClient) GetEventsForDay(ctx context.Context, in *DateRequest, opts ...grpc.CallOption) (*Events, error) {
	out := new(Events)
	err := c.cc.Invoke(ctx, EventsService_GetEventsForDay_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventsServiceClient) GetEventsForWeek(ctx context.Context, in *DateRequest, opts ...grpc.CallOption) (*Events, error) {
	out := new(Events)
	err := c.cc.Invoke(ctx, EventsService_GetEventsForWeek_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventsServiceClient) GetEventsForMonth(ctx context.Context, in *DateRequest, opts ...grpc.CallOption) (*Events, error) {
	out := new(Events)
	err := c.cc.Invoke(ctx, EventsService_GetEventsForMonth_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EventsServiceServer is the server API for EventsService service.
// All implementations must embed UnimplementedEventsServiceServer
// for forward compatibility
type EventsServiceServer interface {
	Health(context.Context, *empty.Empty) (*empty.Empty, error)
	CreateEvent(context.Context, *Event) (*CreateEventResponse, error)
	ModifyEvent(context.Context, *Event) (*empty.Empty, error)
	RemoveEvent(context.Context, *EventIdRequest) (*empty.Empty, error)
	GetEvent(context.Context, *EventIdRequest) (*Event, error)
	GetEventsForDay(context.Context, *DateRequest) (*Events, error)
	GetEventsForWeek(context.Context, *DateRequest) (*Events, error)
	GetEventsForMonth(context.Context, *DateRequest) (*Events, error)
	mustEmbedUnimplementedEventsServiceServer()
}

// UnimplementedEventsServiceServer must be embedded to have forward compatible implementations.
type UnimplementedEventsServiceServer struct {
}

func (UnimplementedEventsServiceServer) Health(context.Context, *empty.Empty) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Health not implemented")
}
func (UnimplementedEventsServiceServer) CreateEvent(context.Context, *Event) (*CreateEventResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateEvent not implemented")
}
func (UnimplementedEventsServiceServer) ModifyEvent(context.Context, *Event) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ModifyEvent not implemented")
}
func (UnimplementedEventsServiceServer) RemoveEvent(context.Context, *EventIdRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveEvent not implemented")
}
func (UnimplementedEventsServiceServer) GetEvent(context.Context, *EventIdRequest) (*Event, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEvent not implemented")
}
func (UnimplementedEventsServiceServer) GetEventsForDay(context.Context, *DateRequest) (*Events, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEventsForDay not implemented")
}
func (UnimplementedEventsServiceServer) GetEventsForWeek(context.Context, *DateRequest) (*Events, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEventsForWeek not implemented")
}
func (UnimplementedEventsServiceServer) GetEventsForMonth(context.Context, *DateRequest) (*Events, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEventsForMonth not implemented")
}
func (UnimplementedEventsServiceServer) mustEmbedUnimplementedEventsServiceServer() {}

// UnsafeEventsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EventsServiceServer will
// result in compilation errors.
type UnsafeEventsServiceServer interface {
	mustEmbedUnimplementedEventsServiceServer()
}

func RegisterEventsServiceServer(s grpc.ServiceRegistrar, srv EventsServiceServer) {
	s.RegisterService(&EventsService_ServiceDesc, srv)
}

func _EventsService_Health_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventsServiceServer).Health(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventsService_Health_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventsServiceServer).Health(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventsService_CreateEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Event)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventsServiceServer).CreateEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventsService_CreateEvent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventsServiceServer).CreateEvent(ctx, req.(*Event))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventsService_ModifyEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Event)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventsServiceServer).ModifyEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventsService_ModifyEvent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventsServiceServer).ModifyEvent(ctx, req.(*Event))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventsService_RemoveEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventsServiceServer).RemoveEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventsService_RemoveEvent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventsServiceServer).RemoveEvent(ctx, req.(*EventIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventsService_GetEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventsServiceServer).GetEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventsService_GetEvent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventsServiceServer).GetEvent(ctx, req.(*EventIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventsService_GetEventsForDay_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventsServiceServer).GetEventsForDay(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventsService_GetEventsForDay_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventsServiceServer).GetEventsForDay(ctx, req.(*DateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventsService_GetEventsForWeek_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventsServiceServer).GetEventsForWeek(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventsService_GetEventsForWeek_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventsServiceServer).GetEventsForWeek(ctx, req.(*DateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventsService_GetEventsForMonth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventsServiceServer).GetEventsForMonth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventsService_GetEventsForMonth_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventsServiceServer).GetEventsForMonth(ctx, req.(*DateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// EventsService_ServiceDesc is the grpc.ServiceDesc for EventsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EventsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "event.EventsService",
	HandlerType: (*EventsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Health",
			Handler:    _EventsService_Health_Handler,
		},
		{
			MethodName: "CreateEvent",
			Handler:    _EventsService_CreateEvent_Handler,
		},
		{
			MethodName: "ModifyEvent",
			Handler:    _EventsService_ModifyEvent_Handler,
		},
		{
			MethodName: "RemoveEvent",
			Handler:    _EventsService_RemoveEvent_Handler,
		},
		{
			MethodName: "GetEvent",
			Handler:    _EventsService_GetEvent_Handler,
		},
		{
			MethodName: "GetEventsForDay",
			Handler:    _EventsService_GetEventsForDay_Handler,
		},
		{
			MethodName: "GetEventsForWeek",
			Handler:    _EventsService_GetEventsForWeek_Handler,
		},
		{
			MethodName: "GetEventsForMonth",
			Handler:    _EventsService_GetEventsForMonth_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "eventservice/eventservice.proto",
}