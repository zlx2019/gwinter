package utils

// Assert 断言
func Assert(expr bool, message string) {
	if expr {
		handler(message)
	}
}

// IsNil 如果为nil,直接panic
func IsNil(a any, message string) {
	if a == nil {
		handler(message)
	}
}

// IsEmpty 字符串非空断言
func IsEmpty(str, message string) {
	if str == "" {
		handler(message)
	}
}

// 错误日志记录
func handler(message string) {
	// log
	panic(message)
}
