package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Username string `json:"username" gorm:"uniqueIndex;size:50;not null" binding:"required" example:"user123"`
	Password string `json:"-" gorm:"size:255;not null"`
	Nickname string `json:"nickname" gorm:"size:50;not null" binding:"required" example:"小明"`
	RealName string `json:"real_name" gorm:"size:50;not null" binding:"required" example:"张三"`
	College  string `json:"college" gorm:"size:100;not null" binding:"required" example:"计算机学院"`
	StudentID string `json:"student_id" gorm:"size:20;not null" binding:"required" example:"2021001001"`
	QQ       string `json:"qq" gorm:"size:20" example:"123456789"`
	Email    string `json:"email" gorm:"size:100" example:"user@example.com"`
	IsAdmin  bool   `json:"is_admin" gorm:"default:false"`
}

// Direction 方向模型
type Direction struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Name        string `json:"name" gorm:"size:100;not null" binding:"required" example:"前端开发"`
	Description string `json:"description" gorm:"type:text" example:"负责前端页面开发和用户交互"`

	// 关联关系
	Managers []User    `json:"managers" gorm:"many2many:direction_managers;"`
	Problems []Problem `json:"problems,omitempty"`
}

// Problem 题目模型
type Problem struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Title       string `json:"title" gorm:"size:200;not null" binding:"required" example:"实现一个简单的计算器"`
	Description string `json:"description" gorm:"type:text" binding:"required" example:"使用HTML、CSS、JavaScript实现一个基本的计算器功能"`
	DirectionID uint   `json:"direction_id" binding:"required" example:"1"`

	// 关联关系
	Direction        Direction        `json:"direction,omitempty"`
	SubmissionPoints []SubmissionPoint `json:"submission_points,omitempty"`
	Submissions      []Submission     `json:"submissions,omitempty"`
}

// SubmissionPoint 提交点模型
type SubmissionPoint struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Name      string `json:"name" gorm:"size:100;not null" binding:"required" example:"源代码提交"`
	MaxScore  int    `json:"max_score" gorm:"not null" binding:"required,min=1" example:"100"`
	ProblemID uint   `json:"problem_id" binding:"required" example:"1"`

	// 关联关系
	Problem     Problem      `json:"problem,omitempty"`
	Submissions []Submission `json:"submissions,omitempty"`
}

// Submission 提交模型
type Submission struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Content           string `json:"content" gorm:"type:text" binding:"required" example:"https://github.com/user/project"`
	UserID            uint   `json:"user_id" binding:"required" example:"1"`
	ProblemID         uint   `json:"problem_id" binding:"required" example:"1"`
	SubmissionPointID uint   `json:"submission_point_id" binding:"required" example:"1"`

	// 关联关系
	User            User            `json:"user,omitempty"`
	Problem         Problem         `json:"problem,omitempty"`
	SubmissionPoint SubmissionPoint `json:"submission_point,omitempty"`
	Scores          []Score         `json:"scores,omitempty"`
}

// Score 评分模型
type Score struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Score        int    `json:"score" gorm:"not null" binding:"required,min=0" example:"85"`
	Comment      string `json:"comment" gorm:"type:text" example:"代码实现良好，但缺少注释"`
	UserID       uint   `json:"user_id" binding:"required" example:"1"`
	SubmissionID uint   `json:"submission_id" binding:"required" example:"1"`
	ReviewerID   uint   `json:"reviewer_id" binding:"required" example:"2"`

	// 关联关系
	User       User       `json:"user,omitempty"`
	Submission Submission `json:"submission,omitempty"`
	Reviewer   User       `json:"reviewer,omitempty" gorm:"foreignKey:ReviewerID"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

func (Direction) TableName() string {
	return "directions"
}

func (Problem) TableName() string {
	return "problems"
}

func (SubmissionPoint) TableName() string {
	return "submission_points"
}

func (Submission) TableName() string {
	return "submissions"
}

func (Score) TableName() string {
	return "scores"
}