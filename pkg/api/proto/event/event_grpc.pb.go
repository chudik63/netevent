// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: proto/event/event.proto

package event

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	EventService_CreateEvent_FullMethodName           = "/event.EventService/CreateEvent"
	EventService_ReadEvent_FullMethodName             = "/event.EventService/ReadEvent"
	EventService_UpdateEvent_FullMethodName           = "/event.EventService/UpdateEvent"
	EventService_DeleteEvent_FullMethodName           = "/event.EventService/DeleteEvent"
	EventService_ListEvents_FullMethodName            = "/event.EventService/ListEvents"
	EventService_ListEventsByCreator_FullMethodName   = "/event.EventService/ListEventsByCreator"
	EventService_ListEventsByInterests_FullMethodName = "/event.EventService/ListEventsByInterests"
	EventService_RegisterUser_FullMethodName          = "/event.EventService/RegisterUser"
	EventService_SetChatStatus_FullMethodName         = "/event.EventService/SetChatStatus"
	EventService_ListUsersToChat_FullMethodName       = "/event.EventService/ListUsersToChat"
	EventService_ListEventsByUser_FullMethodName      = "/event.EventService/ListEventsByUser"
)

// EventServiceClient is the client API for EventService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EventServiceClient interface {
	CreateEvent(ctx context.Context, in *CreateEventRequest, opts ...grpc.CallOption) (*CreateEventResponse, error)
	ReadEvent(ctx context.Context, in *ReadEventRequest, opts ...grpc.CallOption) (*ReadEventResponse, error)
	UpdateEvent(ctx context.Context, in *UpdateEventRequest, opts ...grpc.CallOption) (*UpdateEventResponse, error)
	DeleteEvent(ctx context.Context, in *DeleteEventRequest, opts ...grpc.CallOption) (*DeleteEventResponse, error)
	ListEvents(ctx context.Context, in *ListEventsRequest, opts ...grpc.CallOption) (*ListEventsResponse, error)
	ListEventsByCreator(ctx context.Context, in *ListEventsByCreatorRequest, opts ...grpc.CallOption) (*ListEventsByCreatorResponse, error)
	ListEventsByInterests(ctx context.Context, in *ListEventsByInterestsRequest, opts ...grpc.CallOption) (*ListEventsByInterestsResponse, error)
	RegisterUser(ctx context.Context, in *RegisterUserRequest, opts ...grpc.CallOption) (*RegisterUserResponse, error)
	SetChatStatus(ctx context.Context, in *SetChatStatusRequest, opts ...grpc.CallOption) (*SetChatStatusResponse, error)
	ListUsersToChat(ctx context.Context, in *ListUsersToChatRequest, opts ...grpc.CallOption) (*ListUsersToChatResponse, error)
	ListEventsByUser(ctx context.Context, in *ListEventsByUserRequest, opts ...grpc.CallOption) (*ListEventsByUserResponse, error)
}

type eventServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEventServiceClient(cc grpc.ClientConnInterface) EventServiceClient {
	return &eventServiceClient{cc}
}

