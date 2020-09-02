package handler

import (
	"github.com/go-netty/go-netty"
	"github.com/imfoos/go-natx/common"
	"github.com/imfoos/go-natx/protocol"
)

type LocalProxyHandler struct {
	common.NatxCommonHandler
	proxyHandler    *NatxClientHandler
	remoteChannelId string
}

func NewLocalProxyHandler(proxyHandler *NatxClientHandler, remoteChannelId string) *LocalProxyHandler {
	return &LocalProxyHandler{
		proxyHandler:    proxyHandler,
		remoteChannelId: remoteChannelId,
	}
}
func (p *LocalProxyHandler) HandleRead(ctx netty.InboundContext, message netty.Message) {
	metaData := make(map[string]interface{})
	metaData["channelId"] = p.remoteChannelId
	natxMessage := &protocol.NatxMessage{
		Type:     protocol.DATA,
		Data:     message.([]byte),
		MetaData: metaData,
	}

	p.proxyHandler.Ctx.Write(natxMessage)
}

func (p *LocalProxyHandler) HandleActive(ctx netty.ActiveContext) {
	//log.Infoln("LocalProxyHandler HandleActive")
	p.NatxCommonHandler.HandleActive(ctx)
}
func (p *LocalProxyHandler) HandleEvent(ctx netty.EventContext, event netty.Event) {
	//log.Infoln("LocalProxyHandler HandleEvent")
	ctx.HandleEvent(event)
}

func (p *LocalProxyHandler) HandleInactive(ctx netty.InactiveContext, ex netty.Exception) {
	//log.Infoln("LocalProxyHandler HandleInactive")
	metaData := make(map[string]interface{})
	metaData["channelId"] = p.remoteChannelId
	message := &protocol.NatxMessage{
		Type:     protocol.DISCONNECTED,
		MetaData: metaData,
	}
	p.proxyHandler.Ctx.Write(message)
}
