package service

import (
	"context"
	"errors"
	v1 "rain-im-server/protogo/core/v1"
	"rain-im-server/protogo/core/v1/corev1connect"

	"connectrpc.com/connect"
)

type BaseServer struct {
	corev1connect.UnimplementedBaseServiceHandler
}

func (BaseServer) ListClient(context.Context, *connect.Request[v1.ListClientRequest]) (*connect.Response[v1.ListClientResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("core.v1.BaseService.ListClient is not implemented"))
}

func (BaseServer) CreateClient(context.Context, *connect.Request[v1.CreateClientRequest]) (*connect.Response[v1.CreateClientResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("core.v1.BaseService.CreateClient is not implemented"))
}

func (BaseServer) ListMessage(context.Context, *connect.Request[v1.ListMessageRequest]) (*connect.Response[v1.ListMessageResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("core.v1.BaseService.ListMessage is not implemented"))
}

func (BaseServer) CreateMessage(context.Context, *connect.Request[v1.CreateMessageRequest]) (*connect.Response[v1.CreateMessageResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("core.v1.BaseService.CreateMessage is not implemented"))
}
