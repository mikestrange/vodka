package gnet

//一些基础的消息
const (
	GNET_HEARTBEAT_PINT = 10 //心跳
)

//一些基础包
var PacketWithHeartBeat = NewPacketWithArgs(GNET_HEARTBEAT_PINT) //心跳
