package mcp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupTestMCPDB(t *testing.T) *gorm.DB {
	common.RedisEnabled = false
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	model.DB = db
	model.LOG_DB = db
	require.NoError(t, db.AutoMigrate(
		&MCPServer{},
		&MCPRole{},
		&model.Token{},
		&model.User{},
		&model.Log{},
	))
	return db
}

func TestMCPServerAndRoleCRUD(t *testing.T) {
	setupTestMCPDB(t)

	// 1. 创建 MCP Server
	s1 := &MCPServer{
		Name:         "internal_kb",
		DisplayName:  "企业内部知识库",
		Description:  "内部文档检索",
		BaseURL:      "http://127.0.0.1:8000",
		Status:       MCPServerStatusEnabled,
		QuotaPerCall: 0,
	}
	require.NoError(t, s1.Insert())
	assert.Greater(t, s1.Id, 0)

	// 重名校验
	sDuplicate := &MCPServer{
		Name:        "internal_kb",
		DisplayName: "重复知识库",
		BaseURL:     "http://127.0.0.1:8001",
	}
	require.Error(t, sDuplicate.Insert())

	// 查询
	found, err := GetMCPServerByName("internal_kb")
	require.NoError(t, err)
	assert.Equal(t, "企业内部知识库", found.DisplayName)

	// 2. 创建 MCP Role
	r1 := &MCPRole{
		Name:        "dev_pack",
		DisplayName: "研发工具包",
		IsDefault:   false,
	}
	require.NoError(t, r1.SetMCPServers([]string{"internal_kb"}))
	require.NoError(t, r1.Insert())

	// 3. 全员默认角色
	rDefault := &MCPRole{
		Name:        "all_staff",
		DisplayName: "全员默认包",
		IsDefault:   true,
	}
	require.NoError(t, rDefault.SetMCPServers([]string{"internal_kb"}))
	require.NoError(t, rDefault.Insert())

	defaultRoles, err := GetDefaultMCPRoles()
	require.NoError(t, err)
	assert.Len(t, defaultRoles, 1)
	assert.Equal(t, "all_staff", defaultRoles[0].Name)
}

func TestPermissionCheckAndResolution(t *testing.T) {
	setupTestMCPDB(t)

	// 准备 MCP 实例
	kb := &MCPServer{
		Name:        "internal_kb",
		DisplayName: "企业知识库",
		BaseURL:     "http://127.0.0.1:8000",
		Status:      MCPServerStatusEnabled,
	}
	require.NoError(t, kb.Insert())

	sqlTool := &MCPServer{
		Name:        "sql_tool",
		DisplayName: "数据库查询工具",
		BaseURL:     "http://127.0.0.1:8001",
		Status:      MCPServerStatusEnabled,
	}
	require.NoError(t, sqlTool.Insert())

	// 创建研发工具包包含 sql_tool
	devRole := &MCPRole{
		Name:        "dev",
		DisplayName: "研发",
		IsDefault:   false,
	}
	require.NoError(t, devRole.SetMCPServers([]string{"sql_tool"}))
	require.NoError(t, devRole.Insert())

	// 创建全员默认工具包包含 internal_kb
	defaultRole := &MCPRole{
		Name:        "public",
		DisplayName: "公共工具",
		IsDefault:   true,
	}
	require.NoError(t, defaultRole.SetMCPServers([]string{"internal_kb"}))
	require.NoError(t, defaultRole.Insert())

	// Case 1: 普通用户令牌（未关联 dev 角色）
	normalToken := &model.Token{
		Id:     1,
		UserId: 10,
		Key:    "sk-normal",
		Name:   "zzq_normal",
	}
	require.NoError(t, normalToken.SetMCPRoles([]string{}))
	require.NoError(t, model.DB.Create(normalToken).Error)

	serversMap, err := GetAuthorizedMCPServersForToken(normalToken)
	require.NoError(t, err)
	assert.Contains(t, serversMap, "internal_kb", "普通令牌应自动获得默认公共工具")
	assert.NotContains(t, serversMap, "sql_tool", "普通令牌不应具有研发专属工具")

	// Case 2: 研发人员令牌（关联 dev 角色）
	devToken := &model.Token{
		Id:     2,
		UserId: 11,
		Key:    "sk-dev",
		Name:   "zzq_dev",
	}
	require.NoError(t, devToken.SetMCPRoles([]string{"dev"}))
	require.NoError(t, model.DB.Create(devToken).Error)

	devServers, err := GetAuthorizedMCPServersForToken(devToken)
	require.NoError(t, err)
	assert.Contains(t, devServers, "internal_kb", "研发人员同时拥有公共工具")
	assert.Contains(t, devServers, "sql_tool", "研发人员拥有研发工具")
}

func TestMCPTransparentReverseProxy(t *testing.T) {
	setupTestMCPDB(t)

	// 1. 模拟内网真实 MCP 服务 (Echo Server)
	mockUpstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Upstream-Path", r.URL.Path)
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("event: message\ndata: hello from mock mcp\n\n"))
	}))
	defer mockUpstream.Close()

	// 注册内网工具指向 mockUpstream
	server := &MCPServer{
		Name:        "test_mcp",
		DisplayName: "测试工具",
		BaseURL:     mockUpstream.URL,
		Status:      MCPServerStatusEnabled,
	}
	require.NoError(t, server.Insert())

	// 全员默认开放
	role := &MCPRole{
		Name:      "all",
		IsDefault: true,
	}
	require.NoError(t, role.SetMCPServers([]string{"test_mcp"}))
	require.NoError(t, role.Insert())

	token := &model.Token{
		Id:     1,
		UserId: 1,
		Key:    "sk-test-proxy",
		Name:   "test_token",
	}
	require.NoError(t, model.DB.Create(token).Error)

	// 构建 Gin Engine
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(func(c *gin.Context) {
		// 模拟 TokenAuth 中间件注入上下文
		c.Set("id", token.UserId)
		c.Set("token_id", token.Id)
		c.Set("token_key", token.Key)
		c.Set("token_name", token.Name)
		c.Set("token", token)
		c.Next()
	})
	engine.Any("/mcp/:name/*any", MCPPermissionCheck(), RelayMCP)

	// 2. 正常代理测试 (GET /mcp/test_mcp/sse)
	req := httptest.NewRequest("GET", "/mcp/test_mcp/sse", nil)
	rec := newCloseNotifyingRecorder()
	engine.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "/sse", rec.Header().Get("X-Upstream-Path"))
	assert.Contains(t, rec.Body.String(), "hello from mock mcp")

	// 3. 越权拦截测试 (GET /mcp/unauthorized_tool/sse)
	reqForbidden := httptest.NewRequest("GET", "/mcp/unauthorized_tool/sse", nil)
	recForbidden := newCloseNotifyingRecorder()
	engine.ServeHTTP(recForbidden, reqForbidden)

	assert.Equal(t, http.StatusForbidden, recForbidden.Code)
	var errResp map[string]string
	require.NoError(t, json.Unmarshal(recForbidden.Body.Bytes(), &errResp))
	assert.Contains(t, errResp["error"], "Forbidden")
}

type closeNotifyingRecorder struct {
	*httptest.ResponseRecorder
	closed chan bool
}

func newCloseNotifyingRecorder() *closeNotifyingRecorder {
	return &closeNotifyingRecorder{
		ResponseRecorder: httptest.NewRecorder(),
		closed:           make(chan bool, 1),
	}
}

func (c *closeNotifyingRecorder) CloseNotify() <-chan bool {
	return c.closed
}
