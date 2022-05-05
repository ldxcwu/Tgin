package main

import (
	"log"
	"net/http"
	"tgin"
	"time"
)

func onlyForG2() tgin.HandlerFunc {
	return func(ctx *tgin.Context) {
		t := time.Now()
		ctx.Fail(500, "Internal Server Error")
		log.Printf("[%d] %s in %v for group v2", ctx.StatusCode, ctx.Req.RequestURI, time.Since(t))
	}
}

func main() {

	engine := tgin.New()
	engine.Use(func(ctx *tgin.Context) {
		// Start timer
		t := time.Now()
		// Process request
		ctx.Next()
		// Calculate resolution time
		log.Printf("[%d] %s in %v", ctx.StatusCode, ctx.Req.RequestURI, time.Since(t))
	})
	engine.GET("/index", func(ctx *tgin.Context) {
		ctx.HTML(http.StatusOK, "<h1>Index page</h1>")
	})

	g1 := engine.Group("g1")
	{
		g1.GET("/", func(ctx *tgin.Context) {
			ctx.HTML(http.StatusOK, "<h1>Hello Tgin</h1>")
		})

		g1.GET("/hello", func(ctx *tgin.Context) {
			ctx.String(http.StatusOK, "Hello %s!, you're at %s\n", ctx.Query("name"), ctx.Path)
		})
	}

	g2 := engine.Group("g2")
	g2.Use(onlyForG2())
	{
		g2.GET("/hello/:name", func(ctx *tgin.Context) {
			ctx.String(http.StatusOK, "hello %s!, you're at %s...\n", ctx.Param("name"), ctx.Path)
		})

		g2.GET("/asserts/*filepath", func(ctx *tgin.Context) {
			ctx.JSON(http.StatusOK, tgin.H{
				"filepath": ctx.Param("filepath"),
			})
		})
	}

	engine.POST("/login", func(ctx *tgin.Context) {
		ctx.JSON(http.StatusOK, tgin.H{
			"username": ctx.PostForm("username"),
			"password": ctx.PostForm("password"),
		})
	})

	engine.Run(":9999")
}
