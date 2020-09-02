package codec

import (
	"github.com/go-netty/go-netty"
	"github.com/imfoos/go-natx/common"
)

type NatxMessageDecoder struct {
	common.NatxCommonHandler
}

func NewNatxMessageDecoder() netty.ChannelInboundHandler {
	return &NatxMessageDecoder{}
}

func (p *NatxMessageDecoder) CodecName() string {
	return "NatxMessageDecoder"
}

func (p *NatxMessageDecoder) HandleRead(ctx netty.InboundContext, message netty.Message) {
	ctx.HandleRead(p.Decode(ctx, message))
}
