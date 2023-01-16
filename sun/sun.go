package sun

import (
	"fmt"
	. "gwinter/utils"
	"log"
	"net/http"
	"strings"
)

// Application 框架服务引擎
//
// 所有功能都建立在引擎中
type Application struct {
	router
}

// HandlerFunc 请求处理函数
type HandlerFunc func(ctx *Context)

// router 请求路由
type router struct {
	// 分组路由信息
	groups []*routerGroup
}

// 路由组
type routerGroup struct {
	// 组名
	basePath string
	// 路由映射和处理函数
	handlerMapping map[string]HandlerFunc
	// 将路由根据不同的请求类型进行存储
	handlerMethod map[string][]string
}

// Group 创建路由组,并且将路由组返回
// name: 组名
func (r *router) Group(name string) *routerGroup {
	g := &routerGroup{basePath: name, handlerMapping: make(map[string]HandlerFunc), handlerMethod: make(map[string][]string)}
	r.groups = append(r.groups, g)
	return g
}

// BindHandler 绑定路由和处理函数
// method  请求类型
// path    请求路由
// handler 处理函数
func (r *routerGroup) BindHandler(method, path string, handler HandlerFunc) {
	IsEmpty(method, "请求类型不能为空!")
	IsEmpty(path, "path不能为空!")
	IsNil(handler, "处理函数不能为nil!")
	r.handlerMapping[path] = handler
	r.handlerMethod[method] = append(r.handlerMethod[method], r.basePath+path)
}

// Any 绑定任意类型路由和处理函数
func (r *routerGroup) Any(path string, handler HandlerFunc) {
	r.BindHandler("ANY", path, handler)
}

// GET 绑定GET类型路由和处理函数
func (r *routerGroup) GET(path string, handler HandlerFunc) {
	r.BindHandler(http.MethodGet, path, handler)
}

// New 创建服务引擎
func New() *Application {
	return &Application{
		router{},
	}
}

// 重写server的ServeHTTP接口
// 实现所有请求的函数代理调用
func (a *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, group := range a.groups {
		// 拼接本分组的Path前缀
		var groupPath string
		if strings.HasPrefix(group.basePath, "/") {
			groupPath = group.basePath
		} else {
			groupPath = "/" + group.basePath
		}
		// 在本分组内寻找是否有本次请求所要执行的路由和函数
		for path, handler := range group.handlerMapping {
			url := groupPath + path
			if url == r.RequestURI {
				// 匹配成功
				ctx := &Context{
					W: w,
					R: r,
				}
				// 判断AYN类型请求方式,如果有可以任意放行
				paths := group.handlerMethod["ANY"]
				if paths != nil {
					for _, p := range paths {
						if url == p {
							handler(ctx)
							return
						}
					}
				}
				// 匹配Method
				paths = group.handlerMethod[r.Method]
				if paths != nil {
					for _, p := range paths {
						if url == p {
							handler(ctx)
							return
						}
					}
				}
				w.Write([]byte("请求类型不支持."))
				return
			}
		}
	}
	w.Write([]byte("Not Found 404"))
	return
}

// Start 服务启动方法
func (a *Application) Start() {
	// 将所有请求转发到 Application的ServeHTTP方法中
	http.Handle("/", a)
	// 启动HTTP服务
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Server Run Success Port: %d \n")
}
