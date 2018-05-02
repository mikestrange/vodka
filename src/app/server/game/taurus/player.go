package taurus

type Player struct {
	uid     int
	gate    int    //客户端
	session uint64 //回话
	name    string
	gift    int
	state   int
}

func newPlayer(uid int, gate int, session uint64) *Player {
	return &Player{uid: uid, gate: gate, session: session, name: "", gift: 0}
}

//重进刷新
func (this *Player) update(gate int, session uint64) {
	this.gate = gate
	this.session = session
	this.state = 0
}

//在游戏的时候设置
func (this *Player) over_leave() {
	this.state = STATE_OFFLINE
}

func (this *Player) isLeave() bool {
	return this.state == STATE_OFFLINE
}
