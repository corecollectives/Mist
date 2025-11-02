import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { FullScreenLoading } from '@/shared/components';
import { useGitStore } from './store';
import { GitHubCard, ProviderCard, CreateAppModal } from './components';
import { useAuth } from '@/context/AuthContext';

export default function GitPage() {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const { user } = useAuth();

  const {
    app,
    isInstalled,
    fetchApp,
    fetchRepositories,
    isLoading,
    error,
    clearError
  } = useGitStore();

  useEffect(() => {
    fetchApp();
    fetchRepositories();
  }, [fetchApp, fetchRepositories]);

  useEffect(() => {
    if (error) {
      toast.error(error);
      clearError();
    }
  }, [error, clearError]);

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
