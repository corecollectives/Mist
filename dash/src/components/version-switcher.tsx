import { SidebarMenu, SidebarMenuButton, SidebarMenuItem } from "@/components/ui/sidebar"

export function VersionSwitcher({
  defaultVersion,
}: {
  defaultVersion: string
}) {
  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <SidebarMenuButton
          size="lg"
          className="cursor-default hover:bg-transparent focus-visible:ring-0 focus-visible:ring-offset-0"
        >
          <img src="/mist.png" alt="Mist Icon" className="size-10" />
          <div className="flex flex-col gap-0.5 leading-none">
            <span className="font-medium">Mist</span>
            <span className="text-sm text-muted-foreground">v{defaultVersion}</span>
          </div>
        </SidebarMenuButton>
      </SidebarMenuItem>
    </SidebarMenu>
  )
}
