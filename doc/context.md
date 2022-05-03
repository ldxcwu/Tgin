# Context
> 对于Web服务来说，无非是根据`*http.Request` 构造 `http.ResponseWriter` 
> 但是这两个对象的接口粒度太细，这样构造一个完整的响应需要考虑消息头`Header`和消息体`Body`
> 而`Header`又包含状态码`StatusCode`,消息类型`Content-Type`等，这些几乎每次请求都需要设置，若不进行有效的封装，将出现大量的重复代码，而且易出错。
## 1. 封装前后比较
### 1.1 封装前
```go
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
### 1.2 封装后
```go
func main() {

	engine := tgin.New()

	engine.GET("/", func(ctx *tgin.Context) {
		ctx.HTML(http.StatusOK, "<h1>Hello Tiny Gin</h1>")
	})

	engine.GET("/hello", func(ctx *tgin.Context) {
		ctx.String(http.StatusOK, "Hello %s!, you're at %s\n", ctx.Query("name"), ctx.Path)
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
可见，封装后，可以直接写入消息主体内容，不同类型的消息也都提供了对应的接口
## 2. 定义Context
将请求和响应进行封装，并直接暴露添加一些基本属性
```go
type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	StatusCode int
}
```
### 2.1 提供一些获得/写入基本属性的接口
```go
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

func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}
```
### 2.2 提供不同类型消息的写入接口
```go
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
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
```
## 3. 定义路由结构
将原先`Engine`中的路由`map`封装为一个结构体   
修改`HandlerFunc`类型
```go
type HandlerFunc func(*Context)

type Engine struct {
	router *router
}

type router struct {
	handlers map[string]HandlerFunc
}
```
## 4. 测试
```shell
Tgin git:(master) ✗ curl -i http://localhost:9999                                       
HTTP/1.1 200 OK
Content-Length: 23
Connection: keep-alive
Content-Type: text/html
Date: Tue, 03 May 2022 07:58:26 GMT
Keep-Alive: timeout=4
Proxy-Connection: keep-alive

<h1>Hello Tiny Gin</h1>%                                                                                                                         
➜  Tgin git:(master) ✗ curl "http://localhost:9999/hello/hello?name=lxw"
404 NOT FOUND: /hello/hello
➜  Tgin git:(master) ✗ curl "http://localhost:9999/hello?name=lxw" 
Hello lxw!, you're at /hello
➜  Tgin git:(master) ✗ curl "http://localhost:9999/login" -X POST -d "username=1&password=1"
{"password":"1","username":"1"}
```