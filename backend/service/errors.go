package service

// AppError 自定义错误类型
type AppError struct {
	Code    string
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}
