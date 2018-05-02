package ghttp

//异步高并发http
import "net/http"
import "ants/base"
import "ants/core"
import "io/ioutil"
import "strings"
import "net/url"

type Conn struct {
	core.Box
}

func (this *Conn) OnReady() {
	this.SetName("Http程序")
	this.SetAgent(this)
}

func (this *Conn) Handle(event interface{}) {
	this.on_http_request(event.(*HttpEvent))
}

//send
func (this *Conn) on_http_request(item *HttpEvent) {
	switch item.method {
	case GET:
		get_http(item)
	case POST:
		post_http(item)
	case FORM:
		form_http(item)
	default:
		get_http(item)
	}
	this.PushSuper(item)
}

func (this *Conn) Get(url string, data interface{}, callback interface{}) {
	this.Push(&HttpEvent{url: url, method: GET, data: data, call: callback})
}

func (this *Conn) Post(url string, val url.Values, data interface{}, callback interface{}) {
	this.Push(&HttpEvent{url: url, method: POST, vals: val, data: data, call: callback})
}

func (this *Conn) Form(url string, val url.Values, data interface{}, callback interface{}) {
	this.Push(&HttpEvent{url: url, method: FORM, vals: val, data: data, call: callback})
}

//=====================static==================
func http_result(ret *http.Response, item *HttpEvent) {
	defer ret.Body.Close()
	body, err := ioutil.ReadAll(ret.Body)
	if err == nil {
		item.body = body
	} else {
		item.err = err
	}
}

//get提交
func get_http(item *HttpEvent) {
	//http get
	ret, err := http.Get(item.url)
	if err == nil {
		http_result(ret, item)
	} else {
		item.err = err
	}
}

//post提交
func post_http(item *HttpEvent) {
	data := item.vals.(url.Values).Encode()
	req, err := http.NewRequest("POST", item.url, strings.NewReader(data))
	if err != nil {
		item.err = err
	} else {
		//必须设置
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		//发起请求
		ret, err := http.DefaultClient.Do(req)
		if err == nil {
			http_result(ret, item)
		} else {
			item.err = err
		}
	}
}

//表单提交
func form_http(item *HttpEvent) {
	data := item.vals.(url.Values)
	ret, err := http.PostForm(item.url, data)
	if err == nil {
		http_result(ret, item)
	} else {
		item.err = err
	}
}

func init() {
	conn := new(Conn)
	core.Main().Join(101, conn)
	//	var body = url.Values{}
	//	body.Add("handle", "sum_pay")
	//	body.Add("begin", "2017-11-19")
	//	body.Add("end", "2017-11-20")
	//test
	str := base.TryFun(func() {
		for i := 0; i < 1; i++ {
			conn.Get("http://localhost/phprpc/index.php", nil, func(str string) {
				println(str)
			})
		}
	})
	println(str)
}
