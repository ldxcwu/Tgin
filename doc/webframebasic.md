# Web框架雏形

## 1. 最基本的组成是路由映射，使用map实现
```go
type Engine struct {
	router map[string]http.HandlerFunc
}
```
## 2. 提供添加路由的功能
```go
func New() *Engine {
	return &Engine{router: make(map[string]http.HandlerFunc)}
}

func (e *Engine) addRoute(method, pattern string, handler http.HandlerFunc) {
	key := method + "-" + pattern
	e.router[key] = handler
}

func (e *Engine) GET(pattern string, handler http.HandlerFunc) {
	e.addRoute("GET", pattern, handler)
}

func (e *Engine) POST(pattern string, handler http.HandlerFunc) {
	e.addRoute("POST", pattern, handler)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := e.router[key]; ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}
```
## 3. 启动框架的入口
```go
func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}
```
## 4. 测试
```go
package main

import (
	"fmt"
	"net/http"
	"tgin"
)

func main() {
	engine := tgin.New()
	engine.GET("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
	})

	engine.GET("/hello", func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})

	engine.Run(":9999")
}
```

## 5. 结果

http://localhost:9999/hello

Header["Connection"] = ["keep-alive"]
Header["Sec-Ch-Ua-Mobile"] = ["?0"]
Header["Accept"] = ["text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"]
Header["Accept-Encoding"] = ["gzip, deflate, br"]
Header["Sec-Ch-Ua"] = ["\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"101\", \"Google Chrome\";v=\"101\""]
Header["Sec-Ch-Ua-Platform"] = ["\"macOS\""]
Header["Upgrade-Insecure-Requests"] = ["1"]
Header["Sec-Fetch-User"] = ["?1"]
Header["Accept-Language"] = ["zh-CN,zh;q=0.9"]
Header["Sec-Fetch-Site"] = ["none"]
Header["Sec-Fetch-Mode"] = ["navigate"]
Header["User-Agent"] = ["Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.41 Safari/537.36"]
Header["Sec-Fetch-Dest"] = ["document"]