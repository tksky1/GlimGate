# GlimGate 工作室招新交题和评分系统

## 项目简介

GlimGate是一个基于Gin和GORM开发的工作室招新交题和评分系统后端，提供用户管理、题目提交、排行榜和管理员评分等功能。

## 功能特性

- **用户管理**: 用户注册、登录、权限控制
- **方向管理**: 多方向管理，支持设置方向负责人
- **题目管理**: 动态创建题目和提交点，支持多种提交方式
- **提交系统**: 用户可提交文本或Git仓库地址，支持重新提交
- **评分系统**: 管理员按题目评分，支持评分记录查询和修改
- **排行榜**: 各方向分数排名展示，仅显示昵称和分数
- **权限控制**: 完整的JWT认证和基于角色的访问控制

## 技术栈

- Go 1.21+
- Gin Web框架
- GORM ORM框架
- MySQL数据库
- JWT身份验证
- Swagger API文档
- bcrypt密码加密

## 系统架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   前端应用      │    │   API网关       │    │   后端服务      │
│                 │◄──►│                 │◄──►│                 │
│  React/Vue等    │    │  Gin Router     │    │  Business Logic │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │                       │
                                ▼                       ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │   中间件层      │    │   数据访问层    │
                       │                 │    │                 │
                       │ Auth/CORS/Log   │    │   GORM ORM      │
                       └─────────────────┘    └─────────────────┘
                                                       │
                                                       ▼
                                              ┌─────────────────┐
                                              │   MySQL数据库   │
                                              │                 │
                                              │   持久化存储    │
                                              └─────────────────┘
```

## 快速开始

### 环境要求

- Go 1.21+
- MySQL 5.7+
- Git

### 安装步骤

1. **克隆项目**
```bash
git clone https://github.com/tksky1/glimgate.git
cd glimgate
```

2. **安装依赖**
```bash
make deps
# 或者
go mod tidy
```

3. **配置数据库**

创建MySQL数据库：
```sql
CREATE DATABASE glimgate CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

修改 `config/config.yaml` 中的数据库配置：
```yaml
database:
  host: localhost
  port: 3306
  username: root
  password: your_password
  dbname: glimgate
  charset: utf8mb4
  parse_time: true
  loc: Local
```

4. **初始化数据库**
```bash
go run cmd/init/main.go
```

这将创建数据表并插入默认管理员账户：
- 用户名: `admin`
- 密码: `admin123`

5. **启动服务**
```bash
make run
# 或者
go run main.go
```

6. **访问服务**
- API服务: http://localhost:20401
- API文档: http://localhost:20401/swagger/index.html

## 配置说明

### 配置文件结构

```yaml
server:
  port: 20401              # 服务端口
  mode: debug             # 运行模式: debug, release, test

database:
  host: localhost         # 数据库主机
  port: 3306             # 数据库端口
  username: root         # 数据库用户名
  password: password     # 数据库密码
  dbname: glimgate       # 数据库名
  charset: utf8mb4       # 字符集
  parse_time: true       # 解析时间
  loc: Local             # 时区

jwt:
  secret: your-secret-key    # JWT密钥
  expire_hours: 24          # Token过期时间（小时）

cors:
  allow_origins: ["*"]      # 允许的源
  allow_methods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
  allow_headers: ["*"]      # 允许的请求头
```

## API接口

### 主要接口分类

1. **认证接口** (`/api/auth/`)
   - 用户注册
   - 用户登录

2. **用户接口** (`/api/user/`)
   - 获取用户信息
   - 用户管理（管理员）

3. **方向接口** (`/api/directions/`)
   - 方向列表查询
   - 方向管理（管理员）

4. **题目接口** (`/api/problems/`)
   - 题目列表查询
   - 题目管理（管理员）
   - 提交点管理

5. **提交接口** (`/api/submissions/`)
   - 创建提交
   - 查询提交记录
   - 提交管理

6. **评分接口** (`/api/scores/`)
   - 创建评分（管理员）
   - 查询评分记录
   - 评分管理

7. **排行榜接口** (`/api/ranking`)
   - 获取排行榜

