package common

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/go-netty/go-netty"
	"github.com/imfoos/go-natx/protocol"
	"github.com/prometheus/common/log"
	"io"
)

type NatxCommonHandler struct {
	//netty.CodecHandler
	Ctx netty.ActiveContext
}

func (p *NatxCommonHandler) HandleActive(ctx netty.ActiveContext) {
	p.Ctx = ctx
	ctx.HandleActive()
}
func (p *NatxCommonHandler) HandleWrite(ctx netty.OutboundContext, message netty.Message) {
	ctx.HandleWrite(message)
}

//func (p *NatxCommonHandler) HandleRead(ctx netty.InboundContext, message netty.Message) {
//	ctx.HandleRead(message)
//}
//
//func (p *NatxCommonHandler) HandleWrite(ctx netty.OutboundContext, message netty.Message) {
//	//ctx.HandleWrite(message)
//	ctx.HandleWrite(message)
//}

func (p *NatxCommonHandler) HandleException(ctx netty.ExceptionContext, ex netty.Exception) {
	fmt.Println(ctx)
	fmt.Println(string(ex.Stack()))
}
func (p *NatxCommonHandler) HandleInactive(ctx netty.InactiveContext, ex netty.Exception) {
	log.Errorln(ctx)
	log.Errorln(ex.Error())
}

func (p *NatxCommonHandler) Encode(ctx netty.HandlerContext, message netty.Message, out *bytes.Buffer) {
	nMesaage := message.(*protocol.NatxMessage)

	binary.Write(out, binary.BigEndian, int32(nMesaage.Type))
	metaBytes, _ := json.Marshal(nMesaage.MetaData)
	binary.Write(out, binary.BigEndian, int32(len(metaBytes)))
	binary.Write(out, binary.BigEndian, metaBytes)
	binary.Write(out, binary.BigEndian, nMesaage.Data)

	//log.Infof("ChannelID: %v ,type: %v, encode: %v\n", ctx.Channel().ID(), nMesaage.Type.String(), out.Bytes())
	//log.Infof("ChannelID: %v ,meta: %v\n", ctx.Channel().ID(), string(metaBytes))
	//log.Infof("ChannelID: %v ,data: %v\n", ctx.Channel().ID(), string(nMesaage.Data))
}

func (p *NatxCommonHandler) Decode(ctx netty.HandlerContext, message netty.Message) *protocol.NatxMessage {

	reader := new(bytes.Buffer)
	reader.ReadFrom(message.(io.Reader))

	var tp int32
	binary.Read(reader, binary.BigEndian, &tp)
	//log.Infof("ChannelID: %v ,type: %v, decode: %v \n", ctx.Channel().ID(), protocol.NatxMessageType(tp).String(), tp)
	var metaDateLength int32
	binary.Read(reader, binary.BigEndian, &metaDateLength)

	metaDataBytes := make([]byte, metaDateLength)
	reader.Read(metaDataBytes)
	var metaData map[string]interface{}
	err := json.Unmarshal(metaDataBytes, &metaData)
	if err != nil {
		log.Errorln(err)
	}
	//log.Infof("ChannelID: %v ,meta: %v \n", ctx.Channel().ID(), string(metaDataBytes))
	data := make([]byte, reader.Len())
	reader.Read(data)
	//log.Infof("ChannelID: %v ,data: %v \n", ctx.Channel().ID(), string(data))
	return &protocol.NatxMessage{
		Type:     protocol.NatxMessageType(tp),
		MetaData: metaData,
		Data:     data,
	}
}
