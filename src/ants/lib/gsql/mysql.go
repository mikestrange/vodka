package gsql

import (
	"database/sql"
	"fmt"
	//第三方
	_ "github.com/go-sql-driver/mysql"
)

type IConn interface {
	//打开链接
	Connect(string, int) bool
	//关闭
	Close() error
	//打开debug
	Debug()
	//具体链接
	Conn() *sql.DB
	//地址
	Addr() string
	//action
	Begin() (*sql.Tx, bool)
	//一般用于update, insert, delete
	Exec(string, ...interface{}) (sql.Result, bool)
	//一般用于select
	Query(string, ...interface{}) (*sql.Rows, bool)
	//一般用于select获取一条
	QueryRow(string, ...interface{}) *sql.Row
	//方案1 表单(Query)
	Form(string, ...interface{}) IRowsForm
	//方案2	简单(Query)
	Result(string, ...interface{}) []map[string]interface{}
	//预处理>不明觉厉
	Prepare(string) (*sql.Stmt, bool)
}

type sql_client struct {
	conn  *sql.DB
	addr  string
	debug bool
}

//本地
func NewConn() IConn {
	this := new(sql_client)
	this.Connect(ToAddr("root", "127.0.0.1:3306", "123456", "user_info"), 20)
	return this
}

func NewConnAddr(addr string, size int) (IConn, bool) {
	this := new(sql_client)
	return this, this.Connect(addr, size)
}

func (this *sql_client) Conn() *sql.DB {
	return this.conn
}

func (this *sql_client) Debug() {
	this.debug = true
}

func (this *sql_client) Connect(addr string, maxOpen int) bool {
	if maxOpen < MysqlMinOpen {
		maxOpen = MysqlMaxOpen
	}
	conn, err := sql.Open("mysql", addr)
	if err != nil {
		fmt.Println("Mysql Err:", addr)
		return false
	}
	fmt.Println("Mysql Ok:", addr)
	this.addr = addr
	this.conn = conn
	this.conn.SetMaxOpenConns(maxOpen)
	this.conn.SetMaxIdleConns(maxOpen / 3)
	return true
}

func (this *sql_client) Close() error {
	return this.conn.Close()
}

func (this *sql_client) Addr() string {
	return this.addr
}

//static
func ToAddr(root string, host string, pwd string, table string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", root, pwd, host, table)
}

//环境，多句并行，集体合成
func (this *sql_client) Begin() (*sql.Tx, bool) {
	tx, err := this.conn.Begin()
	if err == nil {
		return tx, true
	}
	fmt.Println("Tx Err:", err)
	return nil, false
}

func (this *sql_client) Exec(sql string, args ...interface{}) (sql.Result, bool) {
	ret, err := this.conn.Exec(sql, args...)
	if err == nil {
		if this.debug {
			fmt.Println("Exec ok:", sql)
		}
		return ret, true
	}
	fmt.Println("Exec Err:", err, ",sql=", sql)
	return nil, false
}

func (this *sql_client) Query(sql string, args ...interface{}) (*sql.Rows, bool) {
	ret, err := this.conn.Query(sql, args...)
	if err == nil {
		if this.debug {
			fmt.Println("Query ok:", sql)
		}
		return ret, true
	}
	fmt.Println("Query Err:", err)
	return nil, false
}

func (this *sql_client) Form(sql string, args ...interface{}) IRowsForm {
	if rows, ok := this.Query(sql, args...); ok {
		return toForm(rows)
	}
	return newForm()
}

func (this *sql_client) Result(sql string, args ...interface{}) []map[string]interface{} {
	if rows, ok := this.Query(sql, args...); ok {
		return toResult(rows)
	}
	return []map[string]interface{}{}
}

func (this *sql_client) QueryRow(sql string, args ...interface{}) *sql.Row {
	return this.conn.QueryRow(sql, args...)
}

func (this *sql_client) Prepare(sql string) (*sql.Stmt, bool) {
	stm, err := this.conn.Prepare(sql)
	if err == nil {
		if this.debug {
			fmt.Println("Prepare ok:", sql)
		}
		return stm, true
	}
	fmt.Println("Prepare Err:", err)
	return nil, false
}
