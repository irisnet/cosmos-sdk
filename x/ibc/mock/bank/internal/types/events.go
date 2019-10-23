package types

const (
	AttributeKeyDestPort      = "dest_port"
	AttributeKeyDestChannelID = "dest_channel_id"
)

var (
	EventTypeRecvTransferPacket = MsgRecvTransferPacket{}.Type()
)
