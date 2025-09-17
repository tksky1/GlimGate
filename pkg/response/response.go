package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code int         `json:"code" example:"0"`
	Msg  string      `json:"msg" example:"success"`
	Data interface{} `json:"data,omitempty"`
}

// 响应码定义
const (
	CodeSuccess = 0
	CodeError   = 1

	// 用户相关错误码
	CodeUserNotFound     = 1001
	CodeUserExists       = 1002
	CodeInvalidPassword  = 1003
	CodeUnauthorized     = 1004
	CodeForbidden        = 1005
	CodeInvalidToken     = 1006

	// 题目相关错误码
	CodeDirectionNotFound = 2001
	CodeProblemNotFound   = 2002
	CodeSubmissionNotFound = 2003

	// 参数错误码
	CodeInvalidParams = 3001
	CodeBindError     = 3002

	// 系统错误码
	CodeDatabaseError = 5001
	CodeInternalError = 5002
)

// 错误消息映射
var codeMsg = map[int]string{
	CodeSuccess: "success",
	CodeError:   "error",

	CodeUserNotFound:     "用户不存在",
	CodeUserExists:       "用户已存在",
	CodeInvalidPassword:  "密码错误",
	CodeUnauthorized:     "未授权",
	CodeForbidden:        "权限不足",
	CodeInvalidToken:     "无效的token",

	CodeDirectionNotFound:  "方向不存在",
	CodeProblemNotFound:    "题目不存在",
	CodeSubmissionNotFound: "提交不存在",

	CodeInvalidParams: "参数错误",
	CodeBindError:     "参数绑定失败",

	CodeDatabaseError: "数据库错误",
	CodeInternalError: "内部错误",
}

// GetMsg 获取错误消息
func GetMsg(code int) string {
	if msg, ok := codeMsg[code]; ok {
		return msg
	}
	return "未知错误"
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: CodeSuccess,
		Msg:  GetMsg(CodeSuccess),
		Data: data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  GetMsg(code),
	})
}

// ErrorWithMsg 带自定义消息的错误响应
func ErrorWithMsg(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
	})
}

// Unauthorized 未授权响应
func Unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, Response{
		Code: CodeUnauthorized,
		Msg:  GetMsg(CodeUnauthorized),
	})
}

// Forbidden 权限不足响应
func Forbidden(c *gin.Context) {
	c.JSON(http.StatusForbidden, Response{
		Code: CodeForbidden,
		Msg:  GetMsg(CodeForbidden),
	})
}