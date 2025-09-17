package service

import (
	"errors"

	"github.com/tksky1/glimgate/internal/model"
	"github.com/tksky1/glimgate/pkg/database"
	"github.com/tksky1/glimgate/pkg/jwt"
	"github.com/tksky1/glimgate/pkg/utils"
	"gorm.io/gorm"
)

// UserService 用户服务
type UserService struct{}

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	Username  string `json:"username" binding:"required" example:"user123"`
	Password  string `json:"password" binding:"required,min=6" example:"password123"`
	Nickname  string `json:"nickname" binding:"required" example:"小明"`
	RealName  string `json:"real_name" binding:"required" example:"张三"`
	College   string `json:"college" binding:"required" example:"计算机学院"`
	StudentID string `json:"student_id" binding:"required" example:"2021001001"`
	QQ        string `json:"qq" example:"123456789"`
	Email     string `json:"email" example:"user@example.com"`
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"user123"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	Token string     `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  model.User `json:"user"`
}

// UpdateUserRequest 更新用户请求结构
type UpdateUserRequest struct {
	Nickname  string `json:"nickname" example:"小明"`
	RealName  string `json:"real_name" example:"张三"`
	College   string `json:"college" example:"计算机学院"`
	StudentID string `json:"student_id" example:"2021001001"`
	QQ        string `json:"qq" example:"123456789"`
	Email     string `json:"email" example:"user@example.com"`
	IsAdmin   *bool  `json:"is_admin" example:"false"`
}

// NewUserService 创建用户服务实例
func NewUserService() *UserService {
	return &UserService{}
}

// Register 用户注册
func (s *UserService) Register(req *RegisterRequest) (*model.User, error) {
	db := database.GetDB()

	// 检查用户名是否已存在
	var existingUser model.User
	if err := db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := model.User{
		Username:  req.Username,
		Password:  hashedPassword,
		Nickname:  req.Nickname,
		RealName:  req.RealName,
		College:   req.College,
		StudentID: req.StudentID,
		QQ:        req.QQ,
		Email:     req.Email,
		IsAdmin:   false,
	}

	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// Login 用户登录
func (s *UserService) Login(req *LoginRequest) (*LoginResponse, error) {
	db := database.GetDB()

	// 查找用户
	var user model.User
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("密码错误")
	}

	// 生成token
	token, err := jwt.GenerateToken(user.ID, user.Username, user.IsAdmin)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(userID uint) (*model.User, error) {
	db := database.GetDB()

	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	return &user, nil
}

// GetUsers 获取用户列表
func (s *UserService) GetUsers(page, pageSize int) ([]model.User, int64, error) {
	db := database.GetDB()

	var users []model.User
	var total int64

	// 获取总数
	if err := db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(userID uint, req *UpdateUserRequest) (*model.User, error) {
	db := database.GetDB()

	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	// 更新字段
	updates := make(map[string]interface{})
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.RealName != "" {
		updates["real_name"] = req.RealName
	}
	if req.College != "" {
		updates["college"] = req.College
	}
	if req.StudentID != "" {
		updates["student_id"] = req.StudentID
	}
	if req.QQ != "" {
		updates["qq"] = req.QQ
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.IsAdmin != nil {
		updates["is_admin"] = *req.IsAdmin
	}

	if err := db.Model(&user).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(userID uint) error {
	db := database.GetDB()

	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

	return db.Delete(&user).Error
}