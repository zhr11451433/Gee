package gee

import (
	"html/template"
	"net/http"
	"path"
	"strings"
)

type HandlerFunc func(*Context)

type RouterGroup struct {
	prefix      string        //前缀码
	parent      *RouterGroup  //父分组
	middlewares []HandlerFunc //中间件
	engine      *Engine       //所有的分组都会指向唯一的engine实例
}
type Engine struct {
	*RouterGroup  //Engine内置了RouterGroup的所有方法
	router        *router
	groups        []*RouterGroup     //存储所有分组
	htmlTemplates *template.Template // for html render
	funcMap       template.FuncMap   // for html render
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
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

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (engine *Engine) ServeHTTP(response http.ResponseWriter, r *http.Request) {
	//key := r.Method + "-" + r.URL.Path
	//if handler, ok := engine.router[key]; ok {
	//	handler(response, r)
	//} else {
	//	fmt.Fprintf(response, "404 NOT FOUND %s\n", r.URL)
	//}
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(response, r)
	c.handlers = middlewares
	engine.router.Handle(c)
}

func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		// Check if file exists and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// serve static files
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET handlers
	group.GET(urlPattern, handler)
}
