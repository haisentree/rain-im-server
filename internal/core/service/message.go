package service

import (
	"rain-im-server/internal/core/biz"
	"rain-im-server/protogo/core/v1/corev1connect"
)

type MessageServer struct {
	corev1connect.UnimplementedMessageServiceHandler
	messageBiz *biz.MessageBiz
}
