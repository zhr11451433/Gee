package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	Req        *http.Request
	Writer     http.ResponseWriter
	StatusCode int
	Path       string
	Method     string
	Params     map[string]string
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Req:    req,
		Writer: w,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) { //设置响应头//
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.Status(code)
	c.SetHeader("Content-Type", "text/plain")
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)        //创建一个 JSON 编码器
	if err := encoder.Encode(obj); err != nil { //.Encode(obj)：把传入的 obj（可以是 map、结构体）Go 对象序列化成 JSON 字符串，直接输出到 c.Writer（发给浏览器 /curl）
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}
