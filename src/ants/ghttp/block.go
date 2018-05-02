package ghttp

import "ants/base"
import "ants/core"
import "fmt"

const (
	GET  = 1
	POST = 2
	FORM = 3
)

type HttpEvent struct {
	core.IBlock
	url  string      //调用地址
	err  interface{} //错误信息(返回)
	data interface{} //自定义数据(携带)
	body interface{} //返回数据([]byte)
	//private
	method int         //调用方法(get/post/form)
	vals   interface{} //携带的数据(url.Values)
	req    interface{} //发送自定义(如果它存在，先用它)
	call   interface{} //回调函数
}

func (this *HttpEvent) Url() string {
	return this.url
}

func (this *HttpEvent) Err() interface{} {
	return this.err
}

func (this *HttpEvent) IsErr() bool {
	return this.err != nil
}

func (this *HttpEvent) Data() interface{} {
	return this.data
}

func (this *HttpEvent) Body() interface{} {
	return this.body
}

func (this *HttpEvent) Byte() []byte {
	switch b := this.body.(type) {
	case []byte:
		return b
	case string:
		return []byte(b)
	}
	return []byte{}
}

func (this *HttpEvent) Str() string {
	switch b := this.body.(type) {
	case []byte:
		return string(b)
	case string:
		return b
	}
	return ""
}

func (this *HttpEvent) Json() map[string]interface{} {
	jsonMap := make(map[string]interface{})
	if err := base.JsonDecode(this.Byte(), jsonMap); err != nil {
		fmt.Println(err)
		return make(map[string]interface{})
	}
	return jsonMap
}

//执行
func (this *HttpEvent) Result() {
	switch f := this.call.(type) {
	case func():
		f()
	case func(string):
		f(this.Str())
	case func([]byte):
		f(this.Byte())
	case func(interface{}):
		f(this.body)
	case func(map[string]interface{}):
		f(this.Json())
	case func(*HttpEvent):
		f(this)
	default:
		println("http not callback")
	}
}
