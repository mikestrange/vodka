package texas

type Player struct {
	uid     int
	gate    int
	session uint64
	name    string
	gift    int
	state   int //状态
}

func newPlayer(uid int, name string, gift int) *Player {
	return &Player{uid: uid, name: name, gift: gift}
}

func (this *Player) update(name string, gift int) {
	this.name = name
	this.gift = gift
	//this.gate = gate
	//this.session = session
}
