package service

import (
	"errors"

	"github.com/tksky1/glimgate/internal/model"
	"github.com/tksky1/glimgate/pkg/database"
	"gorm.io/gorm"
)

// DirectionService 方向服务
type DirectionService struct{}

// CreateDirectionRequest 创建方向请求结构
type CreateDirectionRequest struct {
	Name        string `json:"name" binding:"required" example:"前端开发"`
	Description string `json:"description" example:"负责前端页面开发和用户交互"`
	ManagerIDs  []uint `json:"manager_ids" example:"[1,2]"`
}

// UpdateDirectionRequest 更新方向请求结构
type UpdateDirectionRequest struct {
	Name        string `json:"name" example:"前端开发"`
	Description string `json:"description" example:"负责前端页面开发和用户交互"`
	ManagerIDs  []uint `json:"manager_ids" example:"[1,2]"`
}

// NewDirectionService 创建方向服务实例
func NewDirectionService() *DirectionService {
	return &DirectionService{}
}

// CreateDirection 创建方向
func (s *DirectionService) CreateDirection(req *CreateDirectionRequest) (*model.Direction, error) {
	db := database.GetDB()

	// 创建方向
	direction := model.Direction{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := db.Create(&direction).Error; err != nil {
		return nil, err
	}

	// 设置负责人
	if len(req.ManagerIDs) > 0 {
		var managers []model.User
		if err := db.Where("id IN ?", req.ManagerIDs).Find(&managers).Error; err != nil {
			return nil, err
		}
		if err := db.Model(&direction).Association("Managers").Replace(managers); err != nil {
			return nil, err
		}
	}

	// 重新加载包含关联数据的方向
	if err := db.Preload("Managers").First(&direction, direction.ID).Error; err != nil {
		return nil, err
	}

	return &direction, nil
}

// GetDirections 获取方向列表
func (s *DirectionService) GetDirections() ([]model.Direction, error) {
	db := database.GetDB()

	var directions []model.Direction
	if err := db.Preload("Managers").Find(&directions).Error; err != nil {
		return nil, err
	}

	return directions, nil
}

// GetDirectionByID 根据ID获取方向
func (s *DirectionService) GetDirectionByID(directionID uint) (*model.Direction, error) {
	db := database.GetDB()

	var direction model.Direction
	if err := db.Preload("Managers").Preload("Problems").First(&direction, directionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("方向不存在")
		}
		return nil, err
	}

	return &direction, nil
}

// UpdateDirection 更新方向
func (s *DirectionService) UpdateDirection(directionID uint, req *UpdateDirectionRequest) (*model.Direction, error) {
	db := database.GetDB()

	var direction model.Direction
	if err := db.First(&direction, directionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("方向不存在")
		}
		return nil, err
	}

	// 更新基本信息
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}

	if len(updates) > 0 {
		if err := db.Model(&direction).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	// 更新负责人
	if req.ManagerIDs != nil {
		var managers []model.User
		if len(req.ManagerIDs) > 0 {
			if err := db.Where("id IN ?", req.ManagerIDs).Find(&managers).Error; err != nil {
				return nil, err
			}
		}
		if err := db.Model(&direction).Association("Managers").Replace(managers); err != nil {
			return nil, err
		}
	}

	// 重新加载包含关联数据的方向
	if err := db.Preload("Managers").First(&direction, direction.ID).Error; err != nil {
		return nil, err
	}

	return &direction, nil
}

// DeleteDirection 删除方向
func (s *DirectionService) DeleteDirection(directionID uint) error {
	db := database.GetDB()

	var direction model.Direction
	if err := db.First(&direction, directionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("方向不存在")
		}
		return err
	}

	// 检查是否有关联的题目
	var problemCount int64
	if err := db.Model(&model.Problem{}).Where("direction_id = ?", directionID).Count(&problemCount).Error; err != nil {
		return err
	}
	if problemCount > 0 {
		return errors.New("该方向下还有题目，无法删除")
	}

	// 清除关联的负责人
	if err := db.Model(&direction).Association("Managers").Clear(); err != nil {
		return err
	}

	return db.Delete(&direction).Error
}

// CheckDirectionManager 检查用户是否为方向负责人
func (s *DirectionService) CheckDirectionManager(directionID, userID uint) (bool, error) {
	db := database.GetDB()

	var count int64
	if err := db.Table("direction_managers").
		Where("direction_id = ? AND user_id = ?", directionID, userID).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}