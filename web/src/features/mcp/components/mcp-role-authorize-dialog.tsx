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
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { Check, Info, Key, Search, ShieldCheck } from 'lucide-react'
import { useEffect, useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import { Dialog } from '@/components/dialog'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import { Input } from '@/components/ui/input'
import { ScrollArea } from '@/components/ui/scroll-area'
import { handleServerError } from '@/lib/handle-server-error'

import { authorizeMCPRoleTokens, getMCPRoleTokens } from '../api'
import type { MCPRole } from '../types'

type MCPRoleAuthorizeDialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  role: MCPRole | null
}

export function MCPRoleAuthorizeDialog({
  open,
  onOpenChange,
  role,
}: MCPRoleAuthorizeDialogProps) {
  const { t } = useTranslation()
  const queryClient = useQueryClient()
  const [search, setSearch] = useState('')
  const [selectedIds, setSelectedIds] = useState<number[]>([])

  const { data, isLoading } = useQuery({
    queryKey: ['mcp-role-tokens', role?.id],
    queryFn: () => getMCPRoleTokens(role!.id),
    enabled: open && !!role,
  })

  useEffect(() => {
    if (data?.authorized_token_ids) {
      setSelectedIds(data.authorized_token_ids)
    }
  }, [data])

  const mutation = useMutation({
    mutationFn: () => authorizeMCPRoleTokens(role!.id, selectedIds),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['mcp-role-tokens', role?.id] })
      queryClient.invalidateQueries({ queryKey: ['api-keys'] })
      toast.success(t('API keys authorization updated successfully'))
      onOpenChange(false)
    },
    onError: (err) => {
      handleServerError(err)
    },
  })

  const tokens = data?.tokens ?? []
  const filteredTokens = useMemo(() => {
    if (!search.trim()) return tokens
    const q = search.toLowerCase()
    return tokens.filter(
      (tk) =>
        tk.name.toLowerCase().includes(q) || tk.key.toLowerCase().includes(q)
    )
  }, [tokens, search])

  const handleToggle = (id: number) => {
    setSelectedIds((prev) =>
      prev.includes(id) ? prev.filter((item) => item !== id) : [...prev, id]
    )
  }

  const handleSelectAll = () => {
    const visibleIds = filteredTokens.map((tk) => tk.id)
    const allSelected = visibleIds.every((id) => selectedIds.includes(id))
    if (allSelected) {
      setSelectedIds((prev) => prev.filter((id) => !visibleIds.includes(id)))
    } else {
      setSelectedIds((prev) => Array.from(new Set([...prev, ...visibleIds])))
    }
  }

  if (!role) return null

  return (
    <Dialog
      open={open}
      onOpenChange={onOpenChange}
      title={
        <div className='flex items-center gap-2'>
          <ShieldCheck className='size-5 text-primary' />
          <span>{t('Authorize API Keys for Role')}</span>
        </div>
      }
      description={t('Manage which API keys have access to {{name}} ({{code}}).', {
        name: role.display_name,
        code: role.name,
      })}
      footer={
        <div className='flex justify-end gap-2'>
          <Button
            type='button'
            variant='outline'
            onClick={() => onOpenChange(false)}
          >
            {t('Cancel')}
          </Button>
          <Button
            type='button'
            onClick={() => mutation.mutate()}
            disabled={mutation.isPending || isLoading}
          >
            {mutation.isPending ? t('Saving...') : t('Save Authorization')}
          </Button>
        </div>
      }
    >
      <div className='space-y-4 py-2'>
        {role.is_default && (
          <div className='flex items-start gap-2 rounded-md border border-primary/20 bg-primary/5 p-3 text-xs text-primary'>
            <Info className='size-4 shrink-0 mt-0.5' />
            <div>
              <p className='font-semibold'>
                {t('This role is configured as a Global Default Role.')}
              </p>
              <p className='text-muted-foreground mt-0.5'>
                {t(
                  'All active API keys automatically have access without explicit authorization. You can still assign specific keys here.'
                )}
              </p>
            </div>
          </div>
        )}

        <div className='flex items-center justify-between gap-2'>
          <div className='relative flex-1'>
            <Search className='absolute left-2.5 top-2.5 size-4 text-muted-foreground' />
            <Input
              placeholder={t('Search API keys by name or key...')}
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className='pl-8 h-9'
            />
          </div>
          {filteredTokens.length > 0 && (
            <Button
              variant='outline'
              size='sm'
              onClick={handleSelectAll}
              className='h-9 text-xs shrink-0'
            >
              {filteredTokens.every((tk) => selectedIds.includes(tk.id))
                ? t('Deselect All')
                : t('Select All')}
            </Button>
          )}
        </div>

        <div className='rounded-md border'>
          <ScrollArea className='h-64 p-2'>
            {isLoading ? (
              <div className='flex h-32 items-center justify-center text-sm text-muted-foreground'>
                {t('Loading API keys...')}
              </div>
            ) : filteredTokens.length === 0 ? (
              <div className='flex h-32 items-center justify-center text-sm text-muted-foreground'>
                {t('No matching API keys found')}
              </div>
            ) : (
              <div className='space-y-1.5'>
                {filteredTokens.map((token) => {
                  const isChecked = selectedIds.includes(token.id)
                  return (
                    <div
                      key={token.id}
                      onClick={() => handleToggle(token.id)}
                      className={`flex items-center justify-between rounded-md p-2.5 transition-colors cursor-pointer border ${
                        isChecked
                          ? 'border-primary/40 bg-primary/5'
                          : 'border-transparent hover:bg-muted/50'
                      }`}
                    >
                      <div className='flex items-center gap-3 min-w-0'>
                        <Checkbox
                          checked={isChecked}
                          onCheckedChange={() => handleToggle(token.id)}
                        />
                        <div className='min-w-0'>
                          <div className='flex items-center gap-1.5'>
                            <span className='font-medium text-sm truncate'>
                              {token.name}
                            </span>
                            {token.status === 1 ? (
                              <Badge
                                variant='outline'
                                className='text-[10px] px-1 py-0 h-4 border-emerald-500/30 text-emerald-600 dark:text-emerald-400'
                              >
                                {t('Enabled')}
                              </Badge>
                            ) : (
                              <Badge
                                variant='secondary'
                                className='text-[10px] px-1 py-0 h-4 text-muted-foreground'
                              >
                                {t('Disabled')}
                              </Badge>
                            )}
                          </div>
                          <div className='flex items-center gap-1 font-mono text-xs text-muted-foreground'>
                            <Key className='size-3 shrink-0' />
                            <span>{token.key}</span>
                          </div>
                        </div>
                      </div>

                      {isChecked && (
                        <Check className='size-4 text-primary shrink-0' />
                      )}
                    </div>
                  )
                })}
              </div>
            )}
          </ScrollArea>
        </div>

        <div className='flex items-center justify-between text-xs text-muted-foreground'>
          <span>
            {t('Selected: {{count}} of {{total}} keys', {
              count: selectedIds.length,
              total: tokens.length,
            })}
          </span>
        </div>
      </div>
    </Dialog>
  )
}
