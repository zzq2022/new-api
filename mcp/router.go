package mcp

import (
	"github.com/QuantumNous/new-api/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 统一挂载所有 MCP 相关路由（包括反向代理入口与管理 API）
func RegisterRoutes(router *gin.Engine) {
	// 1. 初始化数据库表结构
	_ = InitDatabase()

	// 2. 外部代理通道：/mcp/:name/*any (支持 GET, POST, SSE 长连接)
	proxyGroup := router.Group("/mcp")
	proxyGroup.Use(middleware.CORS())
	proxyGroup.Use(middleware.TokenAuth(), MCPPermissionCheck())
	{
		proxyGroup.Any("/:name/*any", RelayMCP)
		proxyGroup.Any("/:name", RelayMCP)
	}

	// 3. 管理端与终端查询 API：/api/mcp/*
	apiGroup := router.Group("/api/mcp")
	apiGroup.Use(middleware.CORS())
	{
		// 客户端自查已授权的 MCP 列表（供 TPCowork / SkillHub 动态拉取）
		apiGroup.GET("/my-servers", middleware.TokenAuth(), GetMyMCPServers)

		// 管理员后台 CRUD 接口
		adminGroup := apiGroup.Group("")
		adminGroup.Use(middleware.AdminAuth())
		{
			// MCPServer 管理
			adminGroup.GET("/servers", AdminGetAllMCPServers)
			adminGroup.GET("/servers/:id", AdminGetMCPServer)
			adminGroup.POST("/servers", AdminCreateMCPServer)
			adminGroup.PUT("/servers", AdminUpdateMCPServer)
			adminGroup.DELETE("/servers/:id", AdminDeleteMCPServer)

			// MCPRole 工具包管理
			adminGroup.GET("/roles", AdminGetAllMCPRoles)
			adminGroup.GET("/roles/:id", AdminGetMCPRole)
			adminGroup.POST("/roles", AdminCreateMCPRole)
			adminGroup.PUT("/roles", AdminUpdateMCPRole)
			adminGroup.DELETE("/roles/:id", AdminDeleteMCPRole)
		}
	}
}
