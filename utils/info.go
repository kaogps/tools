package utils

import "github.com/kdada/tinygo/util"

// SuccessReturn 请求成功返回值
type SuccessReturn struct {
	Code int
	Data interface{}
}

// FailReturn 请求失败返回值
type FailReturn struct {
	Code    int
	Message string
}

var (
	// 10000开头的错误，参数错误

	// FailReturnCode10000 参数错误
	FailReturnCode10000 = &FailReturn{10000, "参数错误，参数不符合接口要求！"}

	// 20000开头的错误，用户权限错误

	// FailReturnCode20001 用户未登录
	FailReturnCode20001 = &FailReturn{20001, "用户未登录，请先登录！"}
	// FailReturnCode20002 用户权限不足
	FailReturnCode20002 = &FailReturn{20002, "操作失败，用户权限不足！"}
	// FailReturnCode20003 在学习中心中,用户只有在未设置区域的情况下(省市均为0)才能设置区域
	FailReturnCode20003 = &FailReturn{20003, "用户设置区域失败(用户不能重复设置区域)"}

	// 50000开头的错误，服务器错误

	// FailReturnCode50000 服务器错误
	FailReturnCode50000 = &FailReturn{50000, "服务器内部错误！"}
)

// NewSuccessReturn 创建一个成功返回值
func NewSuccessReturn(data interface{}) *SuccessReturn {
	return &SuccessReturn{0, data}
}

// NewFailReturn 返回一个自定义错误
func NewFailReturn(data interface{}) *FailReturn {
	var msg, ok = data.(string)
	if ok {
		// 60000开头的错误，自定义错误
		return &FailReturn{60000, msg}
	}
	err, ok := data.(error)
	if ok {
		return &FailReturn{60000, err.Error()}
	}
	// 兼容tinygo.util.Error
	err2, ok := data.(util.Error)
	if ok {
		return &FailReturn{60000, err2.String()}
	}
	return &FailReturn{60001, "未知错误"}
}
