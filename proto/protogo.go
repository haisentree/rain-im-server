package protogo

import (
	_ "connectrpc.com/connect"
	_ "github.com/srikrsna/protoc-gen-gotag/tagger"
	_ "google.golang.org/protobuf/proto"
)

//go:generate buf generate ../proto --template buf.gen.yaml
//go:generate buf generate ../protogo --template buf.gen.tag.yaml
