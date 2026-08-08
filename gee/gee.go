package gee

import (
	"net/http"
)

type HandlerFunc func(*Context)

type RouterGroup struct {
	prefix      string         //前缀码
	parent      *RouterGroup   //父分组
	middlewares []*HandlerFunc //中间件
	engine      *Engine        //所有的分组都会指向唯一的engine实例
}
type Engine struct {
	*RouterGroup //Engine内置了RouterGroup的所有方法
	router       *router
	groups       []*RouterGroup //存储所有分组
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	newGroup := &RouterGroup{
		engine: group.engine, //共享一个引擎
		parent: group,
		prefix: group.prefix + prefix, //当前前缀+父前缀
	}
	group.engine.groups = append(group.engine.groups, newGroup)
	return newGroup
}

// v1.GET("/user", handler)
// group 是 v1 这个分组，它的 prefix 是 "/api/v1"。
// comp 是 "/user"。
// pattern = "/api/v1" + "/user" = "/api/v1/user
func (group *RouterGroup) addRouter(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	group.engine.router.addRouter(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRouter("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRouter("POST", pattern, handler)
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
