package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tksky1/glimgate/internal/service"
	"github.com/tksky1/glimgate/pkg/response"
)

// ProblemAPI 题目API处理器
type ProblemAPI struct {
	problemService   *service.ProblemService
	directionService *service.DirectionService
}

// NewProblemAPI 创建题目API实例
func NewProblemAPI() *ProblemAPI {
	return &ProblemAPI{
		problemService:   service.NewProblemService(),
		directionService: service.NewDirectionService(),
	}
}

// CreateProblem 创建题目（管理员）
// @Summary 创建题目
// @Description 管理员创建新的题目
// @Tags 题目管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body service.CreateProblemRequest true "题目信息"
// @Success 200 {object} response.Response{data=model.Problem} "创建成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "权限不足"
// @Router /api/admin/problems [post]
func (a *ProblemAPI) CreateProblem(c *gin.Context) {
	var req service.CreateProblemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBindError)
		return
	}

	// 检查用户是否为该方向的负责人
	userID, _ := c.Get("user_id")
	isAdmin, _ := c.Get("is_admin")
	
	if !isAdmin.(bool) {
		isManager, err := a.directionService.CheckDirectionManager(req.DirectionID, userID.(uint))
		if err != nil {
			response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
			return
		}
		if !isManager {
			response.Error(c, response.CodeForbidden)
			return
		}
	}

	problem, err := a.problemService.CreateProblem(&req)
	if err != nil {
		if err.Error() == "方向不存在" {
			response.Error(c, response.CodeDirectionNotFound)
			return
		}
		response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, problem)
}

// GetProblems 获取题目列表
// @Summary 获取题目列表
// @Description 获取题目列表，可按方向筛选
// @Tags 题目管理
// @Accept json
// @Produce json
// @Param direction_id query int false "方向ID"
// @Success 200 {object} response.Response{data=[]model.Problem} "获取成功"
// @Failure 500 {object} response.Response "内部错误"
// @Router /api/problems [get]
func (a *ProblemAPI) GetProblems(c *gin.Context) {
	directionID, _ := strconv.ParseUint(c.DefaultQuery("direction_id", "0"), 10, 32)

	problems, err := a.problemService.GetProblems(uint(directionID))
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, problems)
}

// GetProblem 获取题目详情
// @Summary 获取题目详情
// @Description 根据ID获取题目的详细信息
// @Tags 题目管理
// @Accept json
// @Produce json
// @Param id path int true "题目ID"
// @Success 200 {object} response.Response{data=model.Problem} "获取成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 404 {object} response.Response "题目不存在"
// @Router /api/problems/{id} [get]
func (a *ProblemAPI) GetProblem(c *gin.Context) {
	problemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.CodeInvalidParams)
		return
	}

	problem, err := a.problemService.GetProblemByID(uint(problemID))
	if err != nil {
		if err.Error() == "题目不存在" {
			response.Error(c, response.CodeProblemNotFound)
			return
		}
		response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, problem)
}

// UpdateProblem 更新题目（管理员）
// @Summary 更新题目
// @Description 管理员更新题目信息
// @Tags 题目管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "题目ID"
// @Param request body service.UpdateProblemRequest true "更新信息"
// @Success 200 {object} response.Response{data=model.Problem} "更新成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "权限不足"
// @Failure 404 {object} response.Response "题目不存在"
// @Router /api/admin/problems/{id} [put]
func (a *ProblemAPI) UpdateProblem(c *gin.Context) {
	problemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.CodeInvalidParams)
		return
	}

	// 检查权限
	userID, _ := c.Get("user_id")
	isAdmin, _ := c.Get("is_admin")
	
	if !isAdmin.(bool) {
		// 获取题目信息以检查方向
		problem, err := a.problemService.GetProblemByID(uint(problemID))
		if err != nil {
			if err.Error() == "题目不存在" {
				response.Error(c, response.CodeProblemNotFound)
				return
			}
			response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
			return
		}

		isManager, err := a.directionService.CheckDirectionManager(problem.DirectionID, userID.(uint))
		if err != nil {
			response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
			return
		}
		if !isManager {
			response.Error(c, response.CodeForbidden)
			return
		}
	}

	var req service.UpdateProblemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBindError)
		return
	}

	problem, err := a.problemService.UpdateProblem(uint(problemID), &req)
	if err != nil {
		if err.Error() == "题目不存在" {
			response.Error(c, response.CodeProblemNotFound)
			return
		}
		response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, problem)
}

