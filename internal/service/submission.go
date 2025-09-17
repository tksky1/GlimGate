package service

import (
	"errors"

	"github.com/tksky1/glimgate/internal/model"
	"github.com/tksky1/glimgate/pkg/database"
	"gorm.io/gorm"
)

// SubmissionService 提交服务
type SubmissionService struct{}

// CreateSubmissionRequest 创建提交请求结构
type CreateSubmissionRequest struct {
	Content           string `json:"content" binding:"required" example:"https://github.com/user/project"`
	ProblemID         uint   `json:"problem_id" binding:"required" example:"1"`
	SubmissionPointID uint   `json:"submission_point_id" binding:"required" example:"1"`
}

// SubmissionResponse 提交响应结构
type SubmissionResponse struct {
	model.Submission
	TotalScore int `json:"total_score"`
}

// NewSubmissionService 创建提交服务实例
func NewSubmissionService() *SubmissionService {
	return &SubmissionService{}
}

// CreateSubmission 创建提交
func (s *SubmissionService) CreateSubmission(userID uint, req *CreateSubmissionRequest) (*model.Submission, error) {
	db := database.GetDB()

	// 检查题目是否存在
	var problem model.Problem
	if err := db.First(&problem, req.ProblemID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("题目不存在")
		}
		return nil, err
	}

	// 检查提交点是否存在且属于该题目
	var submissionPoint model.SubmissionPoint
	if err := db.Where("id = ? AND problem_id = ?", req.SubmissionPointID, req.ProblemID).First(&submissionPoint).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("提交点不存在或不属于该题目")
		}
		return nil, err
	}

	// 检查是否已有提交，如果有则更新，否则创建
	var submission model.Submission
	err := db.Where("user_id = ? AND problem_id = ? AND submission_point_id = ?", 
		userID, req.ProblemID, req.SubmissionPointID).First(&submission).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 创建新提交
		submission = model.Submission{
			Content:           req.Content,
			UserID:            userID,
			ProblemID:         req.ProblemID,
			SubmissionPointID: req.SubmissionPointID,
		}
		if err := db.Create(&submission).Error; err != nil {
			return nil, err
		}
	} else {
		// 更新现有提交
		if err := db.Model(&submission).Update("content", req.Content).Error; err != nil {
			return nil, err
		}
	}

	// 加载关联数据
	if err := db.Preload("User").Preload("Problem").Preload("SubmissionPoint").First(&submission, submission.ID).Error; err != nil {
		return nil, err
	}

	return &submission, nil
}

// GetUserSubmissions 获取用户提交列表
func (s *SubmissionService) GetUserSubmissions(userID uint, problemID uint) ([]SubmissionResponse, error) {
	db := database.GetDB()

	var submissions []model.Submission
	query := db.Preload("User").Preload("Problem").Preload("SubmissionPoint").Preload("Scores").Where("user_id = ?", userID)

	if problemID > 0 {
		query = query.Where("problem_id = ?", problemID)
	}

	if err := query.Find(&submissions).Error; err != nil {
		return nil, err
	}

	// 计算总分
	var result []SubmissionResponse
	for _, submission := range submissions {
		totalScore := 0
		for _, score := range submission.Scores {
			totalScore += score.Score
		}
		result = append(result, SubmissionResponse{
			Submission: submission,
			TotalScore: totalScore,
		})
	}

	return result, nil
}

// GetSubmissionByID 根据ID获取提交
func (s *SubmissionService) GetSubmissionByID(submissionID uint) (*model.Submission, error) {
	db := database.GetDB()

	var submission model.Submission
	if err := db.Preload("User").Preload("Problem").Preload("SubmissionPoint").Preload("Scores").First(&submission, submissionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("提交不存在")
		}
		return nil, err
	}

	return &submission, nil
}

// GetSubmissionsForReview 获取待评分的提交列表（管理员用）
func (s *SubmissionService) GetSubmissionsForReview(reviewerID uint, problemID uint) ([]model.Submission, error) {
	db := database.GetDB()

	// 首先获取该管理员负责的方向
	var directionIDs []uint
	if err := db.Table("direction_managers").Where("user_id = ?", reviewerID).Pluck("direction_id", &directionIDs).Error; err != nil {
		return nil, err
	}

	if len(directionIDs) == 0 {
		return []model.Submission{}, nil
	}

	// 构建查询
	query := db.Preload("User").Preload("Problem").Preload("SubmissionPoint").Preload("Scores")

	if problemID > 0 {
		// 检查题目是否属于该管理员负责的方向
		var problem model.Problem
		if err := db.First(&problem, problemID).Error; err != nil {
			return nil, err
		}
		
		found := false
		for _, dirID := range directionIDs {
			if problem.DirectionID == dirID {
				found = true
				break
			}
		}
		if !found {
			return []model.Submission{}, nil
		}
		
		query = query.Where("problem_id = ?", problemID)
	} else {
		// 获取所有负责方向下的题目
		var problemIDs []uint
		if err := db.Model(&model.Problem{}).Where("direction_id IN ?", directionIDs).Pluck("id", &problemIDs).Error; err != nil {
			return nil, err
		}
		if len(problemIDs) == 0 {
			return []model.Submission{}, nil
		}
		query = query.Where("problem_id IN ?", problemIDs)
	}

	var submissions []model.Submission
	if err := query.Find(&submissions).Error; err != nil {
		return nil, err
	}

	return submissions, nil
}

// DeleteSubmission 删除提交
func (s *SubmissionService) DeleteSubmission(submissionID uint, userID uint) error {
	db := database.GetDB()

	var submission model.Submission
	if err := db.Where("id = ? AND user_id = ?", submissionID, userID).First(&submission).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("提交不存在或无权限删除")
		}
		return err
	}

	// 删除相关评分
	if err := db.Where("submission_id = ?", submissionID).Delete(&model.Score{}).Error; err != nil {
		return err
	}

	return db.Delete(&submission).Error
}