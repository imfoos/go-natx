package protocol

type NatxMessage struct {
	Type     NatxMessageType
	MetaData map[string]interface{}
	Data     []byte
}