// DeleteProblem 删除题目（管理员）
// @Summary 删除题目
// @Description 管理员删除题目
// @Tags 题目管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "题目ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "权限不足"
// @Failure 404 {object} response.Response "题目不存在"
// @Router /api/admin/problems/{id} [delete]
func (a *ProblemAPI) DeleteProblem(c *gin.Context) {
	problemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.CodeInvalidParams)
		return
	}

	// 检查权限
	userID, _ := c.Get("user_id")
	isAdmin, _ := c.Get("is_admin")
	
	if !isAdmin.(bool) {
		// 获取题目信息以检查方向
		problem, err := a.problemService.GetProblemByID(uint(problemID))
		if err != nil {
			if err.Error() == "题目不存在" {
				response.Error(c, response.CodeProblemNotFound)
				return
			}
			response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
			return
		}

		isManager, err := a.directionService.CheckDirectionManager(problem.DirectionID, userID.(uint))
		if err != nil {
			response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
			return
		}
		if !isManager {
			response.Error(c, response.CodeForbidden)
			return
		}
	}

	if err := a.problemService.DeleteProblem(uint(problemID)); err != nil {
		if err.Error() == "题目不存在" {
			response.Error(c, response.CodeProblemNotFound)
			return
		}
		response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, nil)
}

// CreateSubmissionPoint 创建提交点（管理员）
// @Summary 创建提交点
// @Description 管理员为题目创建提交点
// @Tags 题目管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "题目ID"
// @Param request body service.CreateSubmissionPointRequest true "提交点信息"
// @Success 200 {object} response.Response{data=model.SubmissionPoint} "创建成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "权限不足"
// @Failure 404 {object} response.Response "题目不存在"
// @Router /api/admin/problems/{id}/submission-points [post]
func (a *ProblemAPI) CreateSubmissionPoint(c *gin.Context) {
	problemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.CodeInvalidParams)
		return
	}

	// 检查权限
	userID, _ := c.Get("user_id")
	isAdmin, _ := c.Get("is_admin")
	
	if !isAdmin.(bool) {
		// 获取题目信息以检查方向
		problem, err := a.problemService.GetProblemByID(uint(problemID))
		if err != nil {
			if err.Error() == "题目不存在" {
				response.Error(c, response.CodeProblemNotFound)
				return
			}
			response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
			return
		}

		isManager, err := a.directionService.CheckDirectionManager(problem.DirectionID, userID.(uint))
		if err != nil {
			response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
			return
		}
		if !isManager {
			response.Error(c, response.CodeForbidden)
			return
		}
	}

	var req service.CreateSubmissionPointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBindError)
		return
	}

	submissionPoint, err := a.problemService.CreateSubmissionPoint(uint(problemID), &req)
	if err != nil {
		if err.Error() == "题目不存在" {
			response.Error(c, response.CodeProblemNotFound)
			return
		}
		response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, submissionPoint)
}

// GetSubmissionPoints 获取提交点列表
// @Summary 获取提交点列表
// @Description 获取指定题目的提交点列表
// @Tags 题目管理
// @Accept json
// @Produce json
// @Param id path int true "题目ID"
// @Success 200 {object} response.Response{data=[]model.SubmissionPoint} "获取成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 500 {object} response.Response "内部错误"
// @Router /api/problems/{id}/submission-points [get]
func (a *ProblemAPI) GetSubmissionPoints(c *gin.Context) {
	problemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.CodeInvalidParams)
		return
	}

	submissionPoints, err := a.problemService.GetSubmissionPoints(uint(problemID))
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, submissionPoints)
}

// UpdateSubmissionPoint 更新提交点（管理员）
// @Summary 更新提交点
// @Description 管理员更新提交点信息
// @Tags 题目管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "提交点ID"
// @Param request body service.UpdateSubmissionPointRequest true "更新信息"
// @Success 200 {object} response.Response{data=model.SubmissionPoint} "更新成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "权限不足"
// @Failure 404 {object} response.Response "提交点不存在"
// @Router /api/admin/submission-points/{id} [put]
func (a *ProblemAPI) UpdateSubmissionPoint(c *gin.Context) {
	submissionPointID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.CodeInvalidParams)
		return
	}

	var req service.UpdateSubmissionPointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBindError)
		return
	}

	submissionPoint, err := a.problemService.UpdateSubmissionPoint(uint(submissionPointID), &req)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, submissionPoint)
}

// DeleteSubmissionPoint 删除提交点（管理员）
// @Summary 删除提交点
// @Description 管理员删除提交点
// @Tags 题目管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "提交点ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "权限不足"
// @Router /api/admin/submission-points/{id} [delete]
func (a *ProblemAPI) DeleteSubmissionPoint(c *gin.Context) {
	submissionPointID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.CodeInvalidParams)
		return
	}

	if err := a.problemService.DeleteSubmissionPoint(uint(submissionPointID)); err != nil {
		response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, nil)
}