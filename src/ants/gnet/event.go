package gnet

import "ants/gcode"

type SocketEvent struct {
	tx   Context
	pack gcode.ISocketPacket
}

func NewBytes(s Context, b []byte) *SocketEvent {
	return &SocketEvent{tx: s, pack: gcode.NewPackBytes(b)}
}

func NewPack(s Context, p interface{}) *SocketEvent {
	return &SocketEvent{tx: s, pack: p.(gcode.ISocketPacket)}
}

func (this *SocketEvent) Tx() Context {
	return this.tx
}

func (this *SocketEvent) Pack() gcode.ISocketPacket {
	return this.pack
}

func (this *SocketEvent) BeginPack() gcode.ISocketPacket {
	this.pack.ReadBegin()
	return this.pack
}
