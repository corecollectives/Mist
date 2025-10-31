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
import { Home } from "lucide-react"
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
            icon: Home,
            isActive: location.pathname.startsWith("/projects"),
          },
          {
            title: "Deployments",
            url: "/deployments",
            icon: Home,
            isActive: location.pathname === "/deployments",
          }
        ]
      },
      {
        title: "Settings",
        items: [
          {
            title: "Users",
            url: "/users",
            icon: Home,
            isActive: location.pathname === "/users",
          },
          {
            title: "Mist",
            url: "/settings",
            icon: Home,
            isActive: location.pathname === "/settings",
          },
          {
            title: "Profile",
            url: "/profile",
            icon: Home,
            isActive: location.pathname === "/profile",
          },
          {
            title: "Git",
            url: "/git",
            icon: Home,
            isActive: location.pathname === "/git",
          }
        ]
      },
      {
        title: "Extras",
        items: [
          {
            title: "Documentation",
            url: "/docs",
            icon: Home,
            isActive: location.pathname === "/docs",
          },
          {
            title: "Contribute",
            url: "/support",
            icon: Home,
            isActive: location.pathname === "/support",
          }
        ]
      }

    ],
  }

  return { data };
}
export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  const { data } = useNavData();
  return (
    <Sidebar {...props} collapsible="icon" variant="floating">
      <SidebarHeader>
        <VersionSwitcher
          versions={data.versions}
          defaultVersion={data.versions[0]}
        />
      </SidebarHeader>
      <SidebarContent>
        {/* We create a SidebarGroup for each parent. */}
        {data.navMain.map((item) => (
          <SidebarGroup key={item.title}>
            <SidebarGroupLabel>{item.title}</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {item.items.map((item) => (
                  <SidebarMenuItem key={item.title}>
                    <SidebarMenuButton asChild isActive={item.isActive}>
                      <Link to={item.url}>
                        {item.title}</Link>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                ))}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        ))}
      </SidebarContent>
      <SidebarRail />
    </Sidebar>
  )
}
