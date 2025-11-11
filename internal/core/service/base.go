package service

import (
	"context"
	"database/sql"
	"rain-im-server/internal/core/global"
	"rain-im-server/internal/core/model"
	v1 "rain-im-server/protogo/core/v1"
	"rain-im-server/protogo/core/v1/corev1connect"

	"go.x2ox.com/sorbifolia/bunpgd/builder"

	"connectrpc.com/connect"
	"github.com/uptrace/bun"
)

type BaseServer struct {
	corev1connect.UnimplementedBaseServiceHandler
}

func (BaseServer) ListClient(ctx context.Context, req *connect.Request[v1.ListClientRequest]) (*connect.Response[v1.ListClientResponse], error) {
	var (
		arr        []*v1.Client
		count, err = builder.Select(global.DB.NewSelect().Model(&model.Client{}), req.Msg).ScanAndCount(ctx, &arr)
	)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnavailable, err)
	}

	return connect.NewResponse(&v1.ListClientResponse{Data: arr, Count: int64(count)}), nil
}

func (BaseServer) CreateClient(ctx context.Context, req *connect.Request[v1.CreateClientRequest]) (*connect.Response[v1.CreateClientResponse], error) {
	var client *model.Client

	if err := global.DB.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		client = &model.Client{Client: &v1.Client{}}

		if _, err := tx.NewInsert().Model(client).Exec(ctx); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, connect.NewError(connect.CodeUnimplemented, err)
	}

	return connect.NewResponse(&v1.CreateClientResponse{}), nil
}

func (BaseServer) ListMessage(ctx context.Context, req *connect.Request[v1.ListMessageRequest]) (*connect.Response[v1.ListMessageResponse], error) {
	var (
		arr        []*v1.Message
		count, err = builder.Select(global.DB.NewSelect().Model(&model.Message{}), req.Msg).ScanAndCount(ctx, &arr)
	)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnavailable, err)
	}

	return connect.NewResponse(&v1.ListMessageResponse{Data: arr, Count: int64(count)}), nil
}

func (BaseServer) CreateMessage(ctx context.Context, req *connect.Request[v1.CreateMessageRequest]) (*connect.Response[v1.CreateMessageResponse], error) {
	var message *model.Message

	if err := global.DB.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		message = &model.Message{Message: &v1.Message{
			SourceId: req.Msg.SourceId,
			TargetId: req.Msg.TargetId,
			Seq:      1,
			Content:  req.Msg.Content,
		}}

		if _, err := tx.NewInsert().Model(message).Exec(ctx); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, connect.NewError(connect.CodeUnimplemented, err)
	}

	return connect.NewResponse(&v1.CreateMessageResponse{}), nil
}
