package service

import (
	"errors"

	"github.com/tksky1/glimgate/internal/model"
	"github.com/tksky1/glimgate/pkg/database"
	"gorm.io/gorm"
)

// ScoreService 评分服务
type ScoreService struct{}

// CreateScoreRequest 创建评分请求结构
type CreateScoreRequest struct {
	Score        int    `json:"score" binding:"required,min=0" example:"85"`
	Comment      string `json:"comment" example:"代码实现良好，但缺少注释"`
	SubmissionID uint   `json:"submission_id" binding:"required" example:"1"`
}

// UpdateScoreRequest 更新评分请求结构
type UpdateScoreRequest struct {
	Score   int    `json:"score" binding:"required,min=0" example:"85"`
	Comment string `json:"comment" example:"代码实现良好，但缺少注释"`
}

// RankingItem 排行榜项目
type RankingItem struct {
	UserID   uint   `json:"user_id"`
	Nickname string `json:"nickname"`
	Score    int    `json:"score"`
}

// NewScoreService 创建评分服务实例
func NewScoreService() *ScoreService {
	return &ScoreService{}
}

// CreateScore 创建评分
func (s *ScoreService) CreateScore(reviewerID uint, req *CreateScoreRequest) (*model.Score, error) {
	db := database.GetDB()

	// 获取提交信息
	var submission model.Submission
	if err := db.Preload("Problem").Preload("SubmissionPoint").First(&submission, req.SubmissionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("提交不存在")
		}
		return nil, err
	}

	// 检查评分者是否为该方向的负责人
	directionService := NewDirectionService()
	isManager, err := directionService.CheckDirectionManager(submission.Problem.DirectionID, reviewerID)
	if err != nil {
		return nil, err
	}
	if !isManager {
		return nil, errors.New("无权限评分该提交")
	}

	// 检查分数是否超过最大分值
	if req.Score > submission.SubmissionPoint.MaxScore {
		return nil, errors.New("评分不能超过最大分值")
	}

	// 检查是否已有评分，如果有则更新，否则创建
	var score model.Score
	err = db.Where("submission_id = ? AND reviewer_id = ?", req.SubmissionID, reviewerID).First(&score).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 创建新评分
		score = model.Score{
			Score:        req.Score,
			Comment:      req.Comment,
			UserID:       submission.UserID,
			SubmissionID: req.SubmissionID,
			ReviewerID:   reviewerID,
		}
		if err := db.Create(&score).Error; err != nil {
			return nil, err
		}
	} else {
		// 更新现有评分
		updates := map[string]interface{}{
			"score":   req.Score,
			"comment": req.Comment,
		}
		if err := db.Model(&score).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	// 加载关联数据
	if err := db.Preload("User").Preload("Submission").Preload("Reviewer").First(&score, score.ID).Error; err != nil {
		return nil, err
	}

	return &score, nil
}

// GetScoresBySubmission 获取提交的评分列表
func (s *ScoreService) GetScoresBySubmission(submissionID uint) ([]model.Score, error) {
	db := database.GetDB()

	var scores []model.Score
	if err := db.Preload("User").Preload("Submission").Preload("Reviewer").Where("submission_id = ?", submissionID).Find(&scores).Error; err != nil {
		return nil, err
	}

	return scores, nil
}

// GetScoresByUser 获取用户的评分列表
func (s *ScoreService) GetScoresByUser(userID uint, problemID uint) ([]model.Score, error) {
	db := database.GetDB()

	query := db.Preload("User").Preload("Submission").Preload("Reviewer").Where("user_id = ?", userID)

	if problemID > 0 {
		// 需要通过submission表关联查询
		query = query.Joins("JOIN submissions ON scores.submission_id = submissions.id").
			Where("submissions.problem_id = ?", problemID)
	}

	var scores []model.Score
	if err := query.Find(&scores).Error; err != nil {
		return nil, err
	}

	return scores, nil
}

// GetScoresByReviewer 获取评分者的评分列表
func (s *ScoreService) GetScoresByReviewer(reviewerID uint, problemID uint) ([]model.Score, error) {
	db := database.GetDB()

	query := db.Preload("User").Preload("Submission").Preload("Reviewer").Where("reviewer_id = ?", reviewerID)

	if problemID > 0 {
		// 需要通过submission表关联查询
		query = query.Joins("JOIN submissions ON scores.submission_id = submissions.id").
			Where("submissions.problem_id = ?", problemID)
	}

	var scores []model.Score
	if err := query.Find(&scores).Error; err != nil {
		return nil, err
	}

	return scores, nil
}

// UpdateScore 更新评分
func (s *ScoreService) UpdateScore(scoreID uint, reviewerID uint, req *UpdateScoreRequest) (*model.Score, error) {
	db := database.GetDB()

	var score model.Score
	if err := db.Where("id = ? AND reviewer_id = ?", scoreID, reviewerID).First(&score).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("评分不存在或无权限修改")
		}
		return nil, err
	}

	// 获取提交点信息以检查最大分值
	var submission model.Submission
	if err := db.Preload("SubmissionPoint").First(&submission, score.SubmissionID).Error; err != nil {
		return nil, err
	}

	// 检查分数是否超过最大分值
	if req.Score > submission.SubmissionPoint.MaxScore {
		return nil, errors.New("评分不能超过最大分值")
	}

	// 更新评分
	updates := map[string]interface{}{
		"score":   req.Score,
		"comment": req.Comment,
	}
	if err := db.Model(&score).Updates(updates).Error; err != nil {
		return nil, err
	}

	// 重新加载包含关联数据的评分
	if err := db.Preload("User").Preload("Submission").Preload("Reviewer").First(&score, score.ID).Error; err != nil {
		return nil, err
	}

	return &score, nil
}

// DeleteScore 删除评分
func (s *ScoreService) DeleteScore(scoreID uint, reviewerID uint) error {
	db := database.GetDB()

	var score model.Score
	if err := db.Where("id = ? AND reviewer_id = ?", scoreID, reviewerID).First(&score).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("评分不存在或无权限删除")
		}
		return err
	}

	return db.Delete(&score).Error
}

// GetRanking 获取排行榜
func (s *ScoreService) GetRanking(directionID uint, limit int) ([]RankingItem, error) {
	db := database.GetDB()

	var rankings []RankingItem

	// 构建查询语句
	query := `
		SELECT 
			u.id as user_id,
			u.nickname,
			COALESCE(SUM(s.score), 0) as score
		FROM users u
		LEFT JOIN scores s ON u.id = s.user_id
		LEFT JOIN submissions sub ON s.submission_id = sub.id
		LEFT JOIN problems p ON sub.problem_id = p.id
	`

	args := []interface{}{}
	if directionID > 0 {
		query += " WHERE p.direction_id = ?"
		args = append(args, directionID)
	}

	query += `
		GROUP BY u.id, u.nickname
		ORDER BY score DESC
	`

	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)
	}

	if err := db.Raw(query, args...).Scan(&rankings).Error; err != nil {
		return nil, err
	}

	return rankings, nil
}