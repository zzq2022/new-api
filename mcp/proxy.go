package mcp

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/logger"
	"github.com/QuantumNous/new-api/model"
	"github.com/gin-gonic/gin"
)

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		if a == "" {
			return b
		}
		return a + "/" + b
	}
	return a + b
}

// ProxyMCPRequest 执行通用透明流式反向代理
func ProxyMCPRequest(c *gin.Context, server *MCPServer, token *model.Token) {
	startTime := time.Now()

	// 1. 解析目标后端真实地址
	targetURL, err := url.Parse(server.BaseURL)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"error": fmt.Sprintf("Invalid target BaseURL configured for MCP '%s': %s", server.Name, err.Error()),
		})
		return
	}

	userId := c.GetInt("id")
	if userId == 0 && token != nil {
		userId = token.UserId
	}

	// 2. 配额原子预检与扣减 (如果配置了 quota_per_call > 0)
	quotaToDeduct := server.QuotaPerCall
	if quotaToDeduct > 0 && token != nil {
		reserved, err := model.TryReserveTokenQuota(token.Id, token.Key, quotaToDeduct, token.UnlimitedQuota)
		if err != nil || !reserved {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": fmt.Sprintf("Token quota is insufficient to access MCP '%s' (requires %d quota)", server.Name, quotaToDeduct),
			})
			return
		}
		// 扣减用户总额度
		_ = model.DecreaseUserQuota(userId, quotaToDeduct, false)
	}

	// 3. 构建透明反向代理
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			subPath := c.Param("any")
			req.URL.Scheme = targetURL.Scheme
			req.URL.Host = targetURL.Host
			req.URL.Path = singleJoiningSlash(targetURL.Path, subPath)

			// 合并 query 参数（如 juyuan 内嵌的 ?token=...）
			if targetURL.RawQuery == "" {
				req.URL.RawQuery = c.Request.URL.RawQuery
			} else if c.Request.URL.RawQuery == "" {
				req.URL.RawQuery = targetURL.RawQuery
			} else {
				req.URL.RawQuery = targetURL.RawQuery + "&" + c.Request.URL.RawQuery
			}

			req.Host = targetURL.Host

			// 注入上游配置的鉴权请求头（覆盖客户端传来的 New API Token）
			if server.Headers != "" {
				var headersMap map[string]string
				if err := common.UnmarshalJsonStr(server.Headers, &headersMap); err == nil {
					for k, v := range headersMap {
						req.Header.Set(k, v)
					}
				}
			}
		},
		FlushInterval: -1, // 立即刷新，确保 SSE 长连接实时推送
		ErrorHandler: func(rw http.ResponseWriter, req *http.Request, proxyErr error) {
			logger.LogError(c, fmt.Sprintf("MCP reverse proxy error for '%s': %v", server.Name, proxyErr))
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusBadGateway)
			_, _ = rw.Write([]byte(fmt.Sprintf(`{"error":"MCP upstream gateway error: %s"}`, proxyErr.Error())))
		},
	}

	// 4. 执行代理转发
	proxy.ServeHTTP(c.Writer, c.Request)

	// 5. 记录 MCP 审计日志
	duration := int(time.Since(startTime).Seconds())
	if duration == 0 {
		duration = 1
	}

	tokenName := c.GetString("token_name")
	if tokenName == "" && token != nil {
		tokenName = token.Name
	}
	tokenId := c.GetInt("token_id")
	if tokenId == 0 && token != nil {
		tokenId = token.Id
	}

	model.RecordMCPLog(c, userId, model.RecordMCPLogParams{
		TokenId:        tokenId,
		TokenName:      tokenName,
		MCPName:        server.Name,
		Quota:          quotaToDeduct,
		UseTimeSeconds: duration,
		Content:        fmt.Sprintf("Method=%s, Path=%s, Status=%d", c.Request.Method, c.Request.URL.Path, c.Writer.Status()),
	})
}
