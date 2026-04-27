package service

import (
	"net/http"
	"rain-im-server/global"
	"rain-im-server/internal/core/biz"
	"rain-im-server/protogo/core/v1/corev1connect"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
)

func NewServer() http.Server {
	messageBiz := biz.NewMessageBiz(global.DB)

	baseServer := NewBaseServer(messageBiz)
	mux := http.NewServeMux()

	path, handler := corev1connect.NewBaseServiceHandler(
		baseServer,
		connect.WithInterceptors(validate.NewInterceptor()),
	)

	mux.Handle(path, handler)
	p := new(http.Protocols)
	p.SetHTTP1(true)
	p.SetUnencryptedHTTP2(true)

	return http.Server{
		Addr:      "localhost:8080",
		Handler:   mux,
		Protocols: p,
	}
}
