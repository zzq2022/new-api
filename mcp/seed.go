/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/
package mcp

import (
	"github.com/QuantumNous/new-api/model"
)

// BuiltinPreset 定义内置 MCP 服务预设
type BuiltinPreset struct {
	Name        string
	DisplayName string
	Description string
	BaseURL     string
	Headers     string
}

var BuiltinPresets = []BuiltinPreset{
	{
		Name:        "juyuan",
		DisplayName: "聚源智能助手",
		Description: "聚源金融数据与分析智能助手 API 服务",
		BaseURL:     "https://api.gildata.com/mcp-servers/aidata-assistant-srv-api?token=in3ga4hF3mOwlIvmknOskDy3bto3xS3f",
		Headers:     "{}",
	},
	{
		Name:        "caihui_mcp",
		DisplayName: "财汇数据助手",
		Description: "财汇资讯金融宏观与行业数据检索服务",
		BaseURL:     "https://mcp.finchina.com/finchina-data-mcp-server/mcp",
		Headers:     `{"x-api-key":"ak_4571f4d7da4d485eb35a2670b34a7d61"}`,
	},
	{
		Name:        "hexin-ifind-ds-stock-mcp",
		DisplayName: "同花顺 iFinD 股票",
		Description: "A股、港股、美股行情、财务与深度个股研究数据",
		BaseURL:     "https://api-mcp.51ifind.com:8643/ds-mcp-servers/hexin-ifind-ds-stock-mcp",
		Headers:     `{"Authorization":"eyJraWQiOiJtY3AtYXBpIiwidWlkIjoiODYyMDg1ODM1IiwiYWxnIjoiUlNBLU9BRVAtMjU2IiwiZW5jIjoiQTI1NkdDTSJ9.V3_RyxRpf3ZqDlKoy9TjyglvEKHdA3kO514siQ5IDeCrlSZOPtfzqeme7YKQqj9jrm4qPwoBQABCSFfF357YVBahRMJ7DriRh9tL8n-x0NCBKiTXkcLMfP21poszZBXYU8cl05FC2vqaXYuM21ssKD1uew0YizuZrR4oVeB1Om6dLcuSZPiz-uRB9V21bzUb137Cfsn79aCJ8nXxbXVkuHspdNabAw5qtZUgjcQL3-U4Xpuc78WTZqSvXQWJTVPc_J-sX432_qO2g22CrO0cbnAKSeFeDWF5eZJyhqriBeUxwaauGsa3jGPJQhG-Fpedk00muyXI0QVnsENLyZ-Qlw.IxjdSSkf8x4jxzd4.mNJ8X0u1L9KQNBqCrV7lubsfH4lzaVHp9uogPMF7PjXWouOVDBiJV10mH-ePymBPykyVAjP3wqUciAfqrQJaChs3sRKvaO6Sx0g1B6jxxNnPu4Ny2rAZpdUDdN8D8UzF_RMU154cRd1_NQGYO5CelGOOTsB8-Chp1FHCUGx-Igktok9QXGIq_lPM4RK4CSjoH4UjN7I2jE-Ed1fzJUwitpUz29ei9d7qk90csgCaJiauAUpfIGuFLoX8bROywO8BnH-TJxqlED32743AWy4hWn7r7U8MKIueWqKPEccXQzs3CQ_4tKeTnfBorLy_2UIy2OP2fXwZ1Jl4bQ_As6m-OJOvAzaTxeBcxAIVZtplImh3sew8TqIrGQ0i6Y-t9ruJ6oBJWOnKMWhT85j67ODFkiJQ-rK3G1yfNvu8EN85v-MyqJVjDPw6fK8UnjTrTbwpQQID0Q.bKsTjm5LyraeNnjW47AwFQ"}`,
	},
	{
		Name:        "hexin-ifind-ds-fund-mcp",
		DisplayName: "同花顺 iFinD 基金",
		Description: "公募基金净值、持仓穿透、经理业绩与评级数据",
		BaseURL:     "https://api-mcp.51ifind.com:8643/ds-mcp-servers/hexin-ifind-ds-fund-mcp",
		Headers:     `{"Authorization":"eyJraWQiOiJtY3AtYXBpIiwidWlkIjoiODYyMDg1ODM1IiwiYWxnIjoiUlNBLU9BRVAtMjU2IiwiZW5jIjoiQTI1NkdDTSJ9.V3_RyxRpf3ZqDlKoy9TjyglvEKHdA3kO514siQ5IDeCrlSZOPtfzqeme7YKQqj9jrm4qPwoBQABCSFfF357YVBahRMJ7DriRh9tL8n-x0NCBKiTXkcLMfP21poszZBXYU8cl05FC2vqaXYuM21ssKD1uew0YizuZrR4oVeB1Om6dLcuSZPiz-uRB9V21bzUb137Cfsn79aCJ8nXxbXVkuHspdNabAw5qtZUgjcQL3-U4Xpuc78WTZqSvXQWJTVPc_J-sX432_qO2g22CrO0cbnAKSeFeDWF5eZJyhqriBeUxwaauGsa3jGPJQhG-Fpedk00muyXI0QVnsENLyZ-Qlw.IxjdSSkf8x4jxzd4.mNJ8X0u1L9KQNBqCrV7lubsfH4lzaVHp9uogPMF7PjXWouOVDBiJV10mH-ePymBPykyVAjP3wqUciAfqrQJaChs3sRKvaO6Sx0g1B6jxxNnPu4Ny2rAZpdUDdN8D8UzF_RMU154cRd1_NQGYO5CelGOOTsB8-Chp1FHCUGx-Igktok9QXGIq_lPM4RK4CSjoH4UjN7I2jE-Ed1fzJUwitpUz29ei9d7qk90csgCaJiauAUpfIGuFLoX8bROywO8BnH-TJxqlED32743AWy4hWn7r7U8MKIueWqKPEccXQzs3CQ_4tKeTnfBorLy_2UIy2OP2fXwZ1Jl4bQ_As6m-OJOvAzaTxeBcxAIVZtplImh3sew8TqIrGQ0i6Y-t9ruJ6oBJWOnKMWhT85j67ODFkiJQ-rK3G1yfNvu8EN85v-MyqJVjDPw6fK8UnjTrTbwpQQID0Q.bKsTjm5LyraeNnjW47AwFQ"}`,
	},
	{
		Name:        "hexin-ifind-ds-edb-mcp",
		DisplayName: "同花顺 iFinD 宏观",
		Description: "宏观经济指标库、宏观 EDB、行业高频数据与预测",
		BaseURL:     "https://api-mcp.51ifind.com:8643/ds-mcp-servers/hexin-ifind-ds-edb-mcp",
		Headers:     `{"Authorization":"eyJraWQiOiJtY3AtYXBpIiwidWlkIjoiODYyMDg1ODM1IiwiYWxnIjoiUlNBLU9BRVAtMjU2IiwiZW5jIjoiQTI1NkdDTSJ9.V3_RyxRpf3ZqDlKoy9TjyglvEKHdA3kO514siQ5IDeCrlSZOPtfzqeme7YKQqj9jrm4qPwoBQABCSFfF357YVBahRMJ7DriRh9tL8n-x0NCBKiTXkcLMfP21poszZBXYU8cl05FC2vqaXYuM21ssKD1uew0YizuZrR4oVeB1Om6dLcuSZPiz-uRB9V21bzUb137Cfsn79aCJ8nXxbXVkuHspdNabAw5qtZUgjcQL3-U4Xpuc78WTZqSvXQWJTVPc_J-sX432_qO2g22CrO0cbnAKSeFeDWF5eZJyhqriBeUxwaauGsa3jGPJQhG-Fpedk00muyXI0QVnsENLyZ-Qlw.IxjdSSkf8x4jxzd4.mNJ8X0u1L9KQNBqCrV7lubsfH4lzaVHp9uogPMF7PjXWouOVDBiJV10mH-ePymBPykyVAjP3wqUciAfqrQJaChs3sRKvaO6Sx0g1B6jxxNnPu4Ny2rAZpdUDdN8D8UzF_RMU154cRd1_NQGYO5CelGOOTsB8-Chp1FHCUGx-Igktok9QXGIq_lPM4RK4CSjoH4UjN7I2jE-Ed1fzJUwitpUz29ei9d7qk90csgCaJiauAUpfIGuFLoX8bROywO8BnH-TJxqlED32743AWy4hWn7r7U8MKIueWqKPEccXQzs3CQ_4tKeTnfBorLy_2UIy2OP2fXwZ1Jl4bQ_As6m-OJOvAzaTxeBcxAIVZtplImh3sew8TqIrGQ0i6Y-t9ruJ6oBJWOnKMWhT85j67ODFkiJQ-rK3G1yfNvu8EN85v-MyqJVjDPw6fK8UnjTrTbwpQQID0Q.bKsTjm5LyraeNnjW47AwFQ"}`,
	},
	{
		Name:        "hexin-ifind-ds-news-mcp",
		DisplayName: "同花顺 iFinD 新闻",
		Description: "金融财经实时快讯、公司公告、投研研报与舆情事件",
		BaseURL:     "https://api-mcp.51ifind.com:8643/ds-mcp-servers/hexin-ifind-ds-news-mcp",
		Headers:     `{"Authorization":"eyJraWQiOiJtY3AtYXBpIiwidWlkIjoiODYyMDg1ODM1IiwiYWxnIjoiUlNBLU9BRVAtMjU2IiwiZW5jIjoiQTI1NkdDTSJ9.V3_RyxRpf3ZqDlKoy9TjyglvEKHdA3kO514siQ5IDeCrlSZOPtfzqeme7YKQqj9jrm4qPwoBQABCSFfF357YVBahRMJ7DriRh9tL8n-x0NCBKiTXkcLMfP21poszZBXYU8cl05FC2vqaXYuM21ssKD1uew0YizuZrR4oVeB1Om6dLcuSZPiz-uRB9V21bzUb137Cfsn79aCJ8nXxbXVkuHspdNabAw5qtZUgjcQL3-U4Xpuc78WTZqSvXQWJTVPc_J-sX432_qO2g22CrO0cbnAKSeFeDWF5eZJyhqriBeUxwaauGsa3jGPJQhG-Fpedk00muyXI0QVnsENLyZ-Qlw.IxjdSSkf8x4jxzd4.mNJ8X0u1L9KQNBqCrV7lubsfH4lzaVHp9uogPMF7PjXWouOVDBiJV10mH-ePymBPykyVAjP3wqUciAfqrQJaChs3sRKvaO6Sx0g1B6jxxNnPu4Ny2rAZpdUDdN8D8UzF_RMU154cRd1_NQGYO5CelGOOTsB8-Chp1FHCUGx-Igktok9QXGIq_lPM4RK4CSjoH4UjN7I2jE-Ed1fzJUwitpUz29ei9d7qk90csgCaJiauAUpfIGuFLoX8bROywO8BnH-TJxqlED32743AWy4hWn7r7U8MKIueWqKPEccXQzs3CQ_4tKeTnfBorLy_2UIy2OP2fXwZ1Jl4bQ_As6m-OJOvAzaTxeBcxAIVZtplImh3sew8TqIrGQ0i6Y-t9ruJ6oBJWOnKMWhT85j67ODFkiJQ-rK3G1yfNvu8EN85v-MyqJVjDPw6fK8UnjTrTbwpQQID0Q.bKsTjm5LyraeNnjW47AwFQ"}`,
	},
	{
		Name:        "hexin-ifind-ds-bond-mcp",
		DisplayName: "同花顺 iFinD 债券",
		Description: "国债、信用债、可转债行情收益率曲线与违约预警",
		BaseURL:     "https://api-mcp.51ifind.com:8643/ds-mcp-servers/hexin-ifind-ds-bond-mcp",
		Headers:     `{"Authorization":"eyJraWQiOiJtY3AtYXBpIiwidWlkIjoiODYyMDg1ODM1IiwiYWxnIjoiUlNBLU9BRVAtMjU2IiwiZW5jIjoiQTI1NkdDTSJ9.V3_RyxRpf3ZqDlKoy9TjyglvEKHdA3kO514siQ5IDeCrlSZOPtfzqeme7YKQqj9jrm4qPwoBQABCSFfF357YVBahRMJ7DriRh9tL8n-x0NCBKiTXkcLMfP21poszZBXYU8cl05FC2vqaXYuM21ssKD1uew0YizuZrR4oVeB1Om6dLcuSZPiz-uRB9V21bzUb137Cfsn79aCJ8nXxbXVkuHspdNabAw5qtZUgjcQL3-U4Xpuc78WTZqSvXQWJTVPc_J-sX432_qO2g22CrO0cbnAKSeFeDWF5eZJyhqriBeUxwaauGsa3jGPJQhG-Fpedk00muyXI0QVnsENLyZ-Qlw.IxjdSSkf8x4jxzd4.mNJ8X0u1L9KQNBqCrV7lubsfH4lzaVHp9uogPMF7PjXWouOVDBiJV10mH-ePymBPykyVAjP3wqUciAfqrQJaChs3sRKvaO6Sx0g1B6jxxNnPu4Ny2rAZpdUDdN8D8UzF_RMU154cRd1_NQGYO5CelGOOTsB8-Chp1FHCUGx-Igktok9QXGIq_lPM4RK4CSjoH4UjN7I2jE-Ed1fzJUwitpUz29ei9d7qk90csgCaJiauAUpfIGuFLoX8bROywO8BnH-TJxqlED32743AWy4hWn7r7U8MKIueWqKPEccXQzs3CQ_4tKeTnfBorLy_2UIy2OP2fXwZ1Jl4bQ_As6m-OJOvAzaTxeBcxAIVZtplImh3sew8TqIrGQ0i6Y-t9ruJ6oBJWOnKMWhT85j67ODFkiJQ-rK3G1yfNvu8EN85v-MyqJVjDPw6fK8UnjTrTbwpQQID0Q.bKsTjm5LyraeNnjW47AwFQ"}`,
	},
	{
		Name:        "hexin-ifind-ds-global-stock-mcp",
		DisplayName: "同花顺 iFinD 全球股票",
		Description: "全球主要股票市场、指数、成分股及跨境联动数据",
		BaseURL:     "https://api-mcp.51ifind.com:8643/ds-mcp-servers/hexin-ifind-ds-global-stock-mcp",
		Headers:     `{"Authorization":"eyJraWQiOiJtY3AtYXBpIiwidWlkIjoiODYyMDg1ODM1IiwiYWxnIjoiUlNBLU9BRVAtMjU2IiwiZW5jIjoiQTI1NkdDTSJ9.V3_RyxRpf3ZqDlKoy9TjyglvEKHdA3kO514siQ5IDeCrlSZOPtfzqeme7YKQqj9jrm4qPwoBQABCSFfF357YVBahRMJ7DriRh9tL8n-x0NCBKiTXkcLMfP21poszZBXYU8cl05FC2vqaXYuM21ssKD1uew0YizuZrR4oVeB1Om6dLcuSZPiz-uRB9V21bzUb137Cfsn79aCJ8nXxbXVkuHspdNabAw5qtZUgjcQL3-U4Xpuc78WTZqSvXQWJTVPc_J-sX432_qO2g22CrO0cbnAKSeFeDWF5eZJyhqriBeUxwaauGsa3jGPJQhG-Fpedk00muyXI0QVnsENLyZ-Qlw.IxjdSSkf8x4jxzd4.mNJ8X0u1L9KQNBqCrV7lubsfH4lzaVHp9uogPMF7PjXWouOVDBiJV10mH-ePymBPykyVAjP3wqUciAfqrQJaChs3sRKvaO6Sx0g1B6jxxNnPu4Ny2rAZpdUDdN8D8UzF_RMU154cRd1_NQGYO5CelGOOTsB8-Chp1FHCUGx-Igktok9QXGIq_lPM4RK4CSjoH4UjN7I2jE-Ed1fzJUwitpUz29ei9d7qk90csgCaJiauAUpfIGuFLoX8bROywO8BnH-TJxqlED32743AWy4hWn7r7U8MKIueWqKPEccXQzs3CQ_4tKeTnfBorLy_2UIy2OP2fXwZ1Jl4bQ_As6m-OJOvAzaTxeBcxAIVZtplImh3sew8TqIrGQ0i6Y-t9ruJ6oBJWOnKMWhT85j67ODFkiJQ-rK3G1yfNvu8EN85v-MyqJVjDPw6fK8UnjTrTbwpQQID0Q.bKsTjm5LyraeNnjW47AwFQ"}`,
	},
	{
		Name:        "hexin-ifind-ds-index-mcp",
		DisplayName: "同花顺 iFinD 指数",
		Description: "指数基本信息、权重走势、行业分类及指数估值体系",
		BaseURL:     "https://api-mcp.51ifind.com:8643/ds-mcp-servers/hexin-ifind-ds-index-mcp",
		Headers:     `{"Authorization":"eyJraWQiOiJtY3AtYXBpIiwidWlkIjoiODYyMDg1ODM1IiwiYWxnIjoiUlNBLU9BRVAtMjU2IiwiZW5jIjoiQTI1NkdDTSJ9.V3_RyxRpf3ZqDlKoy9TjyglvEKHdA3kO514siQ5IDeCrlSZOPtfzqeme7YKQqj9jrm4qPwoBQABCSFfF357YVBahRMJ7DriRh9tL8n-x0NCBKiTXkcLMfP21poszZBXYU8cl05FC2vqaXYuM21ssKD1uew0YizuZrR4oVeB1Om6dLcuSZPiz-uRB9V21bzUb137Cfsn79aCJ8nXxbXVkuHspdNabAw5qtZUgjcQL3-U4Xpuc78WTZqSvXQWJTVPc_J-sX432_qO2g22CrO0cbnAKSeFeDWF5eZJyhqriBeUxwaauGsa3jGPJQhG-Fpedk00muyXI0QVnsENLyZ-Qlw.IxjdSSkf8x4jxzd4.mNJ8X0u1L9KQNBqCrV7lubsfH4lzaVHp9uogPMF7PjXWouOVDBiJV10mH-ePymBPykyVAjP3wqUciAfqrQJaChs3sRKvaO6Sx0g1B6jxxNnPu4Ny2rAZpdUDdN8D8UzF_RMU154cRd1_NQGYO5CelGOOTsB8-Chp1FHCUGx-Igktok9QXGIq_lPM4RK4CSjoH4UjN7I2jE-Ed1fzJUwitpUz29ei9d7qk90csgCaJiauAUpfIGuFLoX8bROywO8BnH-TJxqlED32743AWy4hWn7r7U8MKIueWqKPEccXQzs3CQ_4tKeTnfBorLy_2UIy2OP2fXwZ1Jl4bQ_As6m-OJOvAzaTxeBcxAIVZtplImh3sew8TqIrGQ0i6Y-t9ruJ6oBJWOnKMWhT85j67ODFkiJQ-rK3G1yfNvu8EN85v-MyqJVjDPw6fK8UnjTrTbwpQQID0Q.bKsTjm5LyraeNnjW47AwFQ"}`,
	},
	{
		Name:        "anysearch",
		DisplayName: "AnySearch 联网搜索",
		Description: "高质量互联网实时聚合搜索与网页内容抽取服务",
		BaseURL:     "https://api.anysearch.com/mcp",
		Headers:     `{"Authorization":"Bearer as_sk_0e0ce195a1ee5ea7b07b0d226ed81547"}`,
	},
}

// SeedBuiltinPresets 检查并初始化内置 MCP 服务与默认角色
func SeedBuiltinPresets() error {
	if model.DB == nil {
		return nil
	}

	for _, preset := range BuiltinPresets {
		var existing MCPServer
		err := model.DB.Where("name = ?", preset.Name).First(&existing).Error
		if err != nil {
			// 未找到则自动创建
			server := MCPServer{
				Name:         preset.Name,
				DisplayName:  preset.DisplayName,
				Description:  preset.Description,
				BaseURL:      preset.BaseURL,
				Headers:      preset.Headers,
				Status:       MCPServerStatusEnabled,
				QuotaPerCall: 0,
			}
			_ = model.DB.Create(&server).Error
		}
	}

	// 1. 初始化默认金融投研工具包
	financeServers := "juyuan,caihui_mcp,hexin-ifind-ds-stock-mcp,hexin-ifind-ds-fund-mcp,hexin-ifind-ds-edb-mcp,hexin-ifind-ds-news-mcp,hexin-ifind-ds-bond-mcp,hexin-ifind-ds-global-stock-mcp,hexin-ifind-ds-index-mcp"
	var financeRole MCPRole
	if err := model.DB.Where("name = ?", "finance_suite").First(&financeRole).Error; err != nil {
		_ = model.DB.Create(&MCPRole{
			Name:        "finance_suite",
			DisplayName: "金融投研工具包",
			Description: "包含同花顺 iFinD、财汇数据与聚源金融的全面金融数据能力",
			MCPServers:  financeServers,
			IsDefault:   true,
		}).Error
	}

	// 2. 初始化默认联网搜索工具包
	var searchRole MCPRole
	if err := model.DB.Where("name = ?", "search_suite").First(&searchRole).Error; err != nil {
		_ = model.DB.Create(&MCPRole{
			Name:        "search_suite",
			DisplayName: "联网搜索工具包",
			Description: "包含 AnySearch 高质量实时互联网聚合搜索服务",
			MCPServers:  "anysearch",
			IsDefault:   true,
		}).Error
	}

	return nil
}
