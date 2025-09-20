# GlimGate API 接口文档

## 概述

GlimGate是一个工作室招新交题和评分系统，提供完整的用户管理、题目提交、评分管理和排行榜功能。

## 基础信息

- **Base URL**: `http://localhost:20401`
- **API版本**: v1
- **认证方式**: Bearer Token (JWT)

## 响应格式

所有API接口都返回统一的JSON格式：

```json
{
  "code": 0,
  "msg": "success",
  "data": {}
}
```

### 响应码说明

- `0`: 成功
- `1001`: 用户不存在
- `1002`: 用户已存在
- `1003`: 密码错误
- `1004`: 未授权
- `1005`: 权限不足
- `1006`: 无效的token
- `2001`: 方向不存在
- `2002`: 题目不存在
- `2003`: 提交不存在
- `3001`: 参数错误
- `3002`: 参数绑定失败
- `5001`: 数据库错误
- `5002`: 内部错误

## 认证

除了注册、登录和部分公开接口外，所有接口都需要在请求头中携带JWT token：

```
Authorization: Bearer <your_jwt_token>
```

## 接口分类

### 1. 用户管理

#### 用户注册
- **POST** `/api/auth/register`
- **描述**: 用户注册
- **请求体**:
```json
{
  "username": "user123",
  "password": "password123",
  "nickname": "小明",
  "real_name": "张三",
  "college": "计算机学院",
  "student_id": "2021001001",
  "qq": "123456789",
  "email": "user@example.com"
}
```

#### 用户登录
- **POST** `/api/auth/login`
- **描述**: 用户登录
- **请求体**:
```json
{
  "username": "user123",
  "password": "password123"
}
```

#### 获取用户信息
- **GET** `/api/user/profile`
- **描述**: 获取当前登录用户信息
- **需要认证**: 是

### 2. 方向管理

#### 获取方向列表
- **GET** `/api/directions`
- **描述**: 获取所有方向列表
- **需要认证**: 否

#### 获取方向详情
- **GET** `/api/directions/{id}`
- **描述**: 获取指定方向的详细信息
- **需要认证**: 否

#### 创建方向（管理员）
- **POST** `/api/admin/directions`
- **描述**: 创建新方向
- **需要认证**: 是（管理员）
- **请求体**:
```json
{
  "name": "前端开发",
  "description": "负责前端页面开发和用户交互",
  "manager_ids": [1, 2]
}
```

### 3. 题目管理

#### 获取题目列表
- **GET** `/api/problems?direction_id=1`
- **描述**: 获取题目列表，可按方向筛选
- **需要认证**: 否

#### 获取题目详情
- **GET** `/api/problems/{id}`
- **描述**: 获取指定题目的详细信息
- **需要认证**: 否

#### 创建题目（管理员）
- **POST** `/api/admin/problems`
- **描述**: 创建新题目
- **需要认证**: 是（管理员或方向负责人）
- **请求体**:
```json
{
  "title": "实现一个简单的计算器",
  "description": "使用HTML、CSS、JavaScript实现一个基本的计算器功能",
  "direction_id": 1
}
```

#### 创建提交点（管理员）
- **POST** `/api/admin/problems/{id}/submission-points`
- **描述**: 为题目创建提交点
- **需要认证**: 是（管理员或方向负责人）
- **请求体**:
```json
{
  "name": "源代码提交",
  "max_score": 100
}
```

### 4. 提交管理

#### 创建提交
- **POST** `/api/submissions`
- **描述**: 用户提交作业
- **需要认证**: 是
- **请求体**:
```json
{
  "content": "https://github.com/user/project",
  "problem_id": 1,
  "submission_point_id": 1
}
```

#### 获取我的提交列表
- **GET** `/api/submissions/my?problem_id=1`
- **描述**: 获取当前用户的提交列表
- **需要认证**: 是

#### 获取提交详情
- **GET** `/api/submissions/{id}`
- **描述**: 获取指定提交的详细信息
- **需要认证**: 是

### 5. 评分管理

#### 创建评分（管理员）
- **POST** `/api/admin/scores`
- **描述**: 管理员对提交进行评分
- **需要认证**: 是（管理员或方向负责人）
- **请求体**:
```json
{
  "score": 85,
  "comment": "代码实现良好，但缺少注释",
  "submission_id": 1
}
```

#### 获取我的评分列表
- **GET** `/api/scores/my?problem_id=1`
- **描述**: 获取当前用户的评分记录
- **需要认证**: 是

