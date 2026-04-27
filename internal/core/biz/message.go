package biz

import (
	"context"
	"database/sql"
	"rain-im-server/internal/core/model"

	"github.com/uptrace/bun"
)

type MessageBiz struct {
	db *bun.DB
}

func NewMessageBiz(db *bun.DB) *MessageBiz {
	return &MessageBiz{db: db}
}

func (b *MessageBiz) CreateMessage(ctx context.Context, msg *model.Message) error {
	return b.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewInsert().Model(msg).Exec(ctx)
		return err
	})
}

func (b *MessageBiz) GetMessage(ctx context.Context, msg *model.Message) error {
	return b.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		err := tx.NewSelect().Model(msg).Scan(ctx)
		return err
	})
}

func (b *MessageBiz) UpdateMessage(ctx context.Context, msg *model.Message) error {
	return b.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewUpdate().Model(msg).Exec(ctx)
		return err
	})
}

func (b *MessageBiz) DeleteMessage(ctx context.Context, msg *model.Message) error {
	return b.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewDelete().Model(msg).Exec(ctx)
		return err
	})
}
