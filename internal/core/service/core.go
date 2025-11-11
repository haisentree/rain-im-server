package service

import (
	"net/http"
	"rain-im-server/protogo/core/v1/corev1connect"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
)

func NewServer() http.Server {
	base := &BaseServer{}
	mux := http.NewServeMux()

	path, handler := corev1connect.NewBaseServiceHandler(
		base,
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
