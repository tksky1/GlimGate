package service

import (
	"errors"

	"github.com/tksky1/glimgate/internal/model"
	"github.com/tksky1/glimgate/pkg/database"
	"gorm.io/gorm"
)

// ProblemService 题目服务
type ProblemService struct{}

// CreateProblemRequest 创建题目请求结构
type CreateProblemRequest struct {
	Title       string `json:"title" binding:"required" example:"实现一个简单的计算器"`
	Description string `json:"description" binding:"required" example:"使用HTML、CSS、JavaScript实现一个基本的计算器功能"`
	DirectionID uint   `json:"direction_id" binding:"required" example:"1"`
}

// UpdateProblemRequest 更新题目请求结构
type UpdateProblemRequest struct {
	Title       string `json:"title" example:"实现一个简单的计算器"`
	Description string `json:"description" example:"使用HTML、CSS、JavaScript实现一个基本的计算器功能"`
}

// CreateSubmissionPointRequest 创建提交点请求结构
type CreateSubmissionPointRequest struct {
	Name     string `json:"name" binding:"required" example:"源代码提交"`
	MaxScore int    `json:"max_score" binding:"required,min=1" example:"100"`
}

// UpdateSubmissionPointRequest 更新提交点请求结构
type UpdateSubmissionPointRequest struct {
	Name     string `json:"name" example:"源代码提交"`
	MaxScore int    `json:"max_score" binding:"min=1" example:"100"`
}

// NewProblemService 创建题目服务实例
func NewProblemService() *ProblemService {
	return &ProblemService{}
}

// CreateProblem 创建题目
func (s *ProblemService) CreateProblem(req *CreateProblemRequest) (*model.Problem, error) {
	db := database.GetDB()

	// 检查方向是否存在
	var direction model.Direction
	if err := db.First(&direction, req.DirectionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("方向不存在")
		}
		return nil, err
	}

	// 创建题目
	problem := model.Problem{
		Title:       req.Title,
		Description: req.Description,
		DirectionID: req.DirectionID,
	}

	if err := db.Create(&problem).Error; err != nil {
		return nil, err
	}

	// 加载关联数据
	if err := db.Preload("Direction").First(&problem, problem.ID).Error; err != nil {
		return nil, err
	}

	return &problem, nil
}

// GetProblems 获取题目列表
func (s *ProblemService) GetProblems(directionID uint) ([]model.Problem, error) {
	db := database.GetDB()

	var problems []model.Problem
	query := db.Preload("Direction").Preload("SubmissionPoints")

	if directionID > 0 {
		query = query.Where("direction_id = ?", directionID)
	}

	if err := query.Find(&problems).Error; err != nil {
		return nil, err
	}

	return problems, nil
}

// GetProblemByID 根据ID获取题目
func (s *ProblemService) GetProblemByID(problemID uint) (*model.Problem, error) {
	db := database.GetDB()

	var problem model.Problem
	if err := db.Preload("Direction").Preload("SubmissionPoints").First(&problem, problemID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("题目不存在")
		}
		return nil, err
	}

	return &problem, nil
}

// UpdateProblem 更新题目
func (s *ProblemService) UpdateProblem(problemID uint, req *UpdateProblemRequest) (*model.Problem, error) {
	db := database.GetDB()

	var problem model.Problem
	if err := db.First(&problem, problemID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("题目不存在")
		}
		return nil, err
	}

	// 更新字段
	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}

	if len(updates) > 0 {
		if err := db.Model(&problem).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	// 重新加载包含关联数据的题目
	if err := db.Preload("Direction").Preload("SubmissionPoints").First(&problem, problem.ID).Error; err != nil {
		return nil, err
	}

	return &problem, nil
}

// DeleteProblem 删除题目
func (s *ProblemService) DeleteProblem(problemID uint) error {
	db := database.GetDB()

	var problem model.Problem
	if err := db.First(&problem, problemID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("题目不存在")
		}
		return err
	}

	// 检查是否有提交记录
	var submissionCount int64
	if err := db.Model(&model.Submission{}).Where("problem_id = ?", problemID).Count(&submissionCount).Error; err != nil {
		return err
	}
	if submissionCount > 0 {
		return errors.New("该题目已有提交记录，无法删除")
	}

	// 删除关联的提交点
	if err := db.Where("problem_id = ?", problemID).Delete(&model.SubmissionPoint{}).Error; err != nil {
		return err
	}

	return db.Delete(&problem).Error
}

// CreateSubmissionPoint 创建提交点
func (s *ProblemService) CreateSubmissionPoint(problemID uint, req *CreateSubmissionPointRequest) (*model.SubmissionPoint, error) {
	db := database.GetDB()

	// 检查题目是否存在
	var problem model.Problem
	if err := db.First(&problem, problemID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("题目不存在")
		}
		return nil, err
	}

	// 创建提交点
	submissionPoint := model.SubmissionPoint{
		Name:      req.Name,
		MaxScore:  req.MaxScore,
		ProblemID: problemID,
	}

	if err := db.Create(&submissionPoint).Error; err != nil {
		return nil, err
	}

	// 加载关联数据
	if err := db.Preload("Problem").First(&submissionPoint, submissionPoint.ID).Error; err != nil {
		return nil, err
	}

	return &submissionPoint, nil
}

// GetSubmissionPoints 获取提交点列表
func (s *ProblemService) GetSubmissionPoints(problemID uint) ([]model.SubmissionPoint, error) {
	db := database.GetDB()

	var submissionPoints []model.SubmissionPoint
	if err := db.Preload("Problem").Where("problem_id = ?", problemID).Find(&submissionPoints).Error; err != nil {
		return nil, err
	}

	return submissionPoints, nil
}

// UpdateSubmissionPoint 更新提交点
func (s *ProblemService) UpdateSubmissionPoint(submissionPointID uint, req *UpdateSubmissionPointRequest) (*model.SubmissionPoint, error) {
	db := database.GetDB()

	var submissionPoint model.SubmissionPoint
	if err := db.First(&submissionPoint, submissionPointID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("提交点不存在")
		}
		return nil, err
	}

	// 更新字段
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.MaxScore > 0 {
		updates["max_score"] = req.MaxScore
	}

	if len(updates) > 0 {
		if err := db.Model(&submissionPoint).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	// 重新加载包含关联数据的提交点
	if err := db.Preload("Problem").First(&submissionPoint, submissionPoint.ID).Error; err != nil {
		return nil, err
	}

	return &submissionPoint, nil
}

// DeleteSubmissionPoint 删除提交点
func (s *ProblemService) DeleteSubmissionPoint(submissionPointID uint) error {
	db := database.GetDB()

	var submissionPoint model.SubmissionPoint
	if err := db.First(&submissionPoint, submissionPointID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("提交点不存在")
		}
		return err
	}

	// 检查是否有提交记录
	var submissionCount int64
	if err := db.Model(&model.Submission{}).Where("submission_point_id = ?", submissionPointID).Count(&submissionCount).Error; err != nil {
		return err
	}
	if submissionCount > 0 {
		return errors.New("该提交点已有提交记录，无法删除")
	}

	return db.Delete(&submissionPoint).Error
}