# 前缀树路由
此前的路由都是静态路由，指定路径的映射  
但是有时候我们想要实现模糊匹配，例如正则匹配等路由  
这里实现前缀树路由，主要包括 `':'` 以及 `'*'` 匹配
```
/hello/:name  <= /hello/lxw  /hello/cc ...
/asserts/*filepath <= /asserts/css/cmh.png ...
```
## 1. 定义前缀树节点
```go
type node struct {
	pattern  string //该节点所匹配的路由路径
	part     string //该节点当前路径值（路由路径的一部分）
	children []*node
	isWild   bool //part 含：或 *时为true
}
```
### 1.1 前缀树 节点匹配
```go
//查找node的children中是否满足匹配path，返回第一个匹配的，用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

//返回所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}
```
### 1.2 前缀树 节点插入/搜索
> 搜索：给定 /assert/css/cmh.png 可以匹配到 /assert/*filepath的节点，然后返回该节点，该节点的pattern属性表示路由路径，最后根据该路由路径去查路由map，然后进行服务
```go
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	children := n.matchChildren(part)
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
```
## 2. 修改路由对象
```go
type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}
```
### 2.1 添加路由信息
> 在路由map中添加路由信息，同时记录走前缀树里
```go
func parsePattern(pattern string) []string {
	paths := strings.Split(pattern, "/")
	s := make([]string, 0)
	for _, path := range paths {
		if path != "" {
			s = append(s, path)
			if path[0] == '*' {
				break
			}
		}
	}
	return s
}

func (r *router) addRoute(method, pattern string, handler HandlerFunc) {
	paths := parsePattern(pattern)
	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, paths, 0)
	r.handlers[key] = handler
}

//tgin.go
func (e *Engine) addRoute(method, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	e.router.addRoute(method, pattern, handler)
}
```
### 2.2 获取路由信息
```go
func (r *router) getRoute(method, pattern string) (*node, map[string]string) {
	searchPaths := parsePattern(pattern)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
    //解析出前缀树节点，就可以根据节点的pattern属性进行查表提供服务
	n := root.search(searchPaths, 0)
	if n != nil {
        //这里再进一步解析路径的映射关系
        //例如 /assert/css/cmh.png => /assert/*filepath
        //解析出 params[filepath] = "css/cmh.png"
		paths := parsePattern(n.pattern)
		for index, path := range paths {
			if path[0] == ':' {
				params[path[1:]] = searchPaths[index]
			}
			if path[0] == '*' && len(path) > 1 {
				params[path[1:]] = strings.Join(searchPaths[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}

```

## 3. 测试
```go
package main

import (
	"net/http"
	"tgin"
)

func main() {

	engine := tgin.New()

	engine.GET("/", func(ctx *tgin.Context) {
		ctx.HTML(http.StatusOK, "<h1>Hello Tiny Gin</h1>")
	})

	engine.GET("/hello", func(ctx *tgin.Context) {
		ctx.String(http.StatusOK, "Hello %s!, you're at %s\n", ctx.Query("name"), ctx.Path)
	})

	engine.GET("/hello/:name", func(ctx *tgin.Context) {
		ctx.String(http.StatusOK, "hello %s!, you're at %s...\n", ctx.Param("name"), ctx.Path)
	})

	engine.GET("/asserts/*filepath", func(ctx *tgin.Context) {
		ctx.JSON(http.StatusOK, tgin.H{
			"filepath": ctx.Param("filepath"),
		})
	})

	engine.POST("/login", func(ctx *tgin.Context) {
		ctx.JSON(http.StatusOK, tgin.H{
			"username": ctx.PostForm("username"),
			"password": ctx.PostForm("password"),
		})
	})

	engine.Run(":9999")
}
```
```shell
➜  main git:(master) ✗ curl "http://localhost:9999/hello?name=lxw"
Hello lxw!, you're at /hello
➜  main git:(master) ✗ curl "http://localhost:9999/hello/lxw"
hello lxw!, you're at /hello/lxw...
➜  main git:(master) ✗ curl "http://localhost:9999/asserts/css/cmh.png"
{"filepath":"css/cmh.png"}
```
