package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tksky1/glimgate/internal/service"
	"github.com/tksky1/glimgate/pkg/response"
)

// ScoreAPI 评分API处理器
type ScoreAPI struct {
	scoreService *service.ScoreService
}

// NewScoreAPI 创建评分API实例
func NewScoreAPI() *ScoreAPI {
	return &ScoreAPI{
		scoreService: service.NewScoreService(),
	}
}

// CreateScore 创建评分（管理员）
// @Summary 创建评分
// @Description 管理员对提交进行评分
// @Tags 评分管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body service.CreateScoreRequest true "评分信息"
// @Success 200 {object} response.Response{data=model.Score} "评分成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "权限不足"
// @Router /api/admin/scores [post]
func (a *ScoreAPI) CreateScore(c *gin.Context) {
	reviewerID, _ := c.Get("user_id")

	var req service.CreateScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBindError)
		return
	}

	score, err := a.scoreService.CreateScore(reviewerID.(uint), &req)
	if err != nil {
		if err.Error() == "提交不存在" {
			response.Error(c, response.CodeSubmissionNotFound)
			return
		}
		if err.Error() == "无权限评分该提交" {
			response.Error(c, response.CodeForbidden)
			return
		}
		if err.Error() == "评分不能超过最大分值" {
			response.ErrorWithMsg(c, response.CodeInvalidParams, err.Error())
			return
		}
		response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, score)
}

// GetScoresBySubmission 获取提交的评分列表
// @Summary 获取提交的评分列表
// @Description 获取指定提交的所有评分
// @Tags 评分管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "提交ID"
// @Success 200 {object} response.Response{data=[]model.Score} "获取成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /api/submissions/{id}/scores [get]
func (a *ScoreAPI) GetScoresBySubmission(c *gin.Context) {
	submissionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.CodeInvalidParams)
		return
	}

	scores, err := a.scoreService.GetScoresBySubmission(uint(submissionID))
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, scores)
}

// GetScoresByUser 获取用户的评分列表
// @Summary 获取用户的评分列表
// @Description 获取指定用户的评分记录
// @Tags 评分管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Param problem_id query int false "题目ID"
// @Success 200 {object} response.Response{data=[]model.Score} "获取成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /api/users/{id}/scores [get]
func (a *ScoreAPI) GetScoresByUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.CodeInvalidParams)
		return
	}

	problemID, _ := strconv.ParseUint(c.DefaultQuery("problem_id", "0"), 10, 32)

	scores, err := a.scoreService.GetScoresByUser(uint(userID), uint(problemID))
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, scores)
}

// GetMyScores 获取我的评分列表
// @Summary 获取我的评分列表
// @Description 获取当前用户的评分记录
// @Tags 评分管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param problem_id query int false "题目ID"
// @Success 200 {object} response.Response{data=[]model.Score} "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Router /api/scores/my [get]
func (a *ScoreAPI) GetMyScores(c *gin.Context) {
	userID, _ := c.Get("user_id")
	problemID, _ := strconv.ParseUint(c.DefaultQuery("problem_id", "0"), 10, 32)

	scores, err := a.scoreService.GetScoresByUser(userID.(uint), uint(problemID))
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, scores)
}

// GetScoresByReviewer 获取评分者的评分列表（管理员）
// @Summary 获取评分者的评分列表
// @Description 管理员获取自己的评分记录
// @Tags 评分管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param problem_id query int false "题目ID"
// @Success 200 {object} response.Response{data=[]model.Score} "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "权限不足"
// @Router /api/admin/scores/my [get]
func (a *ScoreAPI) GetScoresByReviewer(c *gin.Context) {
	reviewerID, _ := c.Get("user_id")
	problemID, _ := strconv.ParseUint(c.DefaultQuery("problem_id", "0"), 10, 32)

	scores, err := a.scoreService.GetScoresByReviewer(reviewerID.(uint), uint(problemID))
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, scores)
}

// UpdateScore 更新评分（管理员）
// @Summary 更新评分
// @Description 管理员更新评分信息
// @Tags 评分管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "评分ID"
// @Param request body service.UpdateScoreRequest true "更新信息"
// @Success 200 {object} response.Response{data=model.Score} "更新成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "权限不足"
// @Failure 404 {object} response.Response "评分不存在"
// @Router /api/admin/scores/{id} [put]
func (a *ScoreAPI) UpdateScore(c *gin.Context) {
	scoreID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.CodeInvalidParams)
		return
	}

	reviewerID, _ := c.Get("user_id")

	var req service.UpdateScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBindError)
		return
	}

	score, err := a.scoreService.UpdateScore(uint(scoreID), reviewerID.(uint), &req)
	if err != nil {
		if err.Error() == "评分不存在或无权限修改" {
			response.Error(c, response.CodeForbidden)
			return
		}
		if err.Error() == "评分不能超过最大分值" {
			response.ErrorWithMsg(c, response.CodeInvalidParams, err.Error())
			return
		}
		response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, score)
}

// DeleteScore 删除评分（管理员）
// @Summary 删除评分
// @Description 管理员删除评分
// @Tags 评分管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "评分ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "权限不足"
// @Failure 404 {object} response.Response "评分不存在"
// @Router /api/admin/scores/{id} [delete]
func (a *ScoreAPI) DeleteScore(c *gin.Context) {
	scoreID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.CodeInvalidParams)
		return
	}

	reviewerID, _ := c.Get("user_id")

	if err := a.scoreService.DeleteScore(uint(scoreID), reviewerID.(uint)); err != nil {
		if err.Error() == "评分不存在或无权限删除" {
			response.Error(c, response.CodeForbidden)
			return
		}
		response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetRanking 获取排行榜
// @Summary 获取排行榜
// @Description 获取指定方向的排行榜
// @Tags 评分管理
// @Accept json
// @Produce json
// @Param direction_id query int false "方向ID"
// @Param limit query int false "限制数量" default(10)
// @Success 200 {object} response.Response{data=[]service.RankingItem} "获取成功"
// @Failure 500 {object} response.Response "内部错误"
// @Router /api/ranking [get]
func (a *ScoreAPI) GetRanking(c *gin.Context) {
	directionID, _ := strconv.ParseUint(c.DefaultQuery("direction_id", "0"), 10, 32)
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if limit <= 0 || limit > 100 {
		limit = 10
	}

	rankings, err := a.scoreService.GetRanking(uint(directionID), limit)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, rankings)
}