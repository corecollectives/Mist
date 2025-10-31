import { FormModal } from "@/components/FormModal";
import Loading from "@/components/Loading";
import { Button } from "@/components/ui/button";
import type { Project } from "@/lib/types";
import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { toast } from "sonner";

export const ProjectPage = () => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const params = useParams();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [project, setProject] = useState<Project | null>(null);
  const navigate = useNavigate();

  const fetchProjectDetails = async () => {
    try {
      const response = await fetch(`/api/projects/getFromId?id=${params.projectId}`, {
        method: "GET",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
      });
      const data = await response.json();
      if (!data.success) throw new Error(data.error || "Failed to fetch project details");
      setProject(data.data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch project details");
      toast.error("Failed to fetch project details");
    } finally {
      setLoading(false);
    }
  };

  const deleteProject = async () => {
    try {
      const response = await fetch(`/api/projects/delete?id=${params.projectId}`, {
        method: "DELETE",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
      });
      const data = await response.json();
      if (!data.success) throw new Error(data.error || "Failed to delete project");
      toast.success("Project deleted successfully");
      navigate("/projects");
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to delete project");
    } finally {
      setLoading(false);
    }
  };

  const handleUpdateProject = async (projectData: { name: string; description: string; tags: string[] }) => {
    try {
      const response = await fetch(`/api/projects/update?id=${params.projectId}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(projectData),
        credentials: "include",
      });
      const data = await response.json();
      if (!data.success) throw new Error(data.error || "Failed to update project");

      toast.success(data.message || "Project updated");
      fetchProjectDetails();
      setIsModalOpen(false);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to update project");
    }
  };

  useEffect(() => {
    fetchProjectDetails();
  }, []);

  if (loading) return <Loading />;
  if (error)
    return (
      <div className="min-h-screen bg-[#0D1117] p-6 flex items-center justify-center">
        <div className="bg-[#F8514933] border border-[#F85149] text-[#F85149] p-4 rounded-lg max-w-md text-center">
          {error}
        </div>
      </div>
    );

  if (!project) return null;

  return (
    <div className="flex flex-col h-screen bg-background overflow-hidden">
      {/* Header */}
      <div className="flex flex-col sm:flex-row justify-between sm:items-center gap-4 py-6  border-b border-border flex-shrink-0">
        <div className="flex-1 min-w-0">
          <h1 className="text-xl sm:text-2xl font-bold text-foreground break-words">
            {project.name}
          </h1>
          <p className="text-sm sm:text-base text-muted-foreground mt-1 break-words">
            {project.description}
          </p>
          {project?.tags && (
            <div className="mt-3 flex flex-wrap gap-2">
              {project?.tags?.map((tag, index) => (
                <span
                  key={index}
                  className="bg-accent text-accent-foreground px-2 py-1 rounded-full text-xs sm:text-sm"
                >
                  {tag}
                </span>
              ))}
            </div>
          )}
        </div>

        <div className="flex flex-wrap gap-2 sm:flex-nowrap sm:justify-end">
          <Button
            variant="outline"
            onClick={() => setIsModalOpen(true)}
            className="w-full sm:w-auto transition-colors"
          >
            Edit Project
          </Button>
          <Button
            variant="destructive"
            onClick={deleteProject}
            className="w-full sm:w-auto transition-colors"
          >
            Delete Project
          </Button>
        </div>
      </div>

      <FormModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        title="Edit Project"
        fields={[
          { label: "Name", name: "name", type: "text", defaultValue: project.name },
          { label: "Description", name: "description", type: "textarea", defaultValue: project.description },
          { name: "tags", label: "Tags", type: "tags", defaultValue: project.tags || [] },
        ]}
        onSubmit={(data) =>
          handleUpdateProject(data as { name: string; description: string; tags: string[] })
        }
      />

      {/* Scrollable content */}
      <div className="flex-1 py-6 px-4 sm:px-6 overflow-y-auto">
        {/* Future project details and apps go here */}
        <div className="text-center text-muted-foreground text-sm sm:text-base">
          No additional project details available yet.
        </div>
      </div>
    </div>
  );
};
