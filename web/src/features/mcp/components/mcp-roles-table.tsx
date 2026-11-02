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
import { Edit2, Key, MoreHorizontal, Trash2 } from 'lucide-react'
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

import { deleteMCPRole, listMCPRoles } from '../api'
import type { MCPRole } from '../types'
import { MCPRoleAuthorizeDialog } from './mcp-role-authorize-dialog'
import { MCPRoleDialog } from './mcp-role-dialog'

export function MCPRolesTable() {
  const { t } = useTranslation()
  const queryClient = useQueryClient()
  const [editingRole, setEditingRole] = useState<MCPRole | null>(null)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [authorizingRole, setAuthorizingRole] = useState<MCPRole | null>(null)
  const [authorizeOpen, setAuthorizeOpen] = useState(false)
  const [deleteTarget, setDeleteTarget] = useState<MCPRole | null>(null)

  const { data, isLoading } = useQuery({
    queryKey: ['mcp-roles'],
    queryFn: () => listMCPRoles(1, 100),
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => deleteMCPRole(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['mcp-roles'] })
      toast.success(t('MCP role deleted successfully'))
      setDeleteTarget(null)
    },
    onError: (err) => {
      handleServerError(err)
    },
  })

  if (isLoading) {
    return <LoadingState />
  }

  const roles = data?.items ?? []

  if (roles.length === 0) {
    return (
      <EmptyState
        title={t('No MCP Roles Found')}
        description={t(
          'Create your first toolkit role to package MCP servers for users.'
        )}
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
              <TableHead>{t('Bundled Servers')}</TableHead>
              <TableHead className='w-32'>{t('Default Role')}</TableHead>
              <TableHead className='w-20 text-right'>{t('Actions')}</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {roles.map((role) => {
              const serversList = role.mcp_servers
                ? role.mcp_servers.split(',').filter(Boolean)
                : []

              return (
                <TableRow key={role.id}>
                  <TableCell className='font-mono text-xs text-muted-foreground'>
                    {role.id}
                  </TableCell>
                  <TableCell>
                    <div className='flex items-center gap-1.5'>
                      <span className='font-mono font-medium'>{role.name}</span>
                      <CopyButton value={role.name} />
                    </div>
                  </TableCell>
                  <TableCell>
                    <div className='font-medium'>{role.display_name}</div>
                    {role.description && (
                      <div className='text-xs text-muted-foreground line-clamp-1'>
                        {role.description}
                      </div>
                    )}
                  </TableCell>
                  <TableCell>
                    <div className='flex flex-wrap gap-1'>
                      {serversList.length > 0 ? (
                        serversList.map((srv) => (
                          <Badge
                            key={srv}
                            variant='secondary'
                            className='font-mono text-xs'
                          >
                            {srv}
                          </Badge>
                        ))
                      ) : (
                        <span className='text-xs text-muted-foreground italic'>
                          {t('None')}
                        </span>
                      )}
                    </div>
                  </TableCell>
                  <TableCell>
                    {role.is_default ? (
                      <Badge
                        variant='outline'
                        className='border-primary/40 text-primary bg-primary/10'
                      >
                        {t('Default')}
                      </Badge>
                    ) : (
                      <span className='text-xs text-muted-foreground'>-</span>
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
                            setAuthorizingRole(role)
                            setAuthorizeOpen(true)
                          }}
                        >
                          <Key className='mr-2 size-4' />
                          {t('Authorize API Keys')}
                        </DropdownMenuItem>
                        <DropdownMenuItem
                          onClick={() => {
                            setEditingRole(role)
                            setDialogOpen(true)
                          }}
                        >
                          <Edit2 className='mr-2 size-4' />
                          {t('Edit')}
                        </DropdownMenuItem>
                        <DropdownMenuItem
                          className='text-destructive focus:text-destructive'
                          onClick={() => setDeleteTarget(role)}
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

      <MCPRoleAuthorizeDialog
        open={authorizeOpen}
        onOpenChange={(open) => {
          setAuthorizeOpen(open)
          if (!open) setAuthorizingRole(null)
        }}
        role={authorizingRole}
      />

      <MCPRoleDialog
        open={dialogOpen}
        onOpenChange={(open) => {
          setDialogOpen(open)
          if (!open) setEditingRole(null)
        }}
        role={editingRole}
      />

      <ConfirmDialog
        open={!!deleteTarget}
        onOpenChange={(open) => {
          if (!open) setDeleteTarget(null)
        }}
        title={t('Delete MCP Role')}
        desc={t(
          'Are you sure you want to delete role "{{name}}"? Tokens with this role will lose access to its servers.',
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
