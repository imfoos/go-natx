package codec

import (
	"github.com/go-netty/go-netty"
)

type ByteArrayEncoder struct {
}

func NewByteArrayEncoder() netty.OutboundHandler {
	return &ByteArrayEncoder{}
}

func (p *ByteArrayEncoder) HandleWrite(ctx netty.OutboundContext, message netty.Message) {
	ctx.HandleWrite(message)
}
