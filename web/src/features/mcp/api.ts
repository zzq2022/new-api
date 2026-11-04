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
import { api, type ApiRequestConfig } from '@/lib/api'
import { requireServerSuccess } from '@/lib/server-error-message'

import type { MCPServer, MCPRole } from './types'

const mutationConfig: ApiRequestConfig = {
  skipBusinessError: true,
  skipErrorHandler: true,
}

export async function listMCPServers(p?: number, pageSize?: number) {
  const response = await api.get<{
    success: boolean
    message: string
    data: { items: MCPServer[]; total: number }
  }>('/api/mcp/servers', {
    params: { p, page_size: pageSize },
  })
  return requireServerSuccess(response.data).data
}

export async function getMCPServer(id: number) {
  const response = await api.get<{
    success: boolean
    message: string
    data: MCPServer
  }>(`/api/mcp/servers/${id}`)
  return requireServerSuccess(response.data).data
}

export async function createMCPServer(data: Partial<MCPServer>) {
  const response = await api.post<{
    success: boolean
    message: string
    data: MCPServer
  }>('/api/mcp/servers', data, mutationConfig)
  return requireServerSuccess(response.data).data
}

export async function updateMCPServer(data: Partial<MCPServer>) {
  const response = await api.put<{
    success: boolean
    message: string
    data: MCPServer
  }>('/api/mcp/servers', data, mutationConfig)
  return requireServerSuccess(response.data).data
}

export async function deleteMCPServer(id: number) {
  const response = await api.delete<{
    success: boolean
    message: string
  }>(`/api/mcp/servers/${id}`, mutationConfig)
  return requireServerSuccess(response.data)
}

export async function listMCPRoles(p?: number, pageSize?: number) {
  const response = await api.get<{
    success: boolean
    message: string
    data: { items: MCPRole[]; total: number }
  }>('/api/mcp/roles', {
    params: { p, page_size: pageSize },
  })
  return requireServerSuccess(response.data).data
}

export async function getMCPRole(id: number) {
  const response = await api.get<{
    success: boolean
    message: string
    data: MCPRole
  }>(`/api/mcp/roles/${id}`)
  return requireServerSuccess(response.data).data
}

export async function createMCPRole(data: Partial<MCPRole>) {
  const response = await api.post<{
    success: boolean
    message: string
    data: MCPRole
  }>('/api/mcp/roles', data, mutationConfig)
  return requireServerSuccess(response.data).data
}

export async function updateMCPRole(data: Partial<MCPRole>) {
  const response = await api.put<{
    success: boolean
    message: string
    data: MCPRole
  }>('/api/mcp/roles', data, mutationConfig)
  return requireServerSuccess(response.data).data
}

export async function deleteMCPRole(id: number) {
  const response = await api.delete<{
    success: boolean
    message: string
  }>(`/api/mcp/roles/${id}`, mutationConfig)
  return requireServerSuccess(response.data)
}

export type MCPRoleTokensResponse = {
  role: MCPRole
  authorized_token_ids: number[]
  tokens: Array<{
    id: number
    apikey_name?: string
    name: string
    key: string
    status: number
    selected: boolean
  }>
}

export async function getMCPRoleTokens(id: number) {
  const response = await api.get<{
    success: boolean
    message: string
    data: MCPRoleTokensResponse
  }>(`/api/mcp/roles/${id}/tokens`)
  return requireServerSuccess(response.data).data
}

export async function authorizeMCPRoleTokens(id: number, tokenIds: number[]) {
  const response = await api.post<{
    success: boolean
    message: string
  }>(`/api/mcp/roles/${id}/tokens`, { token_ids: tokenIds }, mutationConfig)
  return requireServerSuccess(response.data)
}
