import { useAuth } from '@/providers';
import { SystemUpdates } from '@/components/common/system-updates';

export const SettingsPage = () => {
  const { user } = useAuth();

  if (!user) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <p className="text-muted-foreground">Loading...</p>
      </div>
    );
  }

  // Only admins can access system settings
  if (user.role === 'user') {
    return (
      <div className="min-h-screen bg-background">
        <div className="py-6 border-b border-border">
          <div>
            <h1 className="text-2xl font-bold tracking-tight text-foreground">
              Settings
            </h1>
            <p className="text-muted-foreground mt-1">
              System settings and configuration
            </p>
          </div>
        </div>
        <div className="py-6 max-w-4xl">
          <div className="text-center py-12">
            <p className="text-muted-foreground">
              You need administrator privileges to access system settings.
            </p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background">
      <div className="py-6 border-b border-border">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-foreground">
            Settings
          </h1>
          <p className="text-muted-foreground mt-1">
            Manage system settings and configuration
          </p>
        </div>
      </div>

      {/* Content */}
      <div className="py-6 max-w-4xl">
        <SystemUpdates />
      </div>
    </div>
  );
};
