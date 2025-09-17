package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tksky1/glimgate/internal/service"
	"github.com/tksky1/glimgate/pkg/response"
)

// UserAPI 用户API处理器
type UserAPI struct {
	userService *service.UserService
}

// NewUserAPI 创建用户API实例
func NewUserAPI() *UserAPI {
	return &UserAPI{
		userService: service.NewUserService(),
	}
}

// Register 用户注册
// @Summary 用户注册
// @Description 用户注册接口
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body service.RegisterRequest true "注册信息"
// @Success 200 {object} response.Response{data=model.User} "注册成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 500 {object} response.Response "内部错误"
// @Router /api/auth/register [post]
func (a *UserAPI) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBindError)
		return
	}

	user, err := a.userService.Register(&req)
	if err != nil {
		if err.Error() == "用户名已存在" {
			response.Error(c, response.CodeUserExists)
			return
		}
		response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, user)
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录接口
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body service.LoginRequest true "登录信息"
// @Success 200 {object} response.Response{data=service.LoginResponse} "登录成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "认证失败"
// @Router /api/auth/login [post]
func (a *UserAPI) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBindError)
		return
	}

	loginResp, err := a.userService.Login(&req)
	if err != nil {
		if err.Error() == "用户不存在" {
			response.Error(c, response.CodeUserNotFound)
			return
		}
		if err.Error() == "密码错误" {
			response.Error(c, response.CodeInvalidPassword)
			return
		}
		response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, loginResp)
}

// GetProfile 获取用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户的信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=model.User} "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "用户不存在"
// @Router /api/user/profile [get]
func (a *UserAPI) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	user, err := a.userService.GetUserByID(userID.(uint))
	if err != nil {
		if err.Error() == "用户不存在" {
			response.Error(c, response.CodeUserNotFound)
			return
		}
		response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, user)
}

// GetUsers 获取用户列表（管理员）
// @Summary 获取用户列表
// @Description 管理员获取用户列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=map[string]interface{}} "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "权限不足"
// @Router /api/admin/users [get]
func (a *UserAPI) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	users, total, err := a.userService.GetUsers(page, pageSize)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
		return
	}

	data := map[string]interface{}{
		"users": users,
		"total": total,
		"page":  page,
		"page_size": pageSize,
	}

	response.Success(c, data)
}

// GetUser 获取指定用户信息（管理员）
// @Summary 获取指定用户信息
// @Description 管理员获取指定用户的详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response{data=model.User} "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "权限不足"
// @Failure 404 {object} response.Response "用户不存在"
// @Router /api/admin/users/{id} [get]
func (a *UserAPI) GetUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.CodeInvalidParams)
		return
	}

	user, err := a.userService.GetUserByID(uint(userID))
	if err != nil {
		if err.Error() == "用户不存在" {
			response.Error(c, response.CodeUserNotFound)
			return
		}
		response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, user)
}

// UpdateUser 更新用户信息（管理员）
// @Summary 更新用户信息
// @Description 管理员更新用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Param request body service.UpdateUserRequest true "更新信息"
// @Success 200 {object} response.Response{data=model.User} "更新成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "权限不足"
// @Failure 404 {object} response.Response "用户不存在"
// @Router /api/admin/users/{id} [put]
func (a *UserAPI) UpdateUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.CodeInvalidParams)
		return
	}

	var req service.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBindError)
		return
	}

	user, err := a.userService.UpdateUser(uint(userID), &req)
	if err != nil {
		if err.Error() == "用户不存在" {
			response.Error(c, response.CodeUserNotFound)
			return
		}
		response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, user)
}

// DeleteUser 删除用户（管理员）
// @Summary 删除用户
// @Description 管理员删除用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "权限不足"
// @Failure 404 {object} response.Response "用户不存在"
// @Router /api/admin/users/{id} [delete]
func (a *UserAPI) DeleteUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.CodeInvalidParams)
		return
	}

	if err := a.userService.DeleteUser(uint(userID)); err != nil {
		if err.Error() == "用户不存在" {
			response.Error(c, response.CodeUserNotFound)
			return
		}
		response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, nil)
}