package main

import (
	"gwinter/sun"
)

// Application Start

func main() {
	// 创建服务
	app := sun.New()
	group := app.Group("/v1")
	group.Any("/test", func(ctx *sun.Context) {
		ctx.W.Write([]byte("123"))
	})
	//group.BindHandler("test", func(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte("测试分组"))
	//})
	// 绑定请求处理函数
	//app.BindHandler("/hello", func(w http.ResponseWriter, r *http.Request) {
	//		w.Write([]byte("Hello"))
	//})
	// 启动服务
	app.Start()
}
