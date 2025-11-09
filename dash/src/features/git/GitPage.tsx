import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { FullScreenLoading } from '@/shared/components';
import { GitHubCard, ProviderCard, CreateAppModal } from './components';
import { useAuth } from '@/context/AuthContext';
import type { GitHubApp } from '@/types';

export default function GitPage() {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [app, setApp] = useState<GitHubApp | null>(null);
  const [isInstalled, setIsInstalled] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { user } = useAuth();

  const fetchApp = async () => {
    try {
      setIsLoading(true);
      setError(null);
      const response = await fetch("/api/github/app", {
        credentials: 'include'
      });
      const data = await response.json();
      console.log("Fetched GitHub App:", data.data.app);
      if (data.success) {
        setApp(data.data.app);
        setIsInstalled(data.data.isInstalled);
      } else {
        setError(data.error || "Failed to load GitHub App info");
      }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to load GitHub App info";
      setError(errorMessage);
      console.error(err);
    } finally {
      setIsLoading(false);
    }
  };

  const fetchRepositories = async () => {
    try {
      const response = await fetch("/api/github/repositories", {
        credentials: 'include'
      });
      const data = await response.json();
      console.log("Fetched repos:", data);
    } catch (err) {
      console.error("Failed to fetch repos:", err);
    }
  };

  useEffect(() => {
    fetchApp();
    fetchRepositories();
  }, []);

  useEffect(() => {
    if (error) {
      toast.error(error);
      setError(null);
    }
  }, [error]);

  const handleCreateApp = async () => {
    try {
      window.open('/api/github/app/create', '_blank');
    } catch (err) {
    }
  };

  if (isLoading && !app) {
    return <FullScreenLoading />;
  }

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <div className="flex items-center justify-between py-6 border-b border-border shrink-0">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-foreground">
            Git Integrations
          </h1>
          <p className="text-muted-foreground mt-1">
            Connect your Git providers to enable automatic deployments.
          </p>
        </div>
      </div>

      {/* Error */}
      {error && (
        <div className="mt-4">
          <Alert variant="destructive">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        </div>
      )}

      {/* Main Grid */}
      <div className="grid gap-4 py-6 md:grid-cols-2 lg:grid-cols-3">
        {/* GitHub Card */}
        <GitHubCard
          app={app}
          isInstalled={isInstalled}
          user={user}
          onCreateApp={() => setIsModalOpen(true)}
        />

        {/* Other Providers */}
        <ProviderCard name="GitLab" icon="Gitlab" />
        <ProviderCard name="Bitbucket" icon="GitFork" />
        <ProviderCard name="Gitea" icon="GitMerge" />
      </div>

      {/* Create App Modal */}
      <CreateAppModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onConfirm={handleCreateApp}
      />
    </div>
  );
}
