package codec

import (
	"bufio"
	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty/utils"
	"github.com/prometheus/common/log"
)

type ByteArrayDecoder struct {
	//common.NatxCommonHandler
}

func NewByteArrayDecoder() netty.InboundHandler {
	return &ByteArrayDecoder{}
}

//
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

	//	bbs, e := bbf.ReadBytes(byte('\n'))
	//	println(string(bbs))
	//	println(e)
	//	//if aa := bbf.Buffered(); aa > 0 {
	//	//	bs := make([]byte, aa)
	//	//	bbf.Read(bs)
	//	//	println(string(bs))
	//	//}
	//}
	//data, _ := ioutil.ReadAll(bbf)
	//scanner := bufio.NewScanner(reader)
	//split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
	//	fmt.Printf("%t\t%d\t%s\n", atEOF, len(data), data)
	//	if atEOF {
	//		return 0, nil, errors.New("bad luck")
	//	}
	//	return 0, nil, nil
	//}
	//scanner.Split(split)
	//buf := make([]byte, 256)
	//scanner.Buffer(buf, bufio.MaxScanTokenSize)

	//var data bytes.Buffer
	//reader.Read()
	//max := 64
	//data := make([]byte, max)
	//dd := make([]byte, 0)
	//for {
	//	dn, err := reader.Read(dd)
	//	log.Printf("%v", dn)
	//
	//n, err := reader.Read(data)
	//	//data, err := ioutil.ReadAll(reader)
	//	if err != nil {
	//		log.Println(n, err)
	//	}
	//	go ctx.HandleRead(data[:])
	//	if n < max {

	//		break
	//	}
	//}
	//io.ReadAl
	//bs := make([]byte, 1000)
	//reader.Read(bs)

	//ctx.HandleRead(data)
}

//func (p *ByteArrayDecoder) HandleWrite(ctx netty.OutboundContext, message netty.Message) {
//	out := new(bytes.Buffer)
//	ctx.HandleWrite(out)
//}
