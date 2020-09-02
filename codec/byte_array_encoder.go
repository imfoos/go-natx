package codec

import (
	"github.com/go-netty/go-netty"
)

//
//import (
//	"bytes"
//	"github.com/go-netty/go-netty"
//	"github.com/go-netty/go-netty/codec"
//	"probe/tunnel/handler"
//)
//
type ByteArrayEncoder struct {
	//common.NatxCommonHandler
}

func NewByteArrayEncoder() netty.OutboundHandler {
	return &ByteArrayEncoder{}
}

func (p *ByteArrayEncoder) HandleWrite(ctx netty.OutboundContext, message netty.Message) {
	//out := new(bytes.Buffer)
	//out.ReadFrom(message)
	ctx.HandleWrite(message)
}
