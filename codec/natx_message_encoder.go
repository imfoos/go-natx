package codec

import (
	"bytes"
	"github.com/imfoos/go-natx/common"
	"github.com/imfoos/go-natx/protocol"

	"github.com/go-netty/go-netty"
)

type NatxMessageEncoder struct {
	common.NatxCommonHandler
	netty.OutboundHandler
}

func NewNatxMessageEncoder() *NatxMessageEncoder {
	return &NatxMessageEncoder{}
}

func (p *NatxMessageEncoder) CodecName() string {
	return "NatxMessageEncoder"
}

func (p *NatxMessageEncoder) HandleWrite(ctx netty.OutboundContext, message netty.Message) {
	out := new(bytes.Buffer)
	switch message.(type) {
	case *protocol.NatxMessage:
		p.Encode(ctx, message.(*protocol.NatxMessage), out)
		//if out.Len() > 150 {
		//fmt.Printf("out: %v \n", out)
		//}
		ctx.HandleWrite(out)
	case []byte:
		ctx.HandleWrite(message)
	}
	//ctx.Write(out)
}

//func (p *NatxMessageEncoder) HandleRead(ctx netty.InboundContext, message netty.Message) {
//
//	fmt.Println("natxMessageEncoder handle read")
//	ctx.HandleRead(message)
//}
