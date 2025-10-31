import { useEffect, useState } from "react"
import { toast } from "sonner"
import { useNavigate } from "react-router-dom"
import { useAuth } from "../context/AuthContext"
import Loading from "../components/Loading"
import type { Project } from "../lib/types"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Separator } from "@/components/ui/separator"
import { FormModal } from "@/components/FormModal"

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { MoreVertical } from "lucide-react"

export const ProjectsPage = () => {
  const { user } = useAuth()
  const [projects, setProjects] = useState<Project[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [editingProject, setEditingProject] = useState<Project | null>(null)
  const navigate = useNavigate()

  const fetchProjects = async () => {
    try {
      const response = await fetch("/api/projects/getAll")
      const data = await response.json()
      if (!data.success) throw new Error(data.error || "Failed to fetch projects")
      setProjects(data.data)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch projects")
      toast.error("Failed to fetch projects")
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchProjects()
    document.body.style.overflow = "hidden"
    return () => {
      document.body.style.overflow = ""
    }
  }, [])

  const handleCreateOrUpdateProject = async (projectData: { name: string; description: string; tags: string[] }) => {
    try {
      const isEditing = !!editingProject
      const url = isEditing ? `/api/projects/update?id=${editingProject.id}` : "/api/projects/create"
      const response = await fetch(url, {
        method: isEditing ? "PUT" : "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(projectData),
      })

      const data = await response.json()
      if (!data.success) throw new Error(data.error || "Failed to save project")

      toast.success(data.message || (isEditing ? "Project updated" : "Project created"))
      fetchProjects()
      setIsModalOpen(false)
      setEditingProject(null)
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to save project")
    }
  }

  const handleDeleteProject = async (projectId: number) => {
    if (!confirm("Are you sure you want to delete this project?")) return

    try {
      const response = await fetch(`/api/projects/delete?id=${projectId}`, {
        method: "DELETE",
      })
      const data = await response.json()
      if (!data.success) throw new Error(data.error || "Failed to delete project")
      setProjects(projects.filter((project) => project.id !== projectId))
      toast.success("Project deleted")
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to delete project")
    }
  }

  if (loading) return <Loading />

  return (
    <div className="flex flex-col h-screen bg-background overflow-hidden">
      {/* Header */}
      <div className="flex justify-between items-center py-6 border-b border-border flex-shrink-0">
        <div>
          <h1 className="text-2xl font-bold text-foreground">Projects</h1>
          <p className="text-muted-foreground mt-1">Create and manage your projects</p>
        </div>
        {user?.isAdmin && (
          <Button
            onClick={() => {
              setEditingProject(null)
              setIsModalOpen(true)
            }}
          >
            New Project
          </Button>
        )}
      </div>

      {/* Scrollable content */}
      <div className="flex-1 py-6 overflow-y-auto">
        {error ? (
          <Card className="border-destructive text-destructive">
            <CardContent className="p-4">{error}</CardContent>
          </Card>
        ) : projects.length === 0 ? (
          <Card className="text-center">
            <CardContent className="p-8">
              <p className="text-lg mb-4">No projects yet</p>
              {user?.isAdmin && (
                <Button
                  onClick={() => {
                    setEditingProject(null)
                    setIsModalOpen(true)
                  }}
                >
                  <svg
                    className="w-5 h-5 mr-2"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M12 4v16m8-8H4"
                    />
                  </svg>
                  Create first project
                </Button>
              )}
            </CardContent>
          </Card>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 pb-16">
            {projects.map((project) => (
              <Card
                key={project.id}
                className="relative transition-colors hover:border-primary"
              >
                <CardHeader className="flex flex-row items-start justify-between space-y-0">
                  <div
                    className="cursor-pointer"
                    onClick={() => navigate(`/projects/${project.id}`)}
                  >
                    <CardTitle>{project.name}</CardTitle>
                    <CardDescription>{project.description}</CardDescription>
                  </div>

                  {user?.isAdmin && (
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button variant="ghost" size="icon" className="h-8 w-8">
                          <MoreVertical className="w-4 h-4" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuItem
                          onClick={() => {
                            setEditingProject(project)
                            setIsModalOpen(true)
                          }}
                        >
                          Edit
                        </DropdownMenuItem>
                        <DropdownMenuItem
                          className="text-destructive"
                          onClick={() => handleDeleteProject(project.id)}
                        >
                          Delete
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  )}
                </CardHeader>

                <CardContent>
                  <div className="flex flex-wrap gap-2 mt-2">
                    {project?.tags?.map((tag) => (
                      <Badge key={tag} variant="secondary" className="capitalize">
                        {tag}
                      </Badge>
                    ))}
                  </div>

                  <Separator className="my-4" />

                  <div className="flex items-center justify-between text-sm text-muted-foreground">
                    <div className="flex items-center gap-2">
                      <div className="w-6 h-6 rounded-full bg-muted flex items-center justify-center">
                        <span className="text-xs font-medium text-foreground">
                          {project.owner?.username[0].toUpperCase()}
                        </span>
                      </div>
                      <span>{project.owner?.username}</span>
                    </div>
                    <span>{new Date(project.updatedAt || "").toLocaleDateString()}</span>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </div>

      {/* Modal for Create / Edit */}
      <FormModal
        isOpen={isModalOpen}
        onClose={() => {
          setIsModalOpen(false)
          setEditingProject(null)
        }}
        onSubmit={handleCreateOrUpdateProject}
        title={editingProject ? "Edit Project" : "Create New Project"}
        fields={[
          { name: "name", label: "Project Name", type: "text", defaultValue: editingProject?.name || "" },
          { name: "description", label: "Description", type: "textarea", defaultValue: editingProject?.description || "" },
          { name: "tags", label: "Tags", type: "tags", defaultValue: editingProject?.tags || [] },
        ]}
      />
    </div>
  )
}
