import { useEffect, useState } from "react"
import { toast } from "react-toastify"
import { useNavigate } from "react-router-dom"
import { useAuth } from "../context/AuthContext"
import { CreateProjectModal } from "../components/CreateProjectModal"
import Loading from "../components/Loading"
import type { Project } from "../lib/types"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Separator } from "@/components/ui/separator"
import { FormModal } from "@/components/FormModal"

export const ProjectsPage = () => {
  const { user } = useAuth()
  const [projects, setProjects] = useState<Project[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [isModalOpen, setIsModalOpen] = useState(false)
  const navigate = useNavigate()

  const fetchProjects = async () => {
    try {
      const response = await fetch("/api/projects/getAll")
      const data = await response.json()

      if (!data.success) {
        toast.error(data.error || "Failed to fetch projects")
        throw new Error(data.error || "Failed to fetch projects")
      }

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
  }, [])

  const handleCreateProject = async (projectData: { name: string; description: string; tags: string[] }) => {
    try {
      const response = await fetch("/api/projects/create", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(projectData),
      })

      const data = await response.json()

      if (!data.success) {
        toast.error(data.error || "Failed to create project")
        throw new Error(data.error || "Failed to create project")
      }

      toast.success(data.message || "Project created successfully")
      fetchProjects()
      setIsModalOpen(false)
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to create project")
    }
  }

  if (loading) return <Loading />

  return (
    <div className="min-h-screen bg-background p-6">
      {/* Header */}
      <div className="flex justify-between items-center mb-6">
        <div>
          <h1 className="text-2xl font-bold text-foreground">Projects</h1>
          <p className="text-muted-foreground mt-1">Manage your projects and deployments</p>
        </div>
        {user?.isAdmin && (
          <Button onClick={() => setIsModalOpen(true)}>New Project</Button>
        )}
      </div>

      {/* Error */}
      {error ? (
        <Card className="border-destructive text-destructive">
          <CardContent className="p-4">{error}</CardContent>
        </Card>
      ) : projects.length === 0 ? (
        // Empty state
        <Card className="text-center">
          <CardContent className="p-8">
            <p className="text-lg mb-4">No projects yet</p>
            {user?.isAdmin && (
              <Button onClick={() => setIsModalOpen(true)}>
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
        // Projects grid
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {projects.map((project) => (
            <Card
              key={project.id}
              className="cursor-pointer transition-colors hover:border-primary"
              onClick={() => navigate(`/projects/${project.id}`)}
            >
              <CardHeader>
                <CardTitle>{project.name}</CardTitle>
                <CardDescription>{project.description}</CardDescription>
              </CardHeader>
              <CardContent>
                {/* Tags */}
                <div className="flex flex-wrap gap-2 mt-2">
                  {project?.tags?.map((tag) => (
                    <Badge key={tag} variant="secondary" className="capitalize">
                      {tag}
                    </Badge>
                  ))}
                </div>

                <Separator className="my-4" />

                {/* Footer: owner + date */}
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

      <FormModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleCreateProject}
        title="Create New Project"
        fields={[
          { name: "name", label: "Project Name", type: "text" },
          { name: "description", label: "Description", type: "textarea" },
          { name: "tags", label: "Tags", type: "tags" },
        ]}
      />
    </div>
  )
}
