package gee

import (
	"net/http"
)

type HandlerFunc func(*Context)

type Engine struct {
	router *router
}

func New() *Engine {
	return &Engine{router: newRouter()}
}

func (engine *Engine) addRouter(method string, path string, handler HandlerFunc) {
	engine.router.addRouter(method, path, handler)
}

func (engine *Engine) Get(path string, handler HandlerFunc) {
	engine.addRouter("GET", path, handler)
}

func (engine *Engine) POST(path string, handler HandlerFunc) {
	engine.addRouter("POST", path, handler)
}

func (engine *Engine) Run(path string) (err error) {
	return http.ListenAndServe(path, engine)
}

func (engine *Engine) ServeHTTP(response http.ResponseWriter, r *http.Request) {
	//key := r.Method + "-" + r.URL.Path
	//if handler, ok := engine.router[key]; ok {
	//	handler(response, r)
	//} else {
	//	fmt.Fprintf(response, "404 NOT FOUND %s\n", r.URL)
	//}
	c := newContext(response, r)
	engine.router.Handle(c)
}
