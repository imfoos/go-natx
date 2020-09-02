package codec

import (
	"github.com/go-netty/go-netty"
)

type LogDecoder struct {
}

func NewLogDecoder() netty.CodecHandler {
	return &LogDecoder{}
}

func (l LogDecoder) CodecName() string {
	return "LogDecoder"
}

func (l LogDecoder) HandleRead(ctx netty.InboundContext, message netty.Message) {
	//fmt.Println(reflect.TypeOf(message))
	//bytes := utils.MustToBytes(message)
	ctx.HandleRead(message)
}

func (l LogDecoder) HandleWrite(ctx netty.OutboundContext, message netty.Message) {
	//fmt.Printf("Logger HandleWrite : %v \n", message)
	ctx.HandleWrite(message)
}
