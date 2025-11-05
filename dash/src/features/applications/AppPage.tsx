import { FormModal } from "@/components/FormModal";
import { FullScreenLoading } from "@/shared/components";
import { Button } from "@/components/ui/button";
import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { toast } from "sonner";
import type { App } from "@/types/app";

export const AppPage = () => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [app, setApp] = useState<App | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const params = useParams();
  const navigate = useNavigate();

  const appId = parseInt(params.appId!);

  // Fetch app details
  const fetchAppDetails = async () => {
    try {
      setLoading(true);
      setError(null);

      const response = await fetch(`/api/apps/getById`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ appId }),
      });

      const data = await response.json();
      if (!data.success) throw new Error(data.error || "Failed to fetch app details");
      setApp(data.data);
    } catch (err) {
      const message = err instanceof Error ? err.message : "Failed to fetch app details";
      setError(message);
      toast.error(message);
    } finally {
      setLoading(false);
    }
  };

  const deleteAppHandler = async () => {
    try {
      const response = await fetch(`/api/apps/delete?id=${appId}`, {
        method: "DELETE",
        credentials: "include",
      });
      const data = await response.json();
      if (!data.success) throw new Error(data.error || "Failed to delete app");

      toast.success("App deleted successfully");
      navigate(-1); // Go back to the project page
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to delete app");
    }
  };

  const handleUpdateApp = async (appData: {
    name: string;
    description: string;
    git_repository?: string;
    git_branch?: string;
  }) => {
    try {
      const response = await fetch(`/api/apps/update?id=${appId}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(appData),
        credentials: "include",
      });
      const data = await response.json();
      if (!data.success) throw new Error(data.error || "Failed to update app");

      toast.success(data.message || "App updated successfully");
      await fetchAppDetails();
      setIsModalOpen(false);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to update app");
    }
  };

  useEffect(() => {
    fetchAppDetails();
  }, [params.appId]);

  if (loading) return <FullScreenLoading />;

  if (error)
    return (
      <div className="min-h-screen flex items-center justify-center bg-background">
        <div className="bg-destructive/10 border border-destructive text-destructive p-4 rounded-lg max-w-md text-center">
          {error}
        </div>
      </div>
    );

  if (!app) return null;

  return (
    <div className="flex flex-col min-h-screen bg-background">
      {/* Header */}
      <header className="border-b border-border py-6 flex flex-col sm:flex-row justify-between gap-4">
        <div>
          <h1 className="text-2xl font-semibold">{app.name}</h1>
          <p className="text-muted-foreground mt-1">{app.description}</p>

          <div className="mt-3 flex flex-wrap gap-2 items-center text-sm text-muted-foreground">
            {app.status && (
              <span className="bg-secondary text-secondary-foreground px-2 py-1 rounded-full">
                Status: {app.status}
              </span>
            )}
            {app.deployment_strategy && (
              <span className="bg-secondary text-secondary-foreground px-2 py-1 rounded-full">
                Strategy: {app.deployment_strategy}
              </span>
            )}
            {app.port && (
              <span className="bg-secondary text-secondary-foreground px-2 py-1 rounded-full">
                Port: {app.port}
              </span>
            )}
          </div>
        </div>

        <div className="flex flex-wrap gap-2 sm:flex-nowrap">
          <Button variant="outline" onClick={() => setIsModalOpen(true)}>
            Edit App
          </Button>
          <Button variant="destructive" onClick={deleteAppHandler}>
            Delete App
          </Button>
          <Button
            variant="secondary"
            onClick={() => navigate(`/projects/${app.project_id}`)}
          >
            Back to Project
          </Button>
        </div>
      </header>

      {/* Main Content */}
      <main className="flex-1 overflow-y-auto py-6 px-4 sm:px-6">
        <h2 className="text-lg font-semibold mb-4">App Details</h2>

        <div className="bg-card border border-border rounded-lg p-6 space-y-4">
          <div>
            <h3 className="text-sm text-muted-foreground">Git Repository</h3>
            <p className="text-foreground font-medium">
              {app.git_repository || "Not linked"}
            </p>
          </div>

          <div>
            <h3 className="text-sm text-muted-foreground">Branch</h3>
            <p className="text-foreground font-medium">
              {app.git_branch || "Not specified"}
            </p>
          </div>

          <div>
            <h3 className="text-sm text-muted-foreground">Created At</h3>
            <p className="text-foreground font-medium">
              {new Date(app.created_at).toLocaleString()}
            </p>
          </div>

          <div>
            <h3 className="text-sm text-muted-foreground">Created By</h3>
            <p className="text-foreground font-medium">{app.created_by}</p>
          </div>
        </div>
      </main>

      {/* Edit Modal */}
      <FormModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        title="Edit App"
        fields={[
          { label: "App Name", name: "name", type: "text", defaultValue: app.name },
          { label: "Description", name: "description", type: "textarea", defaultValue: app.description || "" },
          { label: "Git Repository", name: "git_repository", type: "text", defaultValue: app.git_repository || "" },
          { label: "Branch", name: "git_branch", type: "text", defaultValue: app.git_branch! },
        ]}
        onSubmit={(data) => handleUpdateApp(data as any)}
      />
    </div>
  );
};
