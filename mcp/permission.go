package mcp

import (
	"fmt"
	"net/http"

	"github.com/QuantumNous/new-api/model"
	"github.com/gin-gonic/gin"
)

// GetAuthorizedMCPServersForToken 计算给定令牌最终拥有的所有可用 MCPServer 映射
func GetAuthorizedMCPServersForToken(token *model.Token) (map[string]*MCPServer, error) {
	if token == nil {
		return nil, fmt.Errorf("token is nil")
	}

	roleNameSet := make(map[string]struct{})

	// 1. 获取全局默认的角色/工具包（所有有效用户免勾选自动拥有）
	defaultRoles, err := GetDefaultMCPRoles()
	if err == nil {
		for _, r := range defaultRoles {
			roleNameSet[r.Name] = struct{}{}
		}
	}

	// 2. 获取该 Token 显式关联的 MCP 角色代码列表
	tokenRoles, err := token.GetMCPRoles()
	if err == nil {
		for _, rName := range tokenRoles {
			if rName != "" {
				roleNameSet[rName] = struct{}{}
			}
		}
	}

	if len(roleNameSet) == 0 {
		return map[string]*MCPServer{}, nil
	}

	rolesSlice := make([]string, 0, len(roleNameSet))
	for rName := range roleNameSet {
		rolesSlice = append(rolesSlice, rName)
	}

	// 3. 查出这些角色的定义并收集其包含的 MCP 标识
	roles, err := GetMCPRolesByNames(rolesSlice)
	if err != nil {
		return nil, err
	}

	mcpNameSet := make(map[string]struct{})
	for _, role := range roles {
		servers, err := role.GetMCPServers()
		if err != nil {
			continue
		}
		for _, sName := range servers {
			if sName != "" {
				mcpNameSet[sName] = struct{}{}
			}
		}
	}

	if len(mcpNameSet) == 0 {
		return map[string]*MCPServer{}, nil
	}

	mcpNames := make([]string, 0, len(mcpNameSet))
	for name := range mcpNameSet {
		mcpNames = append(mcpNames, name)
	}

	// 4. 查询处于启用状态的目标 MCPServer
	activeServers, err := GetMCPServersByNames(mcpNames)
	if err != nil {
		return nil, err
	}

	result := make(map[string]*MCPServer, len(activeServers))
	for _, s := range activeServers {
		result[s.Name] = s
	}

	return result, nil
}

// MCPPermissionCheck 前置鉴权中间件：核验当前令牌是否有权访问目标 MCP 服务
func MCPPermissionCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		targetName := c.Param("name")
		if targetName == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Missing MCP server name in URL path",
			})
			return
		}

		// 从上下文获取当前认证的 Token 详情
		var token *model.Token
		var err error

		if tokenVal, ok := c.Get("token"); ok {
			if t, valid := tokenVal.(*model.Token); valid && t != nil {
				token = t
			}
		}

		if token == nil {
			tokenId := c.GetInt("token_id")
			if tokenId > 0 {
				token, err = model.GetTokenById(tokenId)
			}
		}

		if err != nil || token == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized: Invalid or missing token",
			})
			return
		}

		// 计算该 Token 当前允许访问的所有 MCP 服务
		authorizedServers, err := GetAuthorizedMCPServersForToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Failed to resolve authorized MCP servers: %s", err.Error()),
			})
			return
		}

		server, ok := authorizedServers[targetName]
		if !ok || server.Status != MCPServerStatusEnabled {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": fmt.Sprintf("Forbidden: You do not have permission to access MCP server '%s'", targetName),
			})
			return
		}

		// 注入上下文供后续 Proxy 及 Audit 消费
		c.Set("mcp_server", server)
		c.Set("mcp_token", token)
		c.Next()
	}
}
