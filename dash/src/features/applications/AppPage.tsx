import { FormModal } from "@/components/FormModal";
import { FullScreenLoading } from "@/components/common";
import { Button } from "@/components/ui/button";
import { useEffect, useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { toast } from "sonner";
import type { App } from "@/types/app";
import { TabsList, Tabs, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import { AppInfo, GitProviderTab, EnvironmentVariables, Domains, AppSettings, ContainerControls, LiveLogsViewer, AppStats } from "@/components/applications";
import { DeploymentsTab } from "@/components/deployments";


export const AppPage = () => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [app, setApp] = useState<App | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [latestCommit, setLatestCommit] = useState();
  const [previewUrl, setPreviewUrl] = useState<string>("");

  const params = useParams();
  const navigate = useNavigate();

  const appId = useMemo(() => Number(params.appId), [params.appId]);
  const projectId = parseInt(params.projectId!);

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


  const fetchLatestCommit = async () => {
    try {
      const res = await fetch(`/api/apps/getLatestCommit`,
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          credentials: "include",
          body: JSON.stringify({ appID: appId, projectID: projectId }),
        }
      );
      const data = await res.json();
      if (!data.success) throw new Error(data.error || "Failed to fetch latest commit");
      setLatestCommit(data.data);
    }
    catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to fetch latest commit");
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
  }) => {
    try {
      const response = await fetch(`/api/apps/update`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ appId: appId, ...appData }),
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
    fetchLatestCommit()
    fetchAppDetails()
  }, [params.appId]);

  // Fetch preview URL when app changes
  useEffect(() => {
    const fetchPreviewUrl = async () => {
      if (!app || app.status !== "running") return
      
      try {
        const response = await fetch(`/api/apps/getPreviewUrl`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          credentials: "include",
          body: JSON.stringify({ appId: app.id }),
        })
        const data = await response.json()
        if (data.success) {
          setPreviewUrl(data.data.url)
        }
      } catch (err) {
        console.error("Failed to fetch preview URL:", err)
      }
    }

    fetchPreviewUrl()
  }, [app]);




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
        </div>

        <div className="flex flex-wrap gap-2 sm:flex-nowrap">
          <Button variant="outline" onClick={() => setIsModalOpen(true)}>
            Edit App
          </Button>
          <Button variant="destructive" onClick={deleteAppHandler}>
            Delete App
          </Button>
        </div>
      </header>

      {/* App Info */}
      <main className="flex-1 overflow-y-auto py-6">
        <Tabs defaultValue="info" className="w-full">
          <TabsList className="grid w-full grid-cols-6 mb-6">
            <TabsTrigger value="info">Info</TabsTrigger>
            <TabsTrigger value="git">Git</TabsTrigger>
            <TabsTrigger value="environment">Environment</TabsTrigger>
            <TabsTrigger value="deployments">Deployments</TabsTrigger>
            <TabsTrigger value="logs">Logs</TabsTrigger>
            <TabsTrigger value="settings">Settings</TabsTrigger>
          </TabsList>

          {/* ✅ INFO TAB */}
          <TabsContent value="info" className="space-y-6">
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
              <div className="lg:col-span-2">
                <AppInfo app={app} latestCommit={latestCommit} />
              </div>
              <div>
                <AppStats appId={app.id} appStatus={app.status} previewUrl={previewUrl} />
              </div>
            </div>
          </TabsContent>

          <TabsContent value="git" className="space-y-6">
            <GitProviderTab app={app} />
          </TabsContent>

          {/* ✅ ENVIRONMENT TAB */}
          <TabsContent value="environment" className="space-y-6">
            <EnvironmentVariables appId={app.id} />
          </TabsContent>

          {/* ✅ DEPLOYMENTS TAB */}
          <TabsContent value="deployments">
            <DeploymentsTab appId={app.id} />
          </TabsContent>

          <TabsContent value="logs" className="h-full">
            <LiveLogsViewer appId={app.id} enabled={true} />
          </TabsContent>

          <TabsContent value="settings" className="space-y-6">
            <ContainerControls appId={app.id} onStatusChange={fetchAppDetails} />
            <AppSettings app={app} onUpdate={fetchAppDetails} />
            <Domains appId={app.id} />
          </TabsContent>
        </Tabs>
      </main>

      {/* Edit Modal */}
      <FormModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        title="Edit App"
        fields={[
          { label: "App Name", name: "name", type: "text", defaultValue: app.name },
          { label: "Description", name: "description", type: "textarea", defaultValue: app.description || "" },
        ]}
        onSubmit={(data) => handleUpdateApp(data as any)}
      />
    </div>
  );
};