func (c *eventServiceClient) CreateEvent(ctx context.Context, in *CreateEventRequest, opts ...grpc.CallOption) (*CreateEventResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateEventResponse)
	err := c.cc.Invoke(ctx, EventService_CreateEvent_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) ReadEvent(ctx context.Context, in *ReadEventRequest, opts ...grpc.CallOption) (*ReadEventResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReadEventResponse)
	err := c.cc.Invoke(ctx, EventService_ReadEvent_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) UpdateEvent(ctx context.Context, in *UpdateEventRequest, opts ...grpc.CallOption) (*UpdateEventResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateEventResponse)
	err := c.cc.Invoke(ctx, EventService_UpdateEvent_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) DeleteEvent(ctx context.Context, in *DeleteEventRequest, opts ...grpc.CallOption) (*DeleteEventResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteEventResponse)
	err := c.cc.Invoke(ctx, EventService_DeleteEvent_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) ListEvents(ctx context.Context, in *ListEventsRequest, opts ...grpc.CallOption) (*ListEventsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListEventsResponse)
	err := c.cc.Invoke(ctx, EventService_ListEvents_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) ListEventsByCreator(ctx context.Context, in *ListEventsByCreatorRequest, opts ...grpc.CallOption) (*ListEventsByCreatorResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListEventsByCreatorResponse)
	err := c.cc.Invoke(ctx, EventService_ListEventsByCreator_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) ListEventsByInterests(ctx context.Context, in *ListEventsByInterestsRequest, opts ...grpc.CallOption) (*ListEventsByInterestsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListEventsByInterestsResponse)
	err := c.cc.Invoke(ctx, EventService_ListEventsByInterests_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) RegisterUser(ctx context.Context, in *RegisterUserRequest, opts ...grpc.CallOption) (*RegisterUserResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RegisterUserResponse)
	err := c.cc.Invoke(ctx, EventService_RegisterUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) SetChatStatus(ctx context.Context, in *SetChatStatusRequest, opts ...grpc.CallOption) (*SetChatStatusResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SetChatStatusResponse)
	err := c.cc.Invoke(ctx, EventService_SetChatStatus_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) ListUsersToChat(ctx context.Context, in *ListUsersToChatRequest, opts ...grpc.CallOption) (*ListUsersToChatResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListUsersToChatResponse)
	err := c.cc.Invoke(ctx, EventService_ListUsersToChat_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) ListEventsByUser(ctx context.Context, in *ListEventsByUserRequest, opts ...grpc.CallOption) (*ListEventsByUserResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListEventsByUserResponse)
	err := c.cc.Invoke(ctx, EventService_ListEventsByUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EventServiceServer is the server API for EventService service.
// All implementations must embed UnimplementedEventServiceServer
// for forward compatibility.
type EventServiceServer interface {
	CreateEvent(context.Context, *CreateEventRequest) (*CreateEventResponse, error)
	ReadEvent(context.Context, *ReadEventRequest) (*ReadEventResponse, error)
	UpdateEvent(context.Context, *UpdateEventRequest) (*UpdateEventResponse, error)
	DeleteEvent(context.Context, *DeleteEventRequest) (*DeleteEventResponse, error)
	ListEvents(context.Context, *ListEventsRequest) (*ListEventsResponse, error)
	ListEventsByCreator(context.Context, *ListEventsByCreatorRequest) (*ListEventsByCreatorResponse, error)
	ListEventsByInterests(context.Context, *ListEventsByInterestsRequest) (*ListEventsByInterestsResponse, error)
	RegisterUser(context.Context, *RegisterUserRequest) (*RegisterUserResponse, error)
	SetChatStatus(context.Context, *SetChatStatusRequest) (*SetChatStatusResponse, error)
	ListUsersToChat(context.Context, *ListUsersToChatRequest) (*ListUsersToChatResponse, error)
	ListEventsByUser(context.Context, *ListEventsByUserRequest) (*ListEventsByUserResponse, error)
	mustEmbedUnimplementedEventServiceServer()
}

// UnimplementedEventServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedEventServiceServer struct{}

func (UnimplementedEventServiceServer) CreateEvent(context.Context, *CreateEventRequest) (*CreateEventResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateEvent not implemented")
}
func (UnimplementedEventServiceServer) ReadEvent(context.Context, *ReadEventRequest) (*ReadEventResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadEvent not implemented")
}
func (UnimplementedEventServiceServer) UpdateEvent(context.Context, *UpdateEventRequest) (*UpdateEventResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateEvent not implemented")
}
func (UnimplementedEventServiceServer) DeleteEvent(context.Context, *DeleteEventRequest) (*DeleteEventResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteEvent not implemented")
}
func (UnimplementedEventServiceServer) ListEvents(context.Context, *ListEventsRequest) (*ListEventsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEvents not implemented")
}
func (UnimplementedEventServiceServer) ListEventsByCreator(context.Context, *ListEventsByCreatorRequest) (*ListEventsByCreatorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEventsByCreator not implemented")
}
func (UnimplementedEventServiceServer) ListEventsByInterests(context.Context, *ListEventsByInterestsRequest) (*ListEventsByInterestsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEventsByInterests not implemented")
}
func (UnimplementedEventServiceServer) RegisterUser(context.Context, *RegisterUserRequest) (*RegisterUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterUser not implemented")
}
func (UnimplementedEventServiceServer) SetChatStatus(context.Context, *SetChatStatusRequest) (*SetChatStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetChatStatus not implemented")
}
func (UnimplementedEventServiceServer) ListUsersToChat(context.Context, *ListUsersToChatRequest) (*ListUsersToChatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListUsersToChat not implemented")
}
func (UnimplementedEventServiceServer) ListEventsByUser(context.Context, *ListEventsByUserRequest) (*ListEventsByUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEventsByUser not implemented")
}
func (UnimplementedEventServiceServer) mustEmbedUnimplementedEventServiceServer() {}
func (UnimplementedEventServiceServer) testEmbeddedByValue()                      {}

// UnsafeEventServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EventServiceServer will
// result in compilation errors.
type UnsafeEventServiceServer interface {
	mustEmbedUnimplementedEventServiceServer()
}

func RegisterEventServiceServer(s grpc.ServiceRegistrar, srv EventServiceServer) {
	// If the following call pancis, it indicates UnimplementedEventServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&EventService_ServiceDesc, srv)
}

func _EventService_CreateEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).CreateEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventService_CreateEvent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).CreateEvent(ctx, req.(*CreateEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_ReadEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).ReadEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventService_ReadEvent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).ReadEvent(ctx, req.(*ReadEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_UpdateEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).UpdateEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventService_UpdateEvent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).UpdateEvent(ctx, req.(*UpdateEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_DeleteEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).DeleteEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventService_DeleteEvent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).DeleteEvent(ctx, req.(*DeleteEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_ListEvents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListEventsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).ListEvents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventService_ListEvents_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).ListEvents(ctx, req.(*ListEventsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_ListEventsByCreator_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListEventsByCreatorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).ListEventsByCreator(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventService_ListEventsByCreator_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).ListEventsByCreator(ctx, req.(*ListEventsByCreatorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_ListEventsByInterests_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListEventsByInterestsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).ListEventsByInterests(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventService_ListEventsByInterests_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).ListEventsByInterests(ctx, req.(*ListEventsByInterestsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_RegisterUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).RegisterUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventService_RegisterUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).RegisterUser(ctx, req.(*RegisterUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_SetChatStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetChatStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).SetChatStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventService_SetChatStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).SetChatStatus(ctx, req.(*SetChatStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_ListUsersToChat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListUsersToChatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).ListUsersToChat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventService_ListUsersToChat_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).ListUsersToChat(ctx, req.(*ListUsersToChatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_ListEventsByUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListEventsByUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).ListEventsByUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventService_ListEventsByUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).ListEventsByUser(ctx, req.(*ListEventsByUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// EventService_ServiceDesc is the grpc.ServiceDesc for EventService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EventService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "event.EventService",
	HandlerType: (*EventServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateEvent",
			Handler:    _EventService_CreateEvent_Handler,
		},
		{
			MethodName: "ReadEvent",
			Handler:    _EventService_ReadEvent_Handler,
		},
		{
			MethodName: "UpdateEvent",
			Handler:    _EventService_UpdateEvent_Handler,
		},
		{
			MethodName: "DeleteEvent",
			Handler:    _EventService_DeleteEvent_Handler,
		},
		{
			MethodName: "ListEvents",
			Handler:    _EventService_ListEvents_Handler,
		},
		{
			MethodName: "ListEventsByCreator",
			Handler:    _EventService_ListEventsByCreator_Handler,
		},
		{
			MethodName: "ListEventsByInterests",
			Handler:    _EventService_ListEventsByInterests_Handler,
		},
		{
			MethodName: "RegisterUser",
			Handler:    _EventService_RegisterUser_Handler,
		},
		{
			MethodName: "SetChatStatus",
			Handler:    _EventService_SetChatStatus_Handler,
		},
		{
			MethodName: "ListUsersToChat",
			Handler:    _EventService_ListUsersToChat_Handler,
		},
		{
			MethodName: "ListEventsByUser",
			Handler:    _EventService_ListEventsByUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/event/event.proto",
}
