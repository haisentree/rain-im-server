package model

import (
	"context"
	"rain-im-server/internal/core/global"

	"log/slog"
)

func init() {
	ctx := context.Background()
	model := []any{&Client{}, &Message{}}
	global.DB.RegisterModel(model...)

	for _, val := range model {
		_, err := global.DB.NewCreateTable().Model(val).IfNotExists().Exec(ctx)
		if err != nil {
			slog.Error("err",
				slog.String("err", err.Error()),
				slog.Any("model", val),
			)
		}
	}
}
