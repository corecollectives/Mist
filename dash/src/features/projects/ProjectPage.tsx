import { FormModal } from "@/components/FormModal";
import { FullScreenLoading } from "@/components/common";
import { Button } from "@/components/ui/button";
import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { toast } from "sonner";
import type { Project } from "@/types";
import type { App, CreateAppRequest } from "@/types/app";
import { AppCard } from "./components/AppCard";
import { CreateAppModal } from "./components/CreateAppModal";

export const ProjectPage = () => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isAddNewAppModalOpen, setIsAddNewAppModalOpen] = useState(false);
  const [project, setProject] = useState<Project | null>(null);
  const [apps, setApps] = useState<App[]>([]);
  const [loading, setLoading] = useState(true);
  const [fetchingApps, setFetchingApps] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const params = useParams();
  const navigate = useNavigate();

  const projectId = parseInt(params.projectId!);

  // Fetch project details
  const fetchProjectDetails = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await fetch(`/api/projects/getFromId?id=${projectId}`, {
        credentials: "include",
      });
      const data = await response.json();
      if (!data.success) throw new Error(data.error || "Failed to fetch project details");
      setProject(data.data);
    } catch (err) {
      const message = err instanceof Error ? err.message : "Failed to fetch project details";
      setError(message);
      toast.error(message);
    } finally {
      setLoading(false);
    }
  };

  // Fetch all apps in the project
  const getApps = async () => {
    try {
      setFetchingApps(true);
      const response = await fetch(`/api/apps/getByProjectId`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ projectId }),
      });
      const data = await response.json();
      if (!data.success) throw new Error(data.error || "Failed to fetch apps");
      setApps(data.data);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to fetch apps");
    } finally {
      setFetchingApps(false);
    }
  };

  const createNewApp = async (appData: CreateAppRequest) => {
    try {
      const response = await fetch(`/api/apps/create`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(appData),
        credentials: "include",
      });
      const data = await response.json();
      if (!data.success) throw new Error(data.error || "Failed to create app");

      toast.success("App created successfully");
      setIsAddNewAppModalOpen(false);
      await getApps(); // Refresh app list after creating
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to create app");
    }
  };

  const deleteProjectHandler = async () => {
    try {
      const response = await fetch(`/api/projects/delete?id=${projectId}`, {
        method: "DELETE",
        credentials: "include",
      });
      const data = await response.json();
      if (!data.success) throw new Error(data.error || "Failed to delete project");

      toast.success("Project deleted successfully");
      navigate("/projects");
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to delete project");
    }
  };

  const handleUpdateProject = async (projectData: {
    name: string;
    description: string;
    tags: string[];
  }) => {
    try {
      const response = await fetch(`/api/projects/update?id=${projectId}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(projectData),
        credentials: "include",
      });
      const data = await response.json();
      if (!data.success) throw new Error(data.error || "Failed to update project");

      toast.success(data.message || "Project updated");
      await fetchProjectDetails();
      setIsModalOpen(false);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to update project");
    }
  };

  useEffect(() => {
    fetchProjectDetails();
    getApps();
  }, [params.projectId]);

  if (loading) return <FullScreenLoading />;

  if (error)
    return (
      <div className="min-h-screen flex items-center justify-center bg-background">
        <div className="bg-destructive/10 border border-destructive text-destructive p-4 rounded-lg max-w-md text-center">
          {error}
        </div>
      </div>
    );

  if (!project) return null;

  return (
    <div className="flex flex-col min-h-screen bg-background">
      {/* Header */}
      <header className="border-b border-border py-6 flex flex-col sm:flex-row justify-between gap-4">
        <div>
          <h1 className="text-2xl font-semibold">{project.name}</h1>
          <p className="text-muted-foreground mt-1">{project.description}</p>
          {project.projectMembers && (
            <div className="mt-3 flex flex-wrap gap-2 items-center">
              <span className="text-sm font-medium text-foreground">Members:</span>
              {project.projectMembers.map((member: any) => (
                <span
                  key={member.username}
                  className="bg-secondary text-secondary-foreground px-2 py-1 rounded-full text-xs"
                >
                  {member.username}
                </span>
              ))}
              <Button variant="secondary" size="sm" className="ml-2">
                Manage Members
              </Button>
            </div>
          )}
        </div>

        <div className="flex flex-wrap gap-2 sm:flex-nowrap">
          <Button variant="outline" onClick={() => setIsModalOpen(true)}>
            Edit Project
          </Button>
          <Button variant="destructive" onClick={deleteProjectHandler}>
            Delete Project
          </Button>
          <Button variant="secondary" onClick={() => setIsAddNewAppModalOpen(true)}>
            Create App
          </Button>
        </div>
      </header>

      {/* Apps Section */}
      <main className="flex-1 overflow-y-auto py-6">

        {fetchingApps ? (
          <div className="text-muted-foreground text-center py-10">Loading apps...</div>
        ) : apps && apps.length === 0 ? (
          <div className="text-muted-foreground text-center py-10">
            No apps found for this project. Create one to get started.
          </div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            {apps && apps.map((app) => (
              <AppCard
                key={app.id}
                app={app}
                onClick={() => navigate(`/projects/${app.projectId}/apps/${app.id}`)}
              />
            ))}
          </div>
        )}
      </main>

      {/* Modals */}
      <FormModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        title="Edit Project"
        fields={[
          { label: "Name", name: "name", type: "text", defaultValue: project.name },
          { label: "Description", name: "description", type: "textarea", defaultValue: project.description },
          { name: "tags", label: "Tags", type: "tags", defaultValue: project.tags || [] },
        ]}
        onSubmit={(data) => handleUpdateProject(data as any)}
      />

      <CreateAppModal
        isOpen={isAddNewAppModalOpen}
        onClose={() => setIsAddNewAppModalOpen(false)}
        projectId={projectId}
        onSubmit={createNewApp}
      />
    </div>
  );
};
