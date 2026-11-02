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

export interface MCPServer {
  id: number
  name: string
  display_name: string
  description?: string
  base_url: string
  quota_per_call: number
  status: number // 1: enabled, 2: disabled
  created_at: number
  updated_at: number
}

export interface MCPRole {
  id: number
  name: string
  display_name: string
  description?: string
  mcp_servers: string // comma-separated server names
  is_default: boolean
  created_at: number
  updated_at: number
}

export interface MCPAvailableItem {
  name: string
  display_name: string
  description: string
  proxy_url: string
  quota_per_call: number
}
