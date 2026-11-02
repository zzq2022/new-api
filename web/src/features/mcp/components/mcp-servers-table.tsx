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
import { Edit2, MoreHorizontal, Trash2 } from 'lucide-react'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import { ConfirmDialog } from '@/components/confirm-dialog'
import { CopyButton } from '@/components/copy-button'
import { EmptyState } from '@/components/empty-state'
import { LoadingState } from '@/components/loading-state'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { handleServerError } from '@/lib/handle-server-error'

import { deleteMCPServer, listMCPServers } from '../api'
import type { MCPServer } from '../types'
import { MCPServerDialog } from './mcp-server-dialog'

export function MCPServersTable() {
  const { t } = useTranslation()
  const queryClient = useQueryClient()
  const [editingServer, setEditingServer] = useState<MCPServer | null>(null)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [deleteTarget, setDeleteTarget] = useState<MCPServer | null>(null)

  const { data, isLoading } = useQuery({
    queryKey: ['mcp-servers'],
    queryFn: () => listMCPServers(1, 100),
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => deleteMCPServer(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['mcp-servers'] })
      toast.success(t('MCP server deleted successfully'))
      setDeleteTarget(null)
    },
    onError: (err) => {
      handleServerError(err)
    },
  })

  if (isLoading) {
    return <LoadingState />
  }

  const servers = data?.items ?? []

  if (servers.length === 0) {
    return (
      <EmptyState
        title={t('No MCP Servers Found')}
        description={t('Register your first enterprise MCP server to get started.')}
      />
    )
  }

  return (
    <>
      <div className='rounded-md border'>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className='w-16'>ID</TableHead>
              <TableHead>{t('Identifier')}</TableHead>
              <TableHead>{t('Display Name')}</TableHead>
              <TableHead>{t('Upstream Base URL')}</TableHead>
              <TableHead>{t('Gateway Proxy Endpoint')}</TableHead>
              <TableHead className='w-32'>{t('Quota / Call')}</TableHead>
              <TableHead className='w-24'>{t('Status')}</TableHead>
              <TableHead className='w-20 text-right'>{t('Actions')}</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {servers.map((server) => {
              const proxyPath = `/mcp/${server.name}/sse`
              return (
                <TableRow key={server.id}>
                  <TableCell className='font-mono text-xs text-muted-foreground'>
                    {server.id}
                  </TableCell>
                  <TableCell>
                    <div className='flex items-center gap-1.5'>
                      <span className='font-mono font-medium'>{server.name}</span>
                      <CopyButton value={server.name} />
                    </div>
                  </TableCell>
                  <TableCell>
                    <div className='font-medium'>{server.display_name}</div>
                    {server.description && (
                      <div className='text-xs text-muted-foreground line-clamp-1'>
                        {server.description}
                      </div>
                    )}
                  </TableCell>
                  <TableCell>
                    <div className='flex items-center gap-1.5 font-mono text-xs text-muted-foreground'>
                      <span className='max-w-[200px] truncate' title={server.base_url}>
                        {server.base_url}
                      </span>
                      <CopyButton value={server.base_url} />
                    </div>
                  </TableCell>
                  <TableCell>
                    <div className='flex items-center gap-1.5 font-mono text-xs text-muted-foreground'>
                      <span>{proxyPath}</span>
                      <CopyButton value={proxyPath} />
                    </div>
                  </TableCell>
                  <TableCell>
                    <Badge variant='outline' className='font-mono'>
                      {server.quota_per_call ?? 0}
                    </Badge>
                  </TableCell>
                  <TableCell>
                    {server.status === 1 ? (
                      <Badge variant='outline' className='border-emerald-500/30 text-emerald-600 dark:text-emerald-400 bg-emerald-500/10'>
                        {t('Enabled')}
                      </Badge>
                    ) : (
                      <Badge variant='outline' className='border-muted text-muted-foreground'>
                        {t('Disabled')}
                      </Badge>
                    )}
                  </TableCell>
                  <TableCell className='text-right'>
                    <DropdownMenu>
                      <DropdownMenuTrigger
                        render={
                          <Button variant='ghost' size='icon' className='size-8'>
                            <MoreHorizontal className='size-4' />
                          </Button>
                        }
                      />
                      <DropdownMenuContent align='end'>
                        <DropdownMenuItem
                          onClick={() => {
                            setEditingServer(server)
                            setDialogOpen(true)
                          }}
                        >
                          <Edit2 className='mr-2 size-4' />
                          {t('Edit')}
                        </DropdownMenuItem>
                        <DropdownMenuItem
                          className='text-destructive focus:text-destructive'
                          onClick={() => setDeleteTarget(server)}
                        >
                          <Trash2 className='mr-2 size-4' />
                          {t('Delete')}
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </TableCell>
                </TableRow>
              )
            })}
          </TableBody>
        </Table>
      </div>

      <MCPServerDialog
        open={dialogOpen}
        onOpenChange={(open) => {
          setDialogOpen(open)
          if (!open) setEditingServer(null)
        }}
        server={editingServer}
      />

      <ConfirmDialog
        open={!!deleteTarget}
        onOpenChange={(open) => {
          if (!open) setDeleteTarget(null)
        }}
        title={t('Delete MCP Server')}
        desc={t(
          'Are you sure you want to delete MCP server "{{name}}"? This action cannot be undone.',
          { name: deleteTarget?.name ?? '' }
        )}
        destructive
        isLoading={deleteMutation.isPending}
        handleConfirm={() => {
          if (deleteTarget) deleteMutation.mutate(deleteTarget.id)
        }}
      />
    </>
  )
}
