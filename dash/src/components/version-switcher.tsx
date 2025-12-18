import { useEffect, useState } from "react"
import { SidebarMenu, SidebarMenuButton, SidebarMenuItem } from "@/components/ui/sidebar"
import { systemService } from "@/services"

export function VersionSwitcher() {
  const [version, setVersion] = useState<string>("...")

  useEffect(() => {
    const fetchVersion = async () => {
      try {
        const versionInfo = await systemService.getVersion()
        setVersion(versionInfo.version)
      } catch (error) {
        console.error("Failed to fetch version:", error)
        setVersion("unknown")
      }
    }

    fetchVersion()
  }, [])

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
            <span className="text-sm text-muted-foreground">v{version}</span>
          </div>
        </SidebarMenuButton>
      </SidebarMenuItem>
    </SidebarMenu>
  )
}
