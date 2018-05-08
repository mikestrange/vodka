package proxy

import "fmt"
import "ants/base"
import "ants/lib/gsql"

func init_uids() {
	conn := gsql.NewConn()
	conn.Debug()
	conn.Exec("DELETE FROM texas_game.user")
	//
	RandUids(conn)
}

func init() {
	//init_uids()
	conn := gsql.NewConn()
	conn.Debug()
	fmt.Println(base.TryFun(func() {
		for i := 0; i < 10; i++ {
			//go regUser(conn, fmt.Sprintf("qwe%d", 10+i))
		}
	}))
}
