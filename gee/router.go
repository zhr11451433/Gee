package gee

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node       //每种http方法对应一棵trie数 roots["Get"]
	handlers map[string]HandlerFunc // 路由映射表：key="方法-路径"，value=处理函数
}

func newRouter() *router {
	return &router{roots: make(map[string]*node), handlers: make(map[string]HandlerFunc)}
}

// 切割路径函数
func parsePattern(pattern string) []string {
	parts := strings.Split(pattern, "/")
	result := make([]string, 0)
	for _, part := range parts {
		if part != "" {
			result = append(result, part)
			if part[0] == '*' { //遇见通配符就停止，后面所有路径都要算上
				break
			}
		}
	}
	return result
}

func (r *router) addRouter(method string, pattern string, handler HandlerFunc) {
	//log.Printf("Route %4s - %s", method, path)
	path := parsePattern(pattern)
	key := method + "-" + pattern
	_, ok := r.roots[method] //判断是否存在
	if !ok {
		r.roots[method] = &node{} //创建
	}
	r.roots[method].insert(pattern, path, 0) //插入
	r.handlers[key] = handler
}

//注册路由：/p/:lang/doc
//请求路径：/p/go/doc

func (r *router) getRouter(method, pattern string) (*node, map[string]string) {
	//1.把请求路径分割
	parts := parsePattern(pattern)

	//2.准备一个map用于存储
	mapPattern := make(map[string]string)

	//3.拿到该HTTP方法的叶子节点，并判断是否存在
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	//4.在前缀树里查找
	trie := root.search(parts, 0)
	if trie != nil {
		splitPattern := parsePattern(trie.pattern)
		for i, s := range splitPattern {
			if s[0] == ':' {
				mapPattern[s[1:]] = parts[i]
			}
			if s[0] == '*' && len(s) > 1 {
				mapPattern[s[1:]] = strings.Join(parts[i:], "/")
				break
			}
		}
		return trie, mapPattern
	}
	return nil, nil
}

func (r *router) Handle(c *Context) { //负责接收一个请求，找到匹配的路由规则，然后调用对应的处理函数
	n, params := r.getRouter(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
		// 第一步：从 handlers 里取出处理函数
		//handler := r.handlers[key]
		// 第二步：调用这个函数，把 Context 传进去
		//handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
