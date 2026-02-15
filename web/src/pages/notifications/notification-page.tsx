import { PageHeader } from '@/components/layout/page-header'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Bell, FileText } from 'lucide-react'
import { InternalNotificationList } from './internal-notification-list'
import { TemplateList } from './template-list'

export default function NotificationPage() {
  return (
    <div>
      <PageHeader
        title="Notificaciones"
        description="Centro de notificaciones y gestión de plantillas"
      />

      <Tabs defaultValue="internal" className="space-y-4">
        <TabsList>
          <TabsTrigger value="internal" className="flex items-center gap-2">
            <Bell className="h-4 w-4" />
            Mis Notificaciones
          </TabsTrigger>
          <TabsTrigger value="templates" className="flex items-center gap-2">
            <FileText className="h-4 w-4" />
            Plantillas
          </TabsTrigger>
        </TabsList>

        <TabsContent value="internal">
          <InternalNotificationList />
        </TabsContent>

        <TabsContent value="templates">
          <TemplateList />
        </TabsContent>
      </Tabs>
    </div>
  )
}
