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
import { Plus } from 'lucide-react'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'

import { SectionPageLayout } from '@/components/layout'
import { Button } from '@/components/ui/button'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'

import { MCPRoleDialog } from './components/mcp-role-dialog'
import { MCPRolesTable } from './components/mcp-roles-table'
import { MCPServerDialog } from './components/mcp-server-dialog'
import { MCPServersTable } from './components/mcp-servers-table'

export function MCPManagement() {
  const { t } = useTranslation()
  const [activeTab, setActiveTab] = useState<'servers' | 'roles'>('servers')
  const [serverDialogOpen, setServerDialogOpen] = useState(false)
  const [roleDialogOpen, setRoleDialogOpen] = useState(false)

  return (
    <>
      <SectionPageLayout fixedContent>
        <SectionPageLayout.Title>
          {t('MCP Management')}
        </SectionPageLayout.Title>
        <SectionPageLayout.Actions>
          {activeTab === 'servers' ? (
            <Button onClick={() => setServerDialogOpen(true)}>
              <Plus className='mr-1.5 size-4' />
              {t('Add MCP Server')}
            </Button>
          ) : (
            <Button onClick={() => setRoleDialogOpen(true)}>
              <Plus className='mr-1.5 size-4' />
              {t('Add MCP Role')}
            </Button>
          )}
        </SectionPageLayout.Actions>

        <SectionPageLayout.Content>
          <Tabs
            value={activeTab}
            onValueChange={(v) => setActiveTab(v as 'servers' | 'roles')}
            className='space-y-4'
          >
            <TabsList>
              <TabsTrigger value='servers'>
                {t('MCP Servers')}
              </TabsTrigger>
              <TabsTrigger value='roles'>
                {t('MCP Roles & Toolkits')}
              </TabsTrigger>
            </TabsList>

            <TabsContent value='servers' className='space-y-4'>
              <MCPServersTable />
            </TabsContent>

            <TabsContent value='roles' className='space-y-4'>
              <MCPRolesTable />
            </TabsContent>
          </Tabs>
        </SectionPageLayout.Content>
      </SectionPageLayout>

      <MCPServerDialog
        open={serverDialogOpen}
        onOpenChange={setServerDialogOpen}
      />
      <MCPRoleDialog
        open={roleDialogOpen}
        onOpenChange={setRoleDialogOpen}
      />
    </>
  )
}
