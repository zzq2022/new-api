package mcp

import (
	"net/http"
	"strconv"

	"github.com/QuantumNous/new-api/model"
	"github.com/gin-gonic/gin"
)

// RelayMCP 透明反向代理主入口
func RelayMCP(c *gin.Context) {
	serverVal, ok := c.Get("mcp_server")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal error: MCP server context not found",
		})
		return
	}
	server := serverVal.(*MCPServer)

	var token *model.Token
	if tokenVal, ok := c.Get("mcp_token"); ok {
		token = tokenVal.(*model.Token)
	}

	ProxyMCPRequest(c, server, token)
}

// MCPAvailableItem 给终端客户端（TPCowork/SkillHub）返回的可用工具精简结构
type MCPAvailableItem struct {
	Name         string `json:"name"`
	DisplayName  string `json:"display_name"`
	Description  string `json:"description"`
	ProxyURL     string `json:"proxy_url"`
	QuotaPerCall int    `json:"quota_per_call"`
}

// GetMyMCPServers 终端客户端（TPCowork/SkillHub）自查当前用户/Key被授权的可用 MCP 工具列表
func GetMyMCPServers(c *gin.Context) {
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
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized token",
		})
		return
	}

	authorizedServers, err := GetAuthorizedMCPServersForToken(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	items := make([]MCPAvailableItem, 0, len(authorizedServers))
	for _, s := range authorizedServers {
		items = append(items, MCPAvailableItem{
			Name:         s.Name,
			DisplayName:  s.DisplayName,
			Description:  s.Description,
			ProxyURL:     "/mcp/" + s.Name + "/sse",
			QuotaPerCall: s.QuotaPerCall,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    items,
	})
}

// ==========================================
// Admin CRUD: MCPServer 管理
// ==========================================

func AdminGetAllMCPServers(c *gin.Context) {
	p := c.Query("p")
	pageSizeStr := c.Query("page_size")
	page, _ := strconv.Atoi(p)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	servers, total, err := GetAllMCPServers(p != "", page, pageSize)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": gin.H{
			"items": servers,
			"total": total,
		},
	})
}

func AdminGetMCPServer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Invalid ID parameter"})
		return
	}
	server, err := GetMCPServerById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "", "data": server})
}

func AdminCreateMCPServer(c *gin.Context) {
	var server MCPServer
	if err := c.ShouldBindJSON(&server); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	if server.Name == "" || server.BaseURL == "" || server.DisplayName == "" {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Name, DisplayName, and BaseURL are required"})
		return
	}
	if server.Status == 0 {
		server.Status = MCPServerStatusEnabled
	}
	if err := server.Insert(); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "", "data": server})
}

func AdminUpdateMCPServer(c *gin.Context) {
	var server MCPServer
	if err := c.ShouldBindJSON(&server); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	if server.Id <= 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Invalid server ID"})
		return
	}
	if err := server.Update(); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "", "data": server})
}

func AdminDeleteMCPServer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Invalid ID parameter"})
		return
	}
	if err := DeleteMCPServerById(id); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": ""})
}

// ==========================================
// Admin CRUD: MCPRole 工具包管理
// ==========================================

func AdminGetAllMCPRoles(c *gin.Context) {
	p := c.Query("p")
	pageSizeStr := c.Query("page_size")
	page, _ := strconv.Atoi(p)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	roles, total, err := GetAllMCPRoles(p != "", page, pageSize)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": gin.H{
			"items": roles,
			"total": total,
		},
	})
}

func AdminGetMCPRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Invalid ID parameter"})
		return
	}
	role, err := GetMCPRoleById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "", "data": role})
}

func AdminCreateMCPRole(c *gin.Context) {
	var role MCPRole
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	if role.Name == "" || role.DisplayName == "" {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Name and DisplayName are required"})
		return
	}
	if err := role.Insert(); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "", "data": role})
}

func AdminUpdateMCPRole(c *gin.Context) {
	var role MCPRole
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	if role.Id <= 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Invalid role ID"})
		return
	}
	if err := role.Update(); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "", "data": role})
}

func AdminDeleteMCPRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Invalid ID parameter"})
		return
	}
	if err := DeleteMCPRoleById(id); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": ""})
}

// ==========================================
// Admin MCPRole <-> Token 授权关联
// ==========================================

type AuthorizeRoleTokensRequest struct {
	TokenIds []int `json:"token_ids"`
}

func AdminGetMCPRoleTokens(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Invalid ID parameter"})
		return
	}
	role, err := GetMCPRoleById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "MCP role not found"})
		return
	}

	var tokens []model.Token
	if err := model.DB.Select("id", "name", "key", "status", "mcp_roles").Order("id desc").Find(&tokens).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}

	var authorizedIds []int
	type tokenSummary struct {
		Id       int    `json:"id"`
		Name     string `json:"name"`
		Key      string `json:"key"`
		Status   int    `json:"status"`
		Selected bool   `json:"selected"`
	}
	var tokenList []tokenSummary

	for _, t := range tokens {
		roles, _ := t.GetMCPRoles()
		hasRole := false
		for _, r := range roles {
			if r == role.Name {
				hasRole = true
				break
			}
		}
		if hasRole {
			authorizedIds = append(authorizedIds, t.Id)
		}

		maskedKey := t.Key
		if len(t.Key) > 14 {
			maskedKey = t.Key[:8] + "******" + t.Key[len(t.Key)-4:]
		}

		tokenList = append(tokenList, tokenSummary{
			Id:       t.Id,
			Name:     t.Name,
			Key:      maskedKey,
			Status:   t.Status,
			Selected: hasRole,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": gin.H{
			"role":                 role,
			"authorized_token_ids": authorizedIds,
			"tokens":               tokenList,
		},
	})
}

func AdminAuthorizeMCPRoleTokens(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Invalid ID parameter"})
		return
	}
	role, err := GetMCPRoleById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "MCP role not found"})
		return
	}

	var req AuthorizeRoleTokensRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}

	targetMap := make(map[int]bool, len(req.TokenIds))
	for _, tid := range req.TokenIds {
		targetMap[tid] = true
	}

	var tokens []model.Token
	if err := model.DB.Select("id", "mcp_roles").Find(&tokens).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}

	tx := model.DB.Begin()
	for _, t := range tokens {
		roles, _ := t.GetMCPRoles()
		hasRole := false
		filteredRoles := make([]string, 0, len(roles))
		for _, r := range roles {
			if r == role.Name {
				hasRole = true
			} else {
				filteredRoles = append(filteredRoles, r)
			}
		}

		shouldHave := targetMap[t.Id]
		if shouldHave && !hasRole {
			filteredRoles = append(filteredRoles, role.Name)
			_ = t.SetMCPRoles(filteredRoles)
			if err := tx.Model(&model.Token{}).Where("id = ?", t.Id).Update("mcp_roles", t.MCPRoles).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
				return
			}
		} else if !shouldHave && hasRole {
			_ = t.SetMCPRoles(filteredRoles)
			if err := tx.Model(&model.Token{}).Where("id = ?", t.Id).Update("mcp_roles", t.MCPRoles).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
				return
			}
		}
	}
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Tokens authorized successfully"})
}
