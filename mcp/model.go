package mcp

import (
	"errors"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"gorm.io/gorm"
)

type MCPServer struct {
	Id           int            `json:"id" gorm:"primaryKey;autoIncrement"`
	Name         string         `json:"name" gorm:"type:varchar(64);uniqueIndex;not null"` // 英文唯一路由标识
	DisplayName  string         `json:"display_name" gorm:"type:varchar(128);not null"`    // 显示名称
	Description  string         `json:"description" gorm:"type:text"`                      // 描述说明
	BaseURL      string         `json:"base_url" gorm:"type:varchar(512);not null"`        // 内网真实目标基础地址
	Headers      string         `json:"headers" gorm:"type:text"`                          // 上游鉴权与自定义请求头 (JSON 格式)
	Status       int            `json:"status" gorm:"type:int;default:1;index"`            // 1: 启用, 2: 停用
	QuotaPerCall int            `json:"quota_per_call" gorm:"type:int;default:0"`          // 单次调用扣减额度，0 表示免费
	CreatedAt    int64          `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    int64          `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

const (
	MCPServerStatusEnabled  = 1
	MCPServerStatusDisabled = 2
)

type MCPRole struct {
	Id          int            `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string         `json:"name" gorm:"type:varchar(64);uniqueIndex;not null"` // 角色代码，如 developer
	DisplayName string         `json:"display_name" gorm:"type:varchar(128);not null"`    // 角色名称，如 研发工具包
	Description string         `json:"description" gorm:"type:text"`                      // 描述说明
	MCPServers  string         `json:"mcp_servers" gorm:"type:text"`                      // 包含的 MCP Name 列表，JSON 数组如 ["internal_kb", "db_query"]
	IsDefault   bool           `json:"is_default" gorm:"default:false;index"`             // 是否全员默认拥有（免勾选自动赋予）
	CreatedAt   int64          `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   int64          `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func InitDatabase() error {
	if model.DB == nil {
		return nil
	}
	if err := model.DB.AutoMigrate(&MCPServer{}, &MCPRole{}); err != nil {
		return err
	}
	return SeedBuiltinPresets()
}

// ---- MCPServer Operations ----

func (server *MCPServer) Insert() error {
	var count int64
	model.DB.Model(&MCPServer{}).Where("name = ?", server.Name).Count(&count)
	if count > 0 {
		return errors.New("MCP server with the same name already exists")
	}
	return model.DB.Create(server).Error
}

func (server *MCPServer) Update() error {
	var count int64
	model.DB.Model(&MCPServer{}).Where("name = ? AND id != ?", server.Name, server.Id).Count(&count)
	if count > 0 {
		return errors.New("MCP server with the same name already exists")
	}
	return model.DB.Model(server).Select("name", "display_name", "description", "base_url", "status", "quota_per_call").Updates(server).Error
}

func DeleteMCPServerById(id int) error {
	return model.DB.Delete(&MCPServer{}, "id = ?", id).Error
}

func GetMCPServerById(id int) (*MCPServer, error) {
	var server MCPServer
	err := model.DB.Where("id = ?", id).First(&server).Error
	return &server, err
}

func GetMCPServerByName(name string) (*MCPServer, error) {
	var server MCPServer
	err := model.DB.Where("name = ?", name).First(&server).Error
	return &server, err
}

func GetMCPServersByNames(names []string) ([]*MCPServer, error) {
	if len(names) == 0 {
		return []*MCPServer{}, nil
	}
	var servers []*MCPServer
	err := model.DB.Where("name IN ? AND status = ?", names, MCPServerStatusEnabled).Find(&servers).Error
	return servers, err
}

func GetAllMCPServers(paged bool, page int, pageSize int) ([]*MCPServer, int64, error) {
	var servers []*MCPServer
	var total int64
	query := model.DB.Model(&MCPServer{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if paged {
		if page < 1 {
			page = 1
		}
		if pageSize < 1 {
			pageSize = 10
		}
		query = query.Offset((page - 1) * pageSize).Limit(pageSize)
	}
	err := query.Order("id desc").Find(&servers).Error
	return servers, total, err
}

// ---- MCPRole Operations ----

func (role *MCPRole) GetMCPServers() ([]string, error) {
	trimmed := strings.TrimSpace(role.MCPServers)
	if trimmed == "" {
		return []string{}, nil
	}
	var servers []string
	if strings.HasPrefix(trimmed, "[") {
		if err := common.UnmarshalJsonStr(trimmed, &servers); err == nil {
			return servers, nil
		}
	}
	parts := strings.Split(trimmed, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			servers = append(servers, p)
		}
	}
	return servers, nil
}

func (role *MCPRole) SetMCPServers(servers []string) error {
	if len(servers) == 0 {
		role.MCPServers = "[]"
		return nil
	}
	data, err := common.Marshal(servers)
	if err != nil {
		return err
	}
	role.MCPServers = string(data)
	return nil
}

func (role *MCPRole) Insert() error {
	var count int64
	model.DB.Model(&MCPRole{}).Where("name = ?", role.Name).Count(&count)
	if count > 0 {
		return errors.New("MCP role with the same name already exists")
	}
	if role.MCPServers == "" {
		role.MCPServers = "[]"
	}
	return model.DB.Create(role).Error
}

func (role *MCPRole) Update() error {
	var count int64
	model.DB.Model(&MCPRole{}).Where("name = ? AND id != ?", role.Name, role.Id).Count(&count)
	if count > 0 {
		return errors.New("MCP role with the same name already exists")
	}
	if role.MCPServers == "" {
		role.MCPServers = "[]"
	}
	return model.DB.Model(role).Select("name", "display_name", "description", "mcp_servers", "is_default").Updates(role).Error
}

func DeleteMCPRoleById(id int) error {
	return model.DB.Delete(&MCPRole{}, "id = ?", id).Error
}

func GetMCPRoleById(id int) (*MCPRole, error) {
	var role MCPRole
	err := model.DB.Where("id = ?", id).First(&role).Error
	return &role, err
}

func GetMCPRoleByName(name string) (*MCPRole, error) {
	var role MCPRole
	err := model.DB.Where("name = ?", name).First(&role).Error
	return &role, err
}

func GetDefaultMCPRoles() ([]*MCPRole, error) {
	var roles []*MCPRole
	err := model.DB.Where("is_default = ?", true).Find(&roles).Error
	return roles, err
}

func GetMCPRolesByNames(names []string) ([]*MCPRole, error) {
	if len(names) == 0 {
		return []*MCPRole{}, nil
	}
	var roles []*MCPRole
	err := model.DB.Where("name IN ?", names).Find(&roles).Error
	return roles, err
}

func GetAllMCPRoles(paged bool, page int, pageSize int) ([]*MCPRole, int64, error) {
	var roles []*MCPRole
	var total int64
	query := model.DB.Model(&MCPRole{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if paged {
		if page < 1 {
			page = 1
		}
		if pageSize < 1 {
			pageSize = 10
		}
		query = query.Offset((page - 1) * pageSize).Limit(pageSize)
	}
	err := query.Order("id desc").Find(&roles).Error
	return roles, total, err
}
