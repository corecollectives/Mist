import * as React from "react"

import { VersionSwitcher } from "@/components/version-switcher"
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarRail,
} from "@/components/ui/sidebar"
import { Home, Settings, Users, GitBranch, Book, LifeBuoy, FolderGit2, User, Server } from "lucide-react"
import { Link, useLocation } from "react-router-dom"

const useNavData = () => {
  const location = useLocation()

  const data = {
    versions: ["1.0.1", "1.1.0-alpha", "2.0.0-beta1"],
    navMain: [
      {
        title: "Home",
        items: [
          {
            title: "Dashboard",
            url: "/",
            icon: Home,
            isActive: location.pathname === "/",
          },
          {
            title: "Projects",
            url: "/projects",
            icon: FolderGit2,
            isActive: location.pathname.startsWith("/projects"),
          },
          {
            title: "Deployments",
            url: "/deployments",
            icon: Server,
            isActive: location.pathname === "/deployments",
          },
        ],
      },
      {
        title: "Settings",
        items: [
          {
            title: "Users",
            url: "/users",
            icon: Users,
            isActive: location.pathname === "/users",
          },
          {
            title: "Mist",
            url: "/settings",
            icon: Settings,
            isActive: location.pathname === "/settings",
          },
          {
            title: "Profile",
            url: "/profile",
            icon: User,
            isActive: location.pathname === "/profile",
          },
          {
            title: "Git",
            url: "/git",
            icon: GitBranch,
            isActive: location.pathname === "/git",
          },
        ],
      },
      {
        title: "Extras",
        items: [
          {
            title: "Documentation",
            url: "/docs",
            icon: Book,
            isActive: location.pathname === "/docs",
          },
          {
            title: "Contribute",
            url: "/support",
            icon: LifeBuoy,
            isActive: location.pathname === "/support",
          },
        ],
      },
    ],
  }

  return { data }
}

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  const { data } = useNavData()
  return (
    <Sidebar {...props} collapsible="icon" variant="floating">
      <SidebarHeader>
        <VersionSwitcher
          defaultVersion={data.versions[0]}
        />
      </SidebarHeader>

      <SidebarContent>
        {data.navMain.map((group) => (
          <SidebarGroup key={group.title}>
            <SidebarGroupLabel>{group.title}</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {group.items.map((item) => {
                  const Icon = item.icon
                  return (
                    <SidebarMenuItem key={item.title}>
                      <SidebarMenuButton asChild isActive={item.isActive}>
                        <Link to={item.url} className="flex items-center gap-2">
                          <Icon className="h-4 w-4" />
                          <span>{item.title}</span>
                        </Link>
                      </SidebarMenuButton>
                    </SidebarMenuItem>
                  )
                })}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        ))}
      </SidebarContent>

      <SidebarRail />
    </Sidebar>
  )
}
