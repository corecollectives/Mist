import { useState, useEffect } from 'react';
import { useAuth } from '@/providers';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { AlertCircle, CheckCircle2, Globe, ShieldAlert } from 'lucide-react';
import { settingsService } from '@/services';
import { toast } from 'sonner';
import { useNavigate } from 'react-router-dom';

export const SettingsPage = () => {
  const { user } = useAuth();
  const navigate = useNavigate();
  const [isUpdatingSystemSettings, setIsUpdatingSystemSettings] = useState(false);

  // System settings state
  const [wildcardDomain, setWildcardDomain] = useState('');
  const [mistAppName, setMistAppName] = useState('mist');
  const [systemSettingsError, setSystemSettingsError] = useState('');
  const [isLoadingSystemSettings, setIsLoadingSystemSettings] = useState(true);

  // Check if user is owner
  useEffect(() => {
    if (!user) {
      navigate('/');
      return;
    }

    if (user.role !== 'owner') {
      toast.error('Only owners can access system settings');
      navigate('/');
      return;
    }

    loadSystemSettings();
  }, [user, navigate]);

  const loadSystemSettings = async () => {
    try {
      const settings = await settingsService.getSystemSettings();
      setWildcardDomain(settings.wildcardDomain || '');
      setMistAppName(settings.mistAppName);
    } catch (error) {
      console.error('Failed to load system settings:', error);
      toast.error('Failed to load system settings');
    } finally {
      setIsLoadingSystemSettings(false);
    }
  };

  const handleUpdateSystemSettings = async (e: React.FormEvent) => {
    e.preventDefault();
    setSystemSettingsError('');

    if (!mistAppName.trim()) {
      setSystemSettingsError('Mist app name is required');
      return;
    }

    setIsUpdatingSystemSettings(true);

    try {
      const settings = await settingsService.updateSystemSettings(
        wildcardDomain.trim() || null,
        mistAppName.trim()
      );
      setWildcardDomain(settings.wildcardDomain || '');
      setMistAppName(settings.mistAppName);
      toast.success('System settings updated successfully');
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to update system settings';
      setSystemSettingsError(message);
      toast.error(message);
    } finally {
      setIsUpdatingSystemSettings(false);
    }
  };

  if (!user) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <p className="text-muted-foreground">Loading...</p>
      </div>
    );
  }

  if (user.role !== 'owner') {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <Card className="max-w-md">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <ShieldAlert className="h-5 w-5 text-destructive" />
              Access Denied
            </CardTitle>
            <CardDescription>
              Only owners can access system settings
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Button onClick={() => navigate('/')}>Go to Dashboard</Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background">
      <div className="py-6 border-b border-border">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-foreground">
            System Settings
          </h1>
          <p className="text-muted-foreground mt-1">
            Configure system-wide settings for your Mist instance
          </p>
        </div>
      </div>

      {/* Content */}
      <div className="py-6 max-w-4xl">
        {/* System Settings */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Globe className="h-5 w-5 text-primary" />
              Wildcard Domain Configuration
            </CardTitle>
            <CardDescription>
              Configure wildcard domain for automatic app domain generation
            </CardDescription>
          </CardHeader>
          <CardContent>
            {isLoadingSystemSettings ? (
              <div className="flex items-center justify-center py-8">
                <p className="text-muted-foreground">Loading settings...</p>
              </div>
            ) : (
              <form onSubmit={handleUpdateSystemSettings} className="space-y-6">
                {systemSettingsError && (
                  <Alert variant="destructive">
                    <AlertCircle className="h-4 w-4" />
                    <AlertDescription>{systemSettingsError}</AlertDescription>
                  </Alert>
                )}

                <div className="space-y-2">
                  <Label htmlFor="wildcardDomain">Wildcard Domain</Label>
                  <Input
                    id="wildcardDomain"
                    type="text"
                    value={wildcardDomain}
                    onChange={(e) => setWildcardDomain(e.target.value)}
                    placeholder="*.exam.ple or exam.ple"
                    disabled={isUpdatingSystemSettings}
                  />
                  <p className="text-sm text-muted-foreground">
                    When configured, apps will automatically get domains like <code className="bg-muted px-1 py-0.5 rounded">project-app.exam.ple</code>
                  </p>
                  <div className="mt-3 p-3 bg-muted rounded-md">
                    <p className="text-sm font-medium mb-2">Example:</p>
                    <ul className="text-sm text-muted-foreground space-y-1 list-disc list-inside">
                      <li>Wildcard domain: <code className="bg-background px-1 py-0.5 rounded">exam.ple</code></li>
                      <li>Project name: <code className="bg-background px-1 py-0.5 rounded">crux</code></li>
                      <li>App name: <code className="bg-background px-1 py-0.5 rounded">main</code></li>
                      <li>Generated domain: <code className="bg-background px-1 py-0.5 rounded">crux-main.exam.ple</code></li>
                    </ul>
                  </div>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="mistAppName">Mist App Name</Label>
                  <Input
                    id="mistAppName"
                    type="text"
                    value={mistAppName}
                    onChange={(e) => setMistAppName(e.target.value)}
                    placeholder="mist"
                    disabled={isUpdatingSystemSettings}
                  />
                  <p className="text-sm text-muted-foreground">
                    Subdomain name for the Mist dashboard. With wildcard domain <code className="bg-muted px-1 py-0.5 rounded">exam.ple</code> and name <code className="bg-muted px-1 py-0.5 rounded">mist</code>, 
                    Mist will be available at <code className="bg-muted px-1 py-0.5 rounded">mist.exam.ple</code>
                  </p>
                </div>

                <div className="flex items-center gap-2 pt-4">
                  <Button
                    type="submit"
                    disabled={isUpdatingSystemSettings || !mistAppName.trim()}
                  >
                    {isUpdatingSystemSettings ? (
                      <>
                        <span className="animate-spin mr-2">⏳</span>
                        Updating...
                      </>
                    ) : (
                      <>
                        <CheckCircle2 className="h-4 w-4 mr-2" />
                        Update System Settings
                      </>
                    )}
                  </Button>
                  <Button
                    type="button"
                    variant="outline"
                    onClick={loadSystemSettings}
                    disabled={isUpdatingSystemSettings}
                  >
                    Reset
                  </Button>
                </div>
              </form>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
};
