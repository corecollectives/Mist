import { FormModal } from "@/components/FormModal";
import { FullScreenLoading } from "@/shared/components";
import { Button } from "@/components/ui/button";
import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { toast } from "sonner";
import { useProjectStore } from "@/features/projects";

export const ProjectPage = () => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const params = useParams();
  const navigate = useNavigate();
  
  const { 
    currentProject: project, 
    fetchProjectById, 
    updateProject,
    deleteProject,
    isLoading: loading, 
    error 
  } = useProjectStore();

  const fetchProjectDetails = async () => {
    if (!params.projectId) return;
    await fetchProjectById(params.projectId);
  };

  const handleDeleteProject = async () => {
    if (!params.projectId) return;
    
    const success = await deleteProject(Number(params.projectId));
    if (success) {
      toast.success("Project deleted successfully");
      navigate("/projects");
    } else {
      toast.error("Failed to delete project");
    }
  };

  const handleUpdateProject = async (projectData: { name: string; description: string; tags: string[] }) => {
    if (!params.projectId) return;
    
    const success = await updateProject(Number(params.projectId), projectData);
    if (success) {
      toast.success("Project updated");
      setIsModalOpen(false);
    } else {
      toast.error("Failed to update project");
    }
  };

  useEffect(() => {
    fetchProjectDetails();
  }, [params.projectId]);

  if (loading) return <FullScreenLoading />;
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
      <div className="flex flex-col sm:flex-row justify-between sm:items-center gap-4 py-6  border-b border-border shrink-0">
        <div className="flex-1 min-w-0">
          <h1 className="text-xl sm:text-2xl font-bold text-foreground justify-between flex gap-4 flex-wrap">
            <span>{project.name}</span>
          </h1>
          <p className="text-sm sm:text-base text-muted-foreground mt-1 wrap-break-word">
            {project.description}
          </p>
          {project.projectMembers && (
            <div className="mt-3 flex flex-wrap gap-2 items-center">
              <span className="text-sm font-medium text-foreground">Members:</span>
              {project.projectMembers.map((member: any, index: number) => (
                <span
                  key={index}
                  className="bg-secondary text-secondary-foreground px-2 py-1 rounded-full text-xs sm:text-sm"
                >
                  {member.username}
                </span>
              ))}
              <Button variant="secondary" size="sm" className="ml-2 cursor-pointer">
                Manage Members
              </Button>
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
            onClick={handleDeleteProject}
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