#### 获取待评分提交列表（管理员）
- **GET** `/api/admin/submissions/review?problem_id=1`
- **描述**: 管理员获取需要评分的提交列表
- **需要认证**: 是（管理员或方向负责人）

### 6. 排行榜

#### 获取排行榜
- **GET** `/api/ranking?direction_id=1&limit=10`
- **描述**: 获取指定方向的排行榜
- **需要认证**: 否

## 数据模型

### 用户 (User)
```json
{
  "id": 1,
  "username": "user123",
  "nickname": "小明",
  "real_name": "张三",
  "college": "计算机学院",
  "student_id": "2021001001",
  "qq": "123456789",
  "email": "user@example.com",
  "is_admin": false,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### 方向 (Direction)
```json
{
  "id": 1,
  "name": "前端开发",
  "description": "负责前端页面开发和用户交互",
  "managers": [
    {
      "id": 1,
      "nickname": "管理员1"
    }
  ],
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### 题目 (Problem)
```json
{
  "id": 1,
  "title": "实现一个简单的计算器",
  "description": "使用HTML、CSS、JavaScript实现一个基本的计算器功能",
  "direction_id": 1,
  "direction": {
    "id": 1,
    "name": "前端开发"
  },
  "submission_points": [
    {
      "id": 1,
      "name": "源代码提交",
      "max_score": 100
    }
  ],
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### 提交 (Submission)
```json
{
  "id": 1,
  "content": "https://github.com/user/project",
  "user_id": 1,
  "problem_id": 1,
  "submission_point_id": 1,
  "user": {
    "id": 1,
    "nickname": "小明"
  },
  "problem": {
    "id": 1,
    "title": "实现一个简单的计算器"
  },
  "submission_point": {
    "id": 1,
    "name": "源代码提交",
    "max_score": 100
  },
  "scores": [
    {
      "id": 1,
      "score": 85,
      "comment": "代码实现良好"
    }
  ],
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### 评分 (Score)
```json
{
  "id": 1,
  "score": 85,
  "comment": "代码实现良好，但缺少注释",
  "user_id": 1,
  "submission_id": 1,
  "reviewer_id": 2,
  "user": {
    "id": 1,
    "nickname": "小明"
  },
  "reviewer": {
    "id": 2,
    "nickname": "评分老师"
  },
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

## 使用示例

### 1. 用户注册和登录流程

```bash
# 1. 用户注册
curl -X POST http://localhost:20401/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123",
    "nickname": "测试用户",
    "real_name": "张三",
    "college": "计算机学院",
    "student_id": "2021001001"
  }'

# 2. 用户登录
curl -X POST http://localhost:20401/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'

# 3. 使用返回的token访问需要认证的接口
curl -X GET http://localhost:20401/api/user/profile \
  -H "Authorization: Bearer <your_jwt_token>"
```

### 2. 提交作业流程

```bash
# 1. 获取题目列表
curl -X GET http://localhost:20401/api/problems

# 2. 获取题目详情和提交点
curl -X GET http://localhost:20401/api/problems/1

# 3. 提交作业
curl -X POST http://localhost:20401/api/submissions \
  -H "Authorization: Bearer <your_jwt_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "https://github.com/user/project",
    "problem_id": 1,
    "submission_point_id": 1
  }'

# 4. 查看我的提交
curl -X GET http://localhost:20401/api/submissions/my \
  -H "Authorization: Bearer <your_jwt_token>"
```

### 3. 管理员评分流程

```bash
# 1. 获取待评分的提交
curl -X GET http://localhost:20401/api/admin/submissions/review \
  -H "Authorization: Bearer <admin_jwt_token>"

# 2. 对提交进行评分
curl -X POST http://localhost:20401/api/admin/scores \
  -H "Authorization: Bearer <admin_jwt_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "score": 85,
    "comment": "代码实现良好，但缺少注释",
    "submission_id": 1
  }'

# 3. 查看评分记录
curl -X GET http://localhost:20401/api/admin/scores/my \
  -H "Authorization: Bearer <admin_jwt_token>"
```

## 部署说明

1. 配置数据库连接信息在 `config/config.yaml`
2. 运行 `make deps` 安装依赖
3. 运行 `make run` 启动服务
4. 访问 `http://localhost:20401/swagger/index.html` 查看完整的API文档

## 注意事项

1. 所有时间字段都使用UTC时间格式
2. 密码会自动加密存储，不会在API响应中返回
3. JWT token默认有效期为24小时
4. 管理员权限由 `is_admin` 字段控制
5. 方向负责人可以管理自己负责方向下的题目和评分
6. 用户只能查看和修改自己的提交
7. 评分不能超过提交点设置的最大分值