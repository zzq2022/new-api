# 需求规格说明书：API Key 业务名称扩展与 MCP 权限严格收紧

## Problem Statement

在目前的企业大模型网关中，API Key（令牌）的管理与 MCP（Model Context Protocol）工具权限分发存在以下两项业务与安全痛点：

1. **业务命名与系统代号混杂，缺乏直观中文标识**：
   - 目前系统的 API Key 仅具备单一的 `name` 字段（在控制台界面表头显示为“名称”）。在实际企业使用中，该字段通常被用于录入员工账号拼音、工号或调用方系统短代码（如 `zzq01`）。
   - 随着团队与系统对接数量增加，仅凭英文代码难以直接辨识密钥对应的人员中文姓名（如 `张之清01`）或业务归属，导致管理员日常检索、排查调用日志与分配权限时辨识成本高且易发生混淆。

2. **未配置 MCP 的 Key 隐式继承权限，违背最小特权与显式授权原则**：
   - 在既有的 MCP 权限管理中，系统支持将特定工具包设置为“全员默认拥有（`is_default = true`）”。
   - 当某个 API Key 未被配置任何 MCP 角色时，前端列表将其展示为“默认 (Default)”；而在网关鉴权与客户端自查端点中，系统会自动将所有默认工具包下发给该 Key。
   - 用户明确提出“每个 API Key 如果没有配置 MCP，不能是默认，没有就是没有”。对于不需要使用内网 MCP 工具的常规模型调用 Key，隐式赋予公共工具不仅造成权限冗余，也增加了敏感内部接口的暴露面；而在客户端拉取工具配置时，未配置的 Key 应当直接接收空配置（`[]`），而非被动挂载默认工具。

## Solution

在统一模型网关中实施针对 API Key 数据模型、搜索逻辑、界面交互与 MCP 权限鉴权的端到端演进：

1. **数据模型与展示解耦（业务名称 vs 标识代码）**：
   - 在数据实体 `tokens` 中新增物理字段 `apikey_name`（业务名称/姓名），保留既有的 `name` 字段作为系统标识代码。
   - 在控制台 API Key 列表中将原本单一的列拆分为两列独立呈现：**密钥名称**（展示 `apikey_name`，如 `张之清01`）与 **标识代码**（展示 `name`，如 `zzq01`）。
   - 在新建与编辑 API Key 抽屉中，提供独立的“密钥名称”输入项，与“标识代码”共同作为必填项，支持在批量创建时继承与后缀递增。
   - 列表顶部的搜索能力同步升级为双向复合检索，支持输入中文姓名或英文代号快速过滤。

2. **MCP 显式授权与全员默认机制彻底移除**：
   - 彻底废弃 MCP 角色的“全员默认（`is_default`）”隐式赋权机制，所有 MCP 权限必须由管理员在 API Key 上显式勾选绑定，或在角色管理中通过批量反向授权打钩分配。
   - 前端 API Key 列表中的“MCP 角色”列，对于未绑定角色的 Key，状态文案由“默认”彻底纠正为“**无**”。
   - 网关层与客户端配置下发逻辑收紧：
     - **配置自查接口 (`GET /api/mcp/my-servers`)**：未配置 MCP 角色的 Key 收到 `HTTP 200` 响应及**空配置列表 (`data: []`)**，客户端（如 TPCowork）正常启动且不挂载任何工具，大模型聊天调用不受任何影响；
     - **代理鉴权前置拦截 (`/mcp/:name/*any`)**：仅当请求的目标 MCP 存在于该 Key 显式关联的有效角色集合中时放行，未授权的具体工具调用请求一律硬阻断并返回 `403 Forbidden`。

## User Stories

1. 作为一个企业的系统管理员，我想要在为员工创建 API Key 时录入其实名或中文业务名称（如 `张之清01`），以便直观标识该 Key 的实际归属人与业务用途。
2. 作为一个企业的系统管理员，我想要在创建 API Key 时保留短英文代号（如 `zzq01`），以便下游程序或调用脚本保持简短规范的标识符引用。
3. 作为一个企业的系统管理员，我想要在 API Key 列表表格中独立看到“密钥名称”列与“标识代码”列，以便在同一视图下清晰区分中文姓名和英文代号。
4. 作为一个企业的系统管理员，我想要在控制台顶部搜索框中输入“张之清”或“zzq”时均能精确过滤出对应记录，以便根据任意已知信息快速定位 Key。
5. 作为一个企业的系统管理员，我想要在编辑已有 API Key 时能够更新“密钥名称”，以便在员工岗位调整或组织架构变更时平滑维护信息。
6. 作为一个企业的系统管理员，我想要在创建/编辑抽屉中将“密钥名称”与“标识代码”均作为必填校验，以便杜绝空白无效名称的数据录入。
7. 作为一个企业的系统管理员，我想要未分配任何 MCP 工具包的 API Key 在列表中的 MCP 角色列明确标记为“无”，以便一眼确认该 Key 不具备任何 MCP 访问权限。
8. 作为一个企业的安全审计员，我想要未配置 MCP 角色的 API Key 不自动继承任何隐式默认工具，以便遵循最小权限原则并杜绝越权访问内部系统。
9. 作为一个终端员工，我想要在使用未分配 MCP 工具的 API Key 登录桌面客户端时，中台自查接口正常返回空工具清单（`[]`），以便客户端正常运行大模型对话且不报任何权限错误。
10. 作为一个终端员工，我想要日常大模型问答（`/v1/chat/completions`）完全不受 MCP 是否配置的影响，以便保持核心对话功能的稳定高可用。
11. 作为一个企业的安全审计员，我想要当未经授权的 API Key 试图发起具体 MCP 服务的代理连接时，网关立即阻断并返回标准的 403 Forbidden 响应，以便防范非法探测与接口滥用。
12. 作为一个企业的系统管理员，我想要在 MCP 角色管理的“反向授权 API Key”弹窗中同时看到各 Key 的中文密钥名称与标识代码，以便准确批量打钩授权而不误伤其他同名账号。
13. 作为一个企业的系统管理员，我想要系统彻底移除 MCP 角色的“全员默认拥有”勾选属性，以便全网关的工具访问权限模型保持单一、纯粹的显式分配机制。
14. 作为一个企业的运维工程师，我想要存量历史 API Key 在系统启动升级后能够自动平滑迁移，历史记录的 `apikey_name` 默认为空且不影响任何调用，以便无缝完成版本升级。

## Implementation Decisions

### 1. 数据库模型与数据访问层决策
- **模型扩展**：
  在 `model.Token` 结构体中新增字段：
  ```go
  ApikeyName string `json:"apikey_name" gorm:"type:varchar(128);index"`
  ```
  该字段使用 `varchar(128)` 并建立普通索引，通过项目主数据库的 GORM 自动迁移机制（`AutoMigrate`）在应用启动时幂等追加，同时完美兼容 SQLite、MySQL（>= 5.7.8）和 PostgreSQL（>= 9.6）。存量数据自动赋空字符串，不产生破坏性影响。
- **搜索逻辑演进**：
  在 `model.SearchUserTokens` 中，当传入有效搜索关键字时，由原本单一匹配 `name` 扩展为复合条件：
  `(apikey_name LIKE ? ESCAPE '!' OR name LIKE ? ESCAPE '!')`，保持现有的防注入与通配符转义逻辑，同时满足中文名称与代号的双向模糊检索。

### 2. API 控制器与交互契约决策
- **令牌 CRUD 控制器 (`controller/token.go`)**：
  - `AddToken`：从客户端请求中读取 `apikey_name`，经过修剪后注入 `cleanToken.ApikeyName` 落库；
  - `UpdateToken`：允许更新 `apikey_name` 字段，并在审计参数中记录变更；
  - `buildMaskedTokenResponse`：返回结构体透传 `apikey_name`。
- **角色授权查询控制器 (`mcp/controller.go`)**：
  - `AdminGetMCPRoleTokens` 中的 `tokenSummary` 增加 `ApikeyName` 字段并从数据库读出，供角色授权弹窗前端展示。

### 3. MCP 鉴权内核与配置下发收紧决策
- **移除隐式默认工具继承 (`mcp/permission.go`)**：
  彻底移除 `GetAuthorizedMCPServersForToken` 中的 `GetDefaultMCPRoles()` 合并逻辑。当前 Token 所拥有的可用工具集合**完全由且仅由其显式关联的角色列表**决定：
  $$\text{当前可用 MCP 集合} = \bigcup_{r \in \text{Token.mcp\_roles}} r.\text{mcp\_servers}$$
  若 `Token.mcp_roles` 为空，直接短路返回空映射 `map[string]*MCPServer{}`。
- **配置下发行为 (`mcp/controller.go` - `GetMyMCPServers`)**：
  对于空映射，接口以 `HTTP 200 OK` 正常响应 JSON 结构 `{"success": true, "data": []}`。
