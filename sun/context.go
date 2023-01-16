package sun

import "net/http"

// Context 路由上下文
type Context struct {
	// 请求响应
	W http.ResponseWriter
	// 请求
	R *http.Request
}
