package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tksky1/glimgate/internal/service"
	"github.com/tksky1/glimgate/pkg/response"
)

// SubmissionAPI 提交API处理器
type SubmissionAPI struct {
	submissionService *service.SubmissionService
}

// NewSubmissionAPI 创建提交API实例
func NewSubmissionAPI() *SubmissionAPI {
	return &SubmissionAPI{
		submissionService: service.NewSubmissionService(),
	}
}

// CreateSubmission 创建提交
// @Summary 创建提交
// @Description 用户提交作业
// @Tags 提交管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body service.CreateSubmissionRequest true "提交信息"
// @Success 200 {object} response.Response{data=model.Submission} "提交成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /api/submissions [post]
func (a *SubmissionAPI) CreateSubmission(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req service.CreateSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBindError)
		return
	}

	submission, err := a.submissionService.CreateSubmission(userID.(uint), &req)
	if err != nil {
		if err.Error() == "题目不存在" {
			response.Error(c, response.CodeProblemNotFound)
			return
		}
		if err.Error() == "提交点不存在或不属于该题目" {
			response.Error(c, response.CodeInvalidParams)
			return
		}
		response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, submission)
}

// GetMySubmissions 获取我的提交列表
// @Summary 获取我的提交列表
// @Description 获取当前用户的提交列表
// @Tags 提交管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param problem_id query int false "题目ID"
// @Success 200 {object} response.Response{data=[]service.SubmissionResponse} "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Router /api/submissions/my [get]
func (a *SubmissionAPI) GetMySubmissions(c *gin.Context) {
	userID, _ := c.Get("user_id")
	problemID, _ := strconv.ParseUint(c.DefaultQuery("problem_id", "0"), 10, 32)

	submissions, err := a.submissionService.GetUserSubmissions(userID.(uint), uint(problemID))
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, submissions)
}

// GetSubmission 获取提交详情
// @Summary 获取提交详情
// @Description 获取指定提交的详细信息
// @Tags 提交管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "提交ID"
// @Success 200 {object} response.Response{data=model.Submission} "获取成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "提交不存在"
// @Router /api/submissions/{id} [get]
func (a *SubmissionAPI) GetSubmission(c *gin.Context) {
	submissionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.CodeInvalidParams)
		return
	}

	submission, err := a.submissionService.GetSubmissionByID(uint(submissionID))
	if err != nil {
		if err.Error() == "提交不存在" {
			response.Error(c, response.CodeSubmissionNotFound)
			return
		}
		response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
		return
	}

	// 检查权限：只有提交者本人或管理员可以查看
	userID, _ := c.Get("user_id")
	isAdmin, _ := c.Get("is_admin")
	
	if !isAdmin.(bool) && submission.UserID != userID.(uint) {
		response.Error(c, response.CodeForbidden)
		return
	}

	response.Success(c, submission)
}

// GetSubmissionsForReview 获取待评分的提交列表（管理员）
// @Summary 获取待评分的提交列表
// @Description 管理员获取需要评分的提交列表
// @Tags 提交管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param problem_id query int false "题目ID"
// @Success 200 {object} response.Response{data=[]model.Submission} "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "权限不足"
// @Router /api/admin/submissions/review [get]
func (a *SubmissionAPI) GetSubmissionsForReview(c *gin.Context) {
	userID, _ := c.Get("user_id")
	problemID, _ := strconv.ParseUint(c.DefaultQuery("problem_id", "0"), 10, 32)

	submissions, err := a.submissionService.GetSubmissionsForReview(userID.(uint), uint(problemID))
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, submissions)
}

// DeleteSubmission 删除提交
// @Summary 删除提交
// @Description 用户删除自己的提交
// @Tags 提交管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "提交ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "权限不足"
// @Failure 404 {object} response.Response "提交不存在"
// @Router /api/submissions/{id} [delete]
func (a *SubmissionAPI) DeleteSubmission(c *gin.Context) {
	submissionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.CodeInvalidParams)
		return
	}

	userID, _ := c.Get("user_id")

	if err := a.submissionService.DeleteSubmission(uint(submissionID), userID.(uint)); err != nil {
		if err.Error() == "提交不存在或无权限删除" {
			response.Error(c, response.CodeSubmissionNotFound)
			return
		}
		response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, nil)
}