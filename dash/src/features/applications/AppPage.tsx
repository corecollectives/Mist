import { FormModal } from "@/components/FormModal";
import { FullScreenLoading } from "@/shared/components";
import { Button } from "@/components/ui/button";
import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { toast } from "sonner";
import type { App } from "@/types/app";

interface GitHubRepo {
  id: number;
  name: string;
  full_name: string;
  private: boolean;
  html_url: string;
}

interface GitHubBranch {
  name: string;
}

export const AppPage = () => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [app, setApp] = useState<App | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [repositories, setRepositories] = useState<GitHubRepo[]>([]);
  const [branches, setBranches] = useState<GitHubBranch[]>([]);
  const [selectedRepo, setSelectedRepo] = useState<string>("");

  console.log(selectedRepo)
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
      setSelectedRepo(data.data.git_repository || "");
    } catch (err) {
      const message = err instanceof Error ? err.message : "Failed to fetch app details";
      setError(message);
      toast.error(message);
    } finally {
      setLoading(false);
    }
  };

  // Fetch repositories from GitHub
  const fetchRepositories = async () => {
    try {
      const res = await fetch(`/api/github/repositories`, { credentials: "include" });
      const data = await res.json();
      setRepositories(data);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to fetch repositories");
    }
  };

  // Fetch branches for selected repository
  const fetchBranches = async (repoFullName: string) => {
    if (!repoFullName) return;
    try {
      const res = await fetch(`/api/github/branches?repo=${repoFullName}`, { credentials: "include" });
      const data = await res.json();
      if (!data.success) throw new Error(data.error || "Failed to fetch branches");
      setBranches(data.data);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to fetch branches");
    }
  };

  const updateApp = async () => {
    try {
      const res = await fetch(`/api/apps/update`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ appId, gitRepository: selectedRepo }),
      })
      const data = await res.json();
      if (!data.success) throw new Error(data.error || "Failed to update app");
      toast.success("App updated successfully");
      await fetchAppDetails();
    }
    catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to update app");
    }
  }

  const deleteAppHandler = async () => {
    try {
      const response = await fetch(`/api/apps/delete?id=${appId}`, {
        method: "DELETE",
        credentials: "include",
      });
      const data = await response.json();
      if (!data.success) throw new Error(data.error || "Failed to delete app");

      toast.success("App deleted successfully");
      navigate(-1);
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
    fetchRepositories();
  }, [params.appId]);

  useEffect(() => {
    if (selectedRepo) fetchBranches(selectedRepo);
  }, [selectedRepo]);

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
          <Button variant="secondary" onClick={() => navigate(`/projects/${app.project_id}`)}>
            Back to Project
          </Button>
        </div>
      </header>

      {/* App Info */}
      <main className="flex-1 overflow-y-auto py-6 ">

        <div >
          <label className="text-sm text-muted-foreground">Select repo</label>
          <select
            value={selectedRepo}
            onChange={(e) => setSelectedRepo(e.target.value)}
            className="w-full bg-background border rounded-md mt-1 px-3 py-2"
          >
            <option value="">Select a repository</option>
            {repositories.map((repo) => (
              <option key={repo.id} value={repo.full_name}>
                {repo.full_name}
              </option>
            ))}
          </select>
        </div>
        <Button
          onClick={() => updateApp()}
        >
          Save
        </Button>
        <h2 className="text-lg font-semibold mb-4">App Details</h2>

        <div className="bg-card border border-border rounded-lg p-6 space-y-4">
          <div>
            <h3 className="text-sm text-muted-foreground">Git Repository</h3>
            <p className="text-foreground font-medium">{app.git_repository || "Not linked"}</p>
          </div>
          <div>
            <h3 className="text-sm text-muted-foreground">Branch</h3>
            <p className="text-foreground font-medium">{app.git_branch || "Not specified"}</p>
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
          {
            label: "Git Repository",
            name: "git_repository",
            type: "select",
            options: repositories.map((repo) => ({
              label: repo.full_name,
              value: repo.full_name,
            })),
            defaultValue: app.git_repository || "",
            // onChange: (value: string) => setSelectedRepo(valu),
          },
          {
            label: "Branch",
            name: "git_branch",
            type: "select",
            options: branches.map((b) => ({ label: b.name, value: b.name })),
            defaultValue: app.git_branch || "",
          },
        ]}
        onSubmit={(data) => handleUpdateApp(data as any)}
      />
    </div>
  );
};
