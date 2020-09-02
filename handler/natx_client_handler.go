package handler

import (
	"github.com/go-netty/go-netty"
	"github.com/imfoos/go-natx/codec"
	"github.com/imfoos/go-natx/common"
	"github.com/imfoos/go-natx/net"
	"github.com/imfoos/go-natx/protocol"
	"log"
	"sync"
)

type NatxClientHandler struct {
	netty.ActiveHandler
	common.NatxCommonHandler
	port              int
	password          string
	proxyAddress      string
	proxyPort         int
	channelHandlerMap *ProxyHandlerMap
	channelExector    netty.ChannelExecutorFactory
}

type ProxyHandlerMap struct {
	sync.Mutex
	m map[string]*LocalProxyHandler
}

func initProxyHandlerMap() *ProxyHandlerMap {
	return &ProxyHandlerMap{
		m: make(map[string]*LocalProxyHandler),
	}
}

func (p *ProxyHandlerMap) Set(channelId string, proxyHandler *LocalProxyHandler) {
	p.Lock()
	defer p.Unlock()
	p.m[channelId] = proxyHandler
}
func (p *ProxyHandlerMap) Get(channelId string) *LocalProxyHandler {
	p.Lock()
	defer p.Unlock()
	return p.m[channelId]
}
func (p *ProxyHandlerMap) Del(channelId string) {
	p.Lock()
	defer p.Unlock()
	delete(p.m, channelId)
}

func NewNatxClientHandler(port int, password, proxyAddress string, proxyPort int) netty.InboundHandler {
	return &NatxClientHandler{
		port:              port,
		password:          password,
		proxyAddress:      proxyAddress,
		proxyPort:         proxyPort,
		channelHandlerMap: initProxyHandlerMap(),
		channelExector:    netty.NewFixedChannelExecutor(100, 100),
	}
}

func (p *NatxClientHandler) HandleActive(ctx netty.ActiveContext) {

	metaData := make(map[string]interface{})
	metaData["port"] = p.port
	metaData["password"] = p.password

	natxMessage := &protocol.NatxMessage{
		Type:     protocol.REGISTER,
		MetaData: metaData,
	}
	ctx.Write(natxMessage)
	p.NatxCommonHandler.HandleActive(ctx)
}

func (p *NatxClientHandler) HandleEvent(ctx netty.EventContext, event netty.Event) {

	switch event.(type) {
	case netty.WriteIdleEvent:
		//log.Println("HandleEvent WriteIdleEvent  .......")
		natxMessage := &protocol.NatxMessage{
			Type:     protocol.KEEPALIVE,
			MetaData: make(map[string]interface{}),
		}
		ctx.Write(natxMessage)
	case netty.ReadIdleEvent:
		//log.Println("HandleEvent ReadIdleEvent  .......")
		ctx.Close(nil)
	}
}

func (p *NatxClientHandler) HandleRead(ctx netty.InboundContext, message netty.Message) {
	natxMessage := message.(*protocol.NatxMessage)
	//log.Printf("NatxClientHandler %v  \n", protocol.NatxMessageType(natxMessage.Type))
	if natxMessage.Type == protocol.REGISTER_RESULT {
		p.processRegisterResult(natxMessage)
	} else if natxMessage.Type == protocol.CONNECTED {
		p.processConnected(natxMessage)
	} else if natxMessage.Type == protocol.DISCONNECTED {
		p.processDisconnected(natxMessage)
	} else if natxMessage.Type == protocol.DATA {
		p.processData(natxMessage)
	} else if natxMessage.Type == protocol.KEEPALIVE {
		// 心跳包, 不处理
	} else {
	}
}
func (p *NatxClientHandler) HandleWrite(ctx netty.OutboundContext, message netty.Message) {
	ctx.HandleWrite(message)
}

func (p *NatxClientHandler) processRegisterResult(natxMessage *protocol.NatxMessage) {
	if natxMessage.MetaData["success"].(bool) {
		log.Println("Register to Natx server")
	} else {
		log.Printf("Register fail: %v \n", natxMessage.MetaData["reason"])
		p.Ctx.Close(nil)
	}
}

func (p *NatxClientHandler) processConnected(natxMessage *protocol.NatxMessage) {
	ch := make(chan struct{})
	handler := p
	go func() {
		channelId := natxMessage.MetaData["channelId"].(string)
		localConnection := net.NewTcpConnection()
		localConnection.Connect(p.proxyAddress, p.proxyPort, func(channel netty.Channel) {
			localProxyHandler := NewLocalProxyHandler(p, natxMessage.MetaData["channelId"].(string))
			channel.Pipeline().
				//AddLast(frame.LengthFieldCodec(binary.BigEndian, 1048576, 0, 4, 0, 4)).
				AddLast(codec.NewByteArrayDecoder()).
				AddLast(codec.NewByteArrayEncoder()).
				AddLast(localProxyHandler)
			handler.channelHandlerMap.Set(channelId, localProxyHandler)
			log.Printf("LocalProxyHandler Connect , %v \n", channelId)
		}, ch)
		natxMessage := &protocol.NatxMessage{
			Type:     protocol.DISCONNECTED,
			MetaData: map[string]interface{}{"channelId": channelId},
		}
		handler.Ctx.Write(natxMessage)
	}()
	<-ch
	//log.Println("NatxClientHandler processConnected...")
}

func (p *NatxClientHandler) processDisconnected(natxMessage *protocol.NatxMessage) {

	channelId := natxMessage.MetaData["channelId"].(string)
	defer p.channelHandlerMap.Del(channelId)
	handler := p.channelHandlerMap.Get(channelId)
	log.Printf("LocalProxyHandler disconnected   %v \n", channelId)
	if handler != nil {
		handler.Ctx.Close(nil)
	}
}

func (p *NatxClientHandler) processData(natxMessage *protocol.NatxMessage) {

	channelId := natxMessage.MetaData["channelId"].(string)
	handler := p.channelHandlerMap.Get(channelId)
	if nil != handler {
		//log.Printf("processData...... %v \n", natxMessage.Data)
		handler.Ctx.Write(natxMessage.Data)
	}
}
