package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tksky1/glimgate/internal/api"
	"github.com/tksky1/glimgate/internal/middleware"
)

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine) {
	// 创建API实例
	userAPI := api.NewUserAPI()
	directionAPI := api.NewDirectionAPI()
	problemAPI := api.NewProblemAPI()
	submissionAPI := api.NewSubmissionAPI()
	scoreAPI := api.NewScoreAPI()

	// API路由组
	apiGroup := r.Group("/api")
	{
		// 认证相关路由（无需认证）
		authGroup := apiGroup.Group("/auth")
		{
			authGroup.POST("/register", userAPI.Register)
			authGroup.POST("/login", userAPI.Login)
		}

		// 公开路由（无需认证）
		apiGroup.GET("/directions", directionAPI.GetDirections)
		apiGroup.GET("/directions/:id", directionAPI.GetDirection)
		apiGroup.GET("/problems", problemAPI.GetProblems)
		apiGroup.GET("/problems/:id", problemAPI.GetProblem)
		apiGroup.GET("/problems/:id/submission-points", problemAPI.GetSubmissionPoints)
		apiGroup.GET("/ranking", scoreAPI.GetRanking)

		// 需要认证的路由
		authRequired := apiGroup.Group("")
		authRequired.Use(middleware.AuthMiddleware())
		{
			// 用户相关路由
			userGroup := authRequired.Group("/user")
			{
				userGroup.GET("/profile", userAPI.GetProfile)
			}

			// 提交相关路由
			submissionGroup := authRequired.Group("/submissions")
			{
				submissionGroup.POST("", submissionAPI.CreateSubmission)
				submissionGroup.GET("/my", submissionAPI.GetMySubmissions)
				submissionGroup.GET("/:id", submissionAPI.GetSubmission)
				submissionGroup.DELETE("/:id", submissionAPI.DeleteSubmission)
				submissionGroup.GET("/:id/scores", scoreAPI.GetScoresBySubmission)
			}

			// 评分相关路由
			scoreGroup := authRequired.Group("/scores")
			{
				scoreGroup.GET("/my", scoreAPI.GetMyScores)
			}

			// 用户评分查询路由
			authRequired.GET("/users/:id/scores", scoreAPI.GetScoresByUser)

			// 管理员路由
			adminGroup := authRequired.Group("/admin")
			adminGroup.Use(middleware.AdminMiddleware())
			{
				// 用户管理
				adminUserGroup := adminGroup.Group("/users")
				{
					adminUserGroup.GET("", userAPI.GetUsers)
					adminUserGroup.GET("/:id", userAPI.GetUser)
					adminUserGroup.PUT("/:id", userAPI.UpdateUser)
					adminUserGroup.DELETE("/:id", userAPI.DeleteUser)
				}

				// 方向管理
				adminDirectionGroup := adminGroup.Group("/directions")
				{
					adminDirectionGroup.POST("", directionAPI.CreateDirection)
					adminDirectionGroup.PUT("/:id", directionAPI.UpdateDirection)
					adminDirectionGroup.DELETE("/:id", directionAPI.DeleteDirection)
				}

				// 题目管理
				adminProblemGroup := adminGroup.Group("/problems")
				{
					adminProblemGroup.POST("", problemAPI.CreateProblem)
					adminProblemGroup.PUT("/:id", problemAPI.UpdateProblem)
					adminProblemGroup.DELETE("/:id", problemAPI.DeleteProblem)
					adminProblemGroup.POST("/:id/submission-points", problemAPI.CreateSubmissionPoint)
				}

				// 提交点管理
				adminGroup.PUT("/submission-points/:id", problemAPI.UpdateSubmissionPoint)
				adminGroup.DELETE("/submission-points/:id", problemAPI.DeleteSubmissionPoint)

				// 提交管理
				adminSubmissionGroup := adminGroup.Group("/submissions")
				{
					adminSubmissionGroup.GET("/review", submissionAPI.GetSubmissionsForReview)
				}

				// 评分管理
				adminScoreGroup := adminGroup.Group("/scores")
				{
					adminScoreGroup.POST("", scoreAPI.CreateScore)
					adminScoreGroup.GET("/my", scoreAPI.GetScoresByReviewer)
					adminScoreGroup.PUT("/:id", scoreAPI.UpdateScore)
					adminScoreGroup.DELETE("/:id", scoreAPI.DeleteScore)
				}
			}
		}
	}
}