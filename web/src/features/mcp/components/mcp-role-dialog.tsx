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
import { zodResolver } from '@hookform/resolvers/zod'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import { z } from 'zod'

import { Dialog } from '@/components/dialog'
import { MultiSelect } from '@/components/multi-select'
import { Button } from '@/components/ui/button'
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { handleServerError } from '@/lib/handle-server-error'

import { createMCPRole, listMCPServers, updateMCPRole } from '../api'
import type { MCPRole } from '../types'

const mcpRoleSchema = z.object({
  name: z
    .string()
    .min(1, 'Name is required')
    .regex(/^[a-zA-Z0-9_-]+$/, 'Identifier must contain only letters, numbers, hyphens, and underscores'),
  display_name: z.string().min(1, 'Display name is required'),
  description: z.string().optional().default(''),
  mcp_servers: z.array(z.string()).default([]),
  is_default: z.boolean().default(false),
})

type MCPRoleFormValues = z.infer<typeof mcpRoleSchema>

type MCPRoleDialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  role?: MCPRole | null
}

export function MCPRoleDialog({
  open,
  onOpenChange,
  role,
}: MCPRoleDialogProps) {
  const { t } = useTranslation()
  const queryClient = useQueryClient()
  const isEditing = !!role

  const { data: serversData } = useQuery({
    queryKey: ['mcp-servers'],
    queryFn: () => listMCPServers(1, 100),
    enabled: open,
  })

  const serverOptions = (serversData?.items ?? []).map((s) => ({
    label: `${s.display_name} (${s.name})`,
    value: s.name,
    hint: s.base_url,
  }))

  const form = useForm<MCPRoleFormValues>({
    resolver: zodResolver(mcpRoleSchema),
    defaultValues: {
      name: '',
      display_name: '',
      description: '',
      mcp_servers: [],
      is_default: false,
    },
  })

  useEffect(() => {
    if (role) {
      form.reset({
        name: role.name,
        display_name: role.display_name,
        description: role.description ?? '',
        mcp_servers: role.mcp_servers
          ? role.mcp_servers.split(',').filter(Boolean)
          : [],
        is_default: role.is_default ?? false,
      })
    } else {
      form.reset({
        name: '',
        display_name: '',
        description: '',
        mcp_servers: [],
        is_default: false,
      })
    }
  }, [role, form, open])

  const mutation = useMutation({
    mutationFn: async (values: MCPRoleFormValues) => {
      const payload: Partial<MCPRole> = {
        name: values.name,
        display_name: values.display_name,
        description: values.description,
        mcp_servers: values.mcp_servers.join(','),
        is_default: values.is_default,
      }
      if (isEditing && role) {
        return updateMCPRole({
          ...role,
          ...payload,
        })
      }
      return createMCPRole(payload)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['mcp-roles'] })
      toast.success(
        isEditing
          ? t('MCP role updated successfully')
          : t('MCP role created successfully')
      )
      onOpenChange(false)
    },
    onError: (err) => {
      handleServerError(err)
    },
  })

  const onSubmit = (values: MCPRoleFormValues) => {
    mutation.mutate(values)
  }

  return (
    <Dialog
      open={open}
      onOpenChange={onOpenChange}
      title={isEditing ? t('Edit MCP Toolkit Role') : t('Add MCP Toolkit Role')}
      description={t(
        'Define a toolkit package bundling multiple MCP services to assign to API keys'
      )}
      contentClassName='max-w-lg'
    >
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className='space-y-4'>
          <FormField
            control={form.control}
            name='name'
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t('Identifier Name')}</FormLabel>
                <FormControl>
                  <Input
                    {...field}
                    placeholder='e.g. developer'
                    disabled={isEditing}
                  />
                </FormControl>
                <FormDescription>
                  {t('Unique identifier for this MCP role')}
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name='display_name'
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t('Display Name')}</FormLabel>
                <FormControl>
                  <Input {...field} placeholder='e.g. Developer Toolkit' />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name='mcp_servers'
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t('Included MCP Servers')}</FormLabel>
                <FormControl>
                  <MultiSelect
                    options={serverOptions}
                    selected={field.value}
                    onChange={field.onChange}
                    placeholder={t('Select MCP servers to include in this toolkit')}
                  />
                </FormControl>
                <FormDescription>
                  {t('All selected MCP services will be accessible to keys with this role')}
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name='description'
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t('Description')}</FormLabel>
                <FormControl>
                  <Textarea
                    {...field}
                    placeholder={t('Describe the intended usage or audience for this toolkit')}
                    rows={2}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <div className='flex justify-end gap-2 pt-2'>
            <Button
              type='button'
              variant='outline'
              onClick={() => onOpenChange(false)}
            >
              {t('Cancel')}
            </Button>
            <Button type='submit' disabled={mutation.isPending}>
              {mutation.isPending ? t('Saving...') : t('Save')}
            </Button>
          </div>
        </form>
      </Form>
    </Dialog>
  )
}
