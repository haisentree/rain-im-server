package model

import (
	corev1 "rain-im-server/protogo/core/v1"

	"github.com/uptrace/bun"
)

type Client struct {
	bun.BaseModel `bun:"table:client"`
	*corev1.Client
}

type Message struct {
	bun.BaseModel `bun:"table:message"`
	*corev1.Message
}