详细的API文档请查看：[API文档](docs/API.md)

## 数据模型

### 核心实体关系

```
User (用户)
├── 1:N → Submission (提交)
├── 1:N → Score (评分记录)
└── N:M → Direction (方向负责人)

Direction (方向)
├── 1:N → Problem (题目)
└── N:M → User (负责人)

Problem (题目)
├── 1:N → SubmissionPoint (提交点)
└── 1:N → Submission (提交)

Submission (提交)
├── 1:N → Score (评分)
├── N:1 → User (提交者)
├── N:1 → Problem (题目)
└── N:1 → SubmissionPoint (提交点)

Score (评分)
├── N:1 → User (被评分用户)
├── N:1 → User (评分者)
└── N:1 → Submission (提交)
```

## 开发指南

### 项目结构

```
glimgate/
├── cmd/                    # 命令行工具
│   └── init/              # 数据库初始化
├── config/                # 配置文件
├── docs/                  # 文档
├── internal/              # 内部代码
│   ├── api/              # API处理器
│   ├── middleware/       # 中间件
│   ├── model/            # 数据模型
│   ├── router/           # 路由配置
│   └── service/          # 业务逻辑
├── pkg/                   # 公共包
│   ├── config/           # 配置管理
│   ├── database/         # 数据库连接
│   ├── jwt/              # JWT工具
│   ├── response/         # 响应格式
│   └── utils/            # 工具函数
├── .gitignore
├── go.mod
├── go.sum
├── main.go               # 主程序入口
├── Makefile              # 构建脚本
└── README.md
```

### 开发命令

```bash
# 安装依赖
make deps

# 运行项目
make run

# 构建项目
make build

# 运行测试
make test

# 生成API文档
make swagger

# 格式化代码
make fmt

# 代码检查
make lint

# 清理构建文件
make clean
```

### 添加新功能

1. **添加数据模型**: 在 `internal/model/` 中定义新的结构体
2. **添加服务层**: 在 `internal/service/` 中实现业务逻辑
3. **添加API处理器**: 在 `internal/api/` 中实现HTTP处理器
4. **配置路由**: 在 `internal/router/` 中添加路由规则
5. **更新文档**: 添加Swagger注释和API文档

## 部署

### Docker部署

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o glimgate main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/glimgate .
COPY --from=builder /app/config ./config
CMD ["./glimgate"]
```

### 生产环境配置

1. **修改配置文件**
```yaml
server:
  mode: release
  port: 20401

jwt:
  secret: "your-production-secret-key"
  expire_hours: 24

database:
  host: your-db-host
  username: your-db-user
  password: your-db-password
```

2. **设置环境变量**
```bash
export GIN_MODE=release
export DB_PASSWORD=your-secure-password
export JWT_SECRET=your-secure-jwt-secret
```

3. **使用反向代理**
```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    location / {
        proxy_pass http://localhost:20401;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 安全考虑

1. **密码安全**: 使用bcrypt加密存储密码
2. **JWT安全**: 使用强密钥，设置合理的过期时间
3. **权限控制**: 基于角色的访问控制，细粒度权限管理
4. **输入验证**: 所有用户输入都进行验证和过滤
5. **CORS配置**: 生产环境中限制允许的源
6. **HTTPS**: 生产环境中使用HTTPS加密传输

## 常见问题

### Q: 如何重置管理员密码？
A: 可以直接在数据库中修改，或者重新运行初始化脚本。

### Q: 如何添加新的方向负责人？
A: 使用管理员账户调用方向更新接口，在manager_ids中添加用户ID。

### Q: 提交后如何重新提交？
A: 用户可以使用相同的接口重新提交，系统会自动更新原有提交。

### Q: 评分后如何修改分数？
A: 评分者可以使用评分更新接口修改已提交的评分。

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 联系方式

- 项目地址: https://github.com/tksky1/glimgate
- 问题反馈: https://github.com/tksky1/glimgate/issues

## 更新日志

### v1.0.0 (2024-01-01)
- 初始版本发布
- 实现用户管理功能
- 实现题目提交系统
- 实现评分管理功能
- 实现排行榜功能
- 完整的API文档