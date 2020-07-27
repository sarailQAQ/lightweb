package lightweb

import (
	"fmt"
	"github.com/sarailQAQ/MyGin/lightweb/render"
	"log"
	"math"
	"net/http"
	"net/url"
	"strings"
)

const abortIndex int8 = math.MaxInt8 / 2

type Context struct {
	Request        *http.Request
	Writer         http.ResponseWriter

	queryParam map[string]string
//	formParam map[string]string
	Params Params

	index int8
	handlers []Handler
	Keys map[string]interface{}

	engine *Engine

}

func (c *Context) String(s string) {
	_,_ = c.Writer.Write([]byte(s))
}

//下一个方法
func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

//停止请求
func (c *Context) Abort() {
	c.index = abortIndex
}


func (c *Context) Query(key string) string {
	v := c.queryParam[key]
	return v
}

func (c *Context) Render(code int,r render.Render) {
	c.Status(code)

	if !bodyAllowedForStatus(code) {
		r.WriteContentType(c.Writer)
		//c.w.WriteHeaderNow()
		return
	}

	if err := r.Render(c.Writer); err != nil {
		panic(err)
	}
}

//以JSON格式响应请求
func (c *Context) JSON(code int,obj interface{}){
	c.Render(code,render.JSON{Data:obj})
}

func (c *Context) Status(code int) {
	c.Writer.WriteHeader(code)
}

func (c *Context) Set(key string,val interface{}) {
	c.Keys[key] = val
}

func (c *Context) Get(key string) interface{} {
	return c.Keys[key]
}


func NewContext(rw http.ResponseWriter, r *http.Request,hs []Handler,engine *Engine) (ctx Context) {
	ctx = Context{
		Request:       r,
		Writer:         rw,
//		formParam: make(map[string]string),
		index:	   0,
		handlers:  hs,
		engine:    engine,
		Keys:	   make(map[string]interface{}),
	}

	ctx.queryParam = parseQuery(r.RequestURI)
	return
}

func parseQuery(uri string) (res map[string]string) {
	res = make(map[string]string)
	uris := strings.Split(uri, "?")
	if len(uris) == 1 {
		return
	}
	param := uris[len(uris)-1]
	pair := strings.Split(param, "&")

	for _, kv := range pair {
		kvPair := strings.Split(kv, "=")
		if len(kvPair) != 2 {
			fmt.Println(kvPair)
			panic("request error")
		}
		res[kvPair[0]] = kvPair[1]
	}
	return
}

func (c *Context) PostForm(key string) string {
	req := c.Request
	if err := req.ParseMultipartForm(c.engine.MaxMultipartMemory); err != nil {
		if err != http.ErrNotMultipart {
			log.Println(err)
		}
	}
	if values := req.PostForm[key]; len(values) > 0 {
		return values[0]
	}
	return ""
}

func bodyAllowedForStatus(status int) bool {
	switch {
	case status >= 100 && status <= 199:
		return false
	case status == http.StatusNoContent:
		return false
	case status == http.StatusNotModified:
		return false
	}
	return true
}

//设置cookie
func (c *Context) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	if path == "" {
		path = "/"
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}

//获取cookie的值
func (c *Context) Cookie(name string) (string, error) {
	cookie, err := c.Request.Cookie(name)
	if err != nil {
		return "", err
	}
	val, _ := url.QueryUnescape(cookie.Value)
	return val, nil
}

//获取请求头
func (c *Context) GetHeader(key string) string {
	return c.Request.Header.Get(key)
}

//以Key/Value格式设置请求头参数，若Value为空，则默认删除
func (c *Context) Header(key,val string) {
	if val == "" {
		c.Writer.Header().Del(key)
		return
	}
	c.Writer.Header().Set(key,val)
}

//按关键字查找动态路由
func (c *Context) Param(key string) string {
	return c.Params.ByName(key)
}