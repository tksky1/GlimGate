package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tksky1/glimgate/internal/service"
	"github.com/tksky1/glimgate/pkg/response"
)

// DirectionAPI 方向API处理器
type DirectionAPI struct {
	directionService *service.DirectionService
}

// NewDirectionAPI 创建方向API实例
func NewDirectionAPI() *DirectionAPI {
	return &DirectionAPI{
		directionService: service.NewDirectionService(),
	}
}

// CreateDirection 创建方向（管理员）
// @Summary 创建方向
// @Description 管理员创建新的方向
// @Tags 方向管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body service.CreateDirectionRequest true "方向信息"
// @Success 200 {object} response.Response{data=model.Direction} "创建成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "权限不足"
// @Router /api/admin/directions [post]
func (a *DirectionAPI) CreateDirection(c *gin.Context) {
	var req service.CreateDirectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBindError)
		return
	}

	direction, err := a.directionService.CreateDirection(&req)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, direction)
}

// GetDirections 获取方向列表
// @Summary 获取方向列表
// @Description 获取所有方向的列表
// @Tags 方向管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]model.Direction} "获取成功"
// @Failure 500 {object} response.Response "内部错误"
// @Router /api/directions [get]
func (a *DirectionAPI) GetDirections(c *gin.Context) {
	directions, err := a.directionService.GetDirections()
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, directions)
}

// GetDirection 获取方向详情
// @Summary 获取方向详情
// @Description 根据ID获取方向的详细信息
// @Tags 方向管理
// @Accept json
// @Produce json
// @Param id path int true "方向ID"
// @Success 200 {object} response.Response{data=model.Direction} "获取成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 404 {object} response.Response "方向不存在"
// @Router /api/directions/{id} [get]
func (a *DirectionAPI) GetDirection(c *gin.Context) {
	directionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.CodeInvalidParams)
		return
	}

	direction, err := a.directionService.GetDirectionByID(uint(directionID))
	if err != nil {
		if err.Error() == "方向不存在" {
			response.Error(c, response.CodeDirectionNotFound)
			return
		}
		response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, direction)
}

// UpdateDirection 更新方向（管理员）
// @Summary 更新方向
// @Description 管理员更新方向信息
// @Tags 方向管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "方向ID"
// @Param request body service.UpdateDirectionRequest true "更新信息"
// @Success 200 {object} response.Response{data=model.Direction} "更新成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "权限不足"
// @Failure 404 {object} response.Response "方向不存在"
// @Router /api/admin/directions/{id} [put]
func (a *DirectionAPI) UpdateDirection(c *gin.Context) {
	directionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.CodeInvalidParams)
		return
	}

	var req service.UpdateDirectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBindError)
		return
	}

	direction, err := a.directionService.UpdateDirection(uint(directionID), &req)
	if err != nil {
		if err.Error() == "方向不存在" {
			response.Error(c, response.CodeDirectionNotFound)
			return
		}
		response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, direction)
}

// DeleteDirection 删除方向（管理员）
// @Summary 删除方向
// @Description 管理员删除方向
// @Tags 方向管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "方向ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "权限不足"
// @Failure 404 {object} response.Response "方向不存在"
// @Router /api/admin/directions/{id} [delete]
func (a *DirectionAPI) DeleteDirection(c *gin.Context) {
	directionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.CodeInvalidParams)
		return
	}

	if err := a.directionService.DeleteDirection(uint(directionID)); err != nil {
		if err.Error() == "方向不存在" {
			response.Error(c, response.CodeDirectionNotFound)
			return
		}
		response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, nil)
}