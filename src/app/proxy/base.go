package proxy

import "ants/base"
import "ants/lib/gsql"

//库
type DataBase struct {
	conn gsql.IConn
	name string
}

func NewBase(conn gsql.IConn, name string) *DataBase {
	return &DataBase{conn: conn, name: name}
}

func (this *DataBase) Name() string {
	return this.name
}

//建库
func (this *DataBase) Create() {
	this.conn.Exec(base.Format("CREATE SCHEMA IF NOT EXISTS %s DEFAULT CHARACTER SET utf8", this.name))
}

//删库跑路
func (this *DataBase) Drop() {
	this.conn.Exec(base.Format("DROP DATABASE IF EXISTS %s", this.name))
}

func (this *DataBase) NewTable(name string) *Table {
	return &Table{name: this.name + "." + name, conn: this.conn}
}

//表
type Table struct {
	conn gsql.IConn
	name string
}

func NewTable(conn gsql.IConn, base string, name string) *Table {
	return &Table{conn: conn, name: base + "." + name}
}

func (this *Table) Name() string {
	return this.name
}

//建立表
func (this *Table) Create(keys ...string) {
	size := len(keys)
	str := ""
	for i := range keys {
		if i == size-1 {
			str += keys[i]
		} else {
			str += keys[i] + ","
		}
	}
	this.conn.Exec(base.Format("CREATE TABLE IF NOT EXISTS %s (%s)", this.name, str))
}

//删表跑路
func (this *Table) Drop() {
	this.conn.Exec(base.Format("DROP TABLE IF EXISTS %s", this.name))
}
