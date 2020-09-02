package codec

import (
	"bufio"
	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty/utils"
	"github.com/prometheus/common/log"
)

type ByteArrayDecoder struct {
}

func NewByteArrayDecoder() netty.InboundHandler {
	return &ByteArrayDecoder{}
}

func (p *ByteArrayDecoder) HandleRead(ctx netty.InboundContext, message netty.Message) {
	reader := bufio.NewReader(utils.MustToReader(message))
	_, err := reader.Peek(1)
	if err != nil {
		log.Errorln(err)
		ctx.Close(err)
	}
	if buffered := reader.Buffered(); buffered > 0 {
		data := make([]byte, buffered)
		reader.Read(data)
		ctx.HandleRead(data)
	}
}
