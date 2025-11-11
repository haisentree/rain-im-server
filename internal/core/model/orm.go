package model

import (
	"context"
	"errors"

	basev1 "rain-im-server/protogo/base/v1"
	v1 "rain-im-server/protogo/core/v1"

	"github.com/uptrace/bun"
)

func (x *Client) BeforeAppendModel(_ context.Context, query bun.Query) error {
	if x.Client == nil {
		x.Client = &v1.Client{}
	}

	t := basev1.Now()

	switch query.(type) {
	case *bun.InsertQuery:
		if x.Id == nil {
			x.Id = basev1.NewUUID()
		}

		if x.CreatedAt == nil {
			x.CreatedAt = t
		}

		if x.UpdatedAt == nil {
			x.UpdatedAt = t
		}

	case *bun.UpdateQuery:
		x.UpdatedAt = t
	}

	return nil
}

func (x *Client) BeforeUpdate(_ context.Context, query *bun.UpdateQuery) error {
	query.Set("updated_at = ?", basev1.Now())
	return nil
}

func (x *Message) BeforeAppendModel(_ context.Context, query bun.Query) error {
	if x.Message == nil {
		x.Message = &v1.Message{}
	}

	t := basev1.Now()

	switch query.(type) {
	case *bun.InsertQuery:
		if x.Id == nil {
			x.Id = basev1.NewUUID()
		}

		if x.CreatedAt == nil {
			x.CreatedAt = t
		}

		if x.UpdatedAt == nil {
			x.UpdatedAt = t
		}

		// SourceId和TargetId不能为空
		if x.SourceId == nil {
			return errors.New("SourceId must not empty")
		}

		if x.TargetId == nil {
			return errors.New("TargetId must not empty")
		}
	case *bun.UpdateQuery:
		x.UpdatedAt = t
	}

	return nil
}

func (x *Message) BeforeUpdate(_ context.Context, query *bun.UpdateQuery) error {
	query.Set("updated_at = ?", basev1.Now())
	return nil
}
