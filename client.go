package main

import (
	"encoding/binary"
	"fmt"
	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty/codec/frame"
	"github.com/go-netty/go-netty/transport/tcp"
	"github.com/imfoos/go-natx/codec"
	"github.com/imfoos/go-natx/handler"
	"log"
	"math"
	"time"
)

func Start(serverAddress string, serverPort int, proxyAddress string, proxyPort int, port int, password string) {

	for {
		var channelInitializer = func(channel netty.Channel) {
			natxClientHandler := handler.NewNatxClientHandler(port, password, proxyAddress, proxyPort)
			pipeline := channel.Pipeline()
			pipeline.AddLast(codec.NewLogDecoder()).
				AddLast(frame.LengthFieldCodec(binary.BigEndian, math.MaxInt64, 0, 4, 0, 4)).
				AddLast(codec.NewNatxMessageDecoder()).
				AddLast(codec.NewNatxMessageEncoder()).
				AddLast(netty.ReadIdleHandler(60 * time.Second)).
				AddLast(netty.WriteIdleHandler(30 * time.Second)).
				AddLast(natxClientHandler)
		}

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
		channel, err := bootstrap.Connect(fmt.Sprintf("tcp://%v:%v", serverAddress, serverPort), "",
			tcp.WithOptions(tcpOptions))
		if nil != err {
			log.Println(err)
		}
		<-channel.Context().Done()
		log.Println("tunnel close...")
	}
}

func main() {

	//ServerAddress = "192.168.1.186"
	serverAddress := "127.0.0.1"
	serverPort := 7731
	proxyAddress := "127.0.0.1"
	proxyPort := 11111
	//proxyAddress = "172.16.10.201"
	//proxyPort = 3306
	Start(serverAddress, serverPort, proxyAddress, proxyPort, 9191, "1122")
}
