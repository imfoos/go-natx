package net

import (
	"fmt"
	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty/transport/tcp"
	"log"
	"time"
)

type TcpConnection struct {
}

func NewTcpConnection() *TcpConnection {
	return &TcpConnection{}
}

func (p *TcpConnection) Connect(host string, port int, channelInitializer netty.ChannelInitializer) {

	tcpOptions := &tcp.Options{
		Timeout:         time.Second * 3,
		KeepAlive:       true,
		KeepAlivePeriod: time.Second * 5,
		Linger:          0,
		NoDelay:         true,
		SockBuf:         1024,
	}

	bootstrap := netty.NewBootstrap()
	bootstrap.ClientInitializer(channelInitializer)
	bootstrap.Transport(tcp.New())
	bootstrap.Channel(netty.NewBufferedChannel(256, 128))
	channel, err := bootstrap.Connect(fmt.Sprintf("tcp://%v:%v", host, port), nil, tcp.WithOptions(tcpOptions))
	if err != nil {
		log.Println(err)
	} else {
		select {
		case <-channel.Context().Done():
		case <-bootstrap.Context().Done():
		}
		log.Println("TcpConnection channel.Context().Done()  ")
	}
}
