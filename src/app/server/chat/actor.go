package chat

import "ants/actor"

type ChatActor struct {
	actor.BaseActor
}

func (this *ChatActor) init() actor.IActor {
	this.SetMaster(1, "chat", nil, nil)
	return this
}

func (this *ChatActor) OnMessage(args ...interface{}) {

}

func (this *ChatActor) OnClosed() {

}
