package hippo

//处理者
type IHandle interface {
	//tag, target, data
	Handle(IEvent)
}

func NewHandle(block func(IEvent)) IHandle {
	return &BlockHandle{handle: block}
}

type BlockHandle struct {
	handle func(IEvent)
}

func (this *BlockHandle) Handle(event IEvent) {
	this.handle(event)
}