- **代理拦截行为 (`mcp/permission.go` - `MCPPermissionCheck`)**：
  当路由指向 `/mcp/:name/*any` 时，若目标 `:name` 不存在于当前 Token 的授权服务映射中，网关通过 `c.AbortWithStatusJSON(http.StatusForbidden, ...)` 硬阻断并返回标准 403 错误信息。

### 4. 前端交互与视觉排版决策
- **独立双列布局 (`web/src/features/keys/components/api-keys-columns.tsx`)**：
  - 第一列：**密钥名称**（`accessorKey: 'apikey_name'`），表头使用 `t('API Key Name')`，主要展示中文姓名（如 `张之清01`）；
  - 第二列：**标识代码**（`accessorKey: 'name'`），表头调整为 `t('Code')` 或 `t('Key ID')`，展示短代号（如 `zzq01`）。
- **MCP 角色空状态文案修复 (`web/src/features/keys/components/api-keys-cells.tsx`)**：
  - `McpRolesCell` 中检测到未配置角色时，`label` 与 `emptyLabel` 均展示为 `t('None')`（中文下翻译为“**无**”），彻底废除“默认”展示。
- **表单录入与校验更新 (`web/src/features/keys/`)**：
  - 在 `api-key-form.ts` 中，`apiKeySchema` 增加 `apikey_name: z.string().min(1, t('Please enter an API key name'))`；
  - 在 `api-keys-mutate-drawer.tsx` 的基本信息区顶部并排或顺序放置“密钥名称”与“标识代码”两个输入控件；在创建模式（包含单次与批量）中均提供默认占位符与必填校验。
- **角色授权弹窗增强 (`web/src/features/mcp/components/mcp-role-authorize-dialog.tsx`)**：
  - 授权列表项中渲染格式优化为 `张之清01 (zzq01)`，方便管理员核对人员与代号。
- **多语言字典对齐 (`web/src/i18n/locales/*.json`)**：
  - 同步补充 `API Key Name`、`Code`、`None` 等词条在 7 种语言中的标准翻译。

## Testing Decisions

### 1. 良好测试准则与边界
- 遵循“只测试外部暴露契约与不可变约束，不测试内部易变私有实现”原则。
- 重点针对：
  1. 权限硬隔离与空配置下发的正确性；
  2. 显式授权后的可用性；
  3. `apikey_name` 与 `name` 的持久化与双向检索准确性；
  4. 跨 SQLite 内存数据库的迁移兼容性。

### 2. 核心测试场景覆盖
- **MCP 显式授权与无权限测试 (`mcp/mcp_test.go`)**：
  - Case 1（未配置 MCP 的 Token）：断言 `GetAuthorizedMCPServersForToken` 返回 0 个服务，断言调用 `GetMyMCPServers` 返回 HTTP 200 且 `data` 为 `[]`，断言向 `/mcp/:name/sse` 发起请求必定返回 403 Forbidden；
  - Case 2（配置了 MCP 角色的 Token）：断言能且仅能访问该角色包含的目标 MCP，其他未分配工具返回 403。
- **Token 数据模型与双向检索测试 (`controller/token_test.go`)**：
  - 创建携带 `apikey_name="张之清01"` 与 `name="zzq01"` 的 Token，断言接口返回及落库数据中两字段准确无误；
  - 分别以关键词 `"张之清"` 和 `"zzq"` 检索，断言均能命中并准确分页返回该 Token。

### 3. 参考规范与已有工程实践
- 数据库字段与迁移复用 `model/token.go` 已有的 GORM 标签风格。
- 检索测试继承 `model/token.go` 中 `SearchUserTokens` 的安全转义与分页断言规范。

## Out of Scope

1. **不做大模型计费与渠道分组的改动**：本次仅涉及 API Key 基础属性扩充与 MCP 协议网关鉴权，与大模型渠道的分组倍率、额度扣费逻辑完全解耦。
2. **不做底层 API Key 密文生成规则变更**：Key 的 `sk-xxxx` 哈希与加密校验流程保持不变。
3. **不做存量数据的强制全量回填**：允许历史存量 Token 的 `apikey_name` 为空字符串，由管理员后续按需在界面编辑补充。

## Further Notes

- **数据库向后兼容**：本次为 `tokens` 表新增字段为非破坏性的增加列，现有生产环境在升级重新启动时 GORM `AutoMigrate` 会自动执行 `ALTER TABLE tokens ADD COLUMN apikey_name ...`，数据零丢失、业务零停机。
