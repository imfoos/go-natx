package protocol

type NatxMessageType int32

const (
	REGISTER        NatxMessageType = 1
	REGISTER_RESULT NatxMessageType = 2
	CONNECTED       NatxMessageType = 3
	DISCONNECTED    NatxMessageType = 4
	DATA            NatxMessageType = 5
	KEEPALIVE       NatxMessageType = 6
)

func (p NatxMessageType) String() string {
	switch p {
	case REGISTER:
		return "REGISTER"
	case REGISTER_RESULT:
		return "REGISTER_RESULT"
	case CONNECTED:
		return "CONNECTED"
	case DISCONNECTED:
		return "DISCONNECTED"
	case DATA:
		return "DATA"
	case KEEPALIVE:
		return "KEEPALIVE"
	}
	return "NONE"
}
