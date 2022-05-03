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

	engine.POST("/login", func(ctx *tgin.Context) {
		ctx.JSON(http.StatusOK, tgin.H{
			"username": ctx.PostForm("username"),
			"password": ctx.PostForm("password"),
		})
	})

	engine.Run(":9999")
}
