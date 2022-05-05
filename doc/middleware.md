# 中间件

责任链模式，使在处理路由映射请求之前先执行一些中间件服务

由于框架的请求处理作用对象是Context，
用户只需要提供下述方法即可
```go
func(c *Context) {}
```
并将这些中间件服务串联于Context对象中
```go
type Context struct {
	...
	handlers   []HandlerFunc
	index      int
}
```
最后提供服务的时候
```go
func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
        //将真正的路由处理添加至链的末尾，最后处理
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(ctx *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}
```