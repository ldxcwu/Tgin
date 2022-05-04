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
