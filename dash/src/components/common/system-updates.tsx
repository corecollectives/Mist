import { useState, useEffect } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { 
  RefreshCw, 
  Download, 
  CheckCircle2, 
  XCircle, 
  AlertCircle, 
  ArrowUpCircle,
  Clock,
  ChevronDown,
  ChevronUp
} from 'lucide-react';
import { toast } from 'sonner';
import { systemService } from '@/services/system.service';
import type { SystemInfo, UpdateCheck, UpdateHistory } from '@/types/system';

export const SystemUpdates = () => {
  const [systemInfo, setSystemInfo] = useState<SystemInfo | null>(null);
  const [updateCheck, setUpdateCheck] = useState<UpdateCheck | null>(null);
  const [updateHistory, setUpdateHistory] = useState<UpdateHistory[]>([]);
  const [checking, setChecking] = useState(false);
  const [updating, setUpdating] = useState(false);
  const [showHistory, setShowHistory] = useState(false);
  const [showChangelog, setShowChangelog] = useState(false);

  useEffect(() => {
    loadSystemInfo();
    loadUpdateHistory();
  }, []);

  const loadSystemInfo = async () => {
    try {
      const info = await systemService.getVersion();
      setSystemInfo(info);
    } catch (err) {
      console.error('Failed to load system info:', err);
    }
  };

  const loadUpdateHistory = async () => {
    try {
      const history = await systemService.getUpdateHistory();
      setUpdateHistory(history);
    } catch (err) {
      console.error('Failed to load update history:', err);
    }
  };

  const handleCheckForUpdates = async () => {
    setChecking(true);
    try {
      const check = await systemService.checkForUpdates();
      setUpdateCheck(check);
      
      if (check.hasUpdate) {
        toast.success(`New version available: ${check.latestVersion}`);
      } else {
        toast.info('You are running the latest version');
      }
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to check for updates';
      toast.error(message);
    } finally {
      setChecking(false);
    }
  };

  const handleTriggerUpdate = async () => {
    if (!updateCheck?.latestVersion) return;

    const confirmed = confirm(
      `Are you sure you want to update to version ${updateCheck.latestVersion}?\n\n` +
      'The system will be briefly unavailable during the update (30-60 seconds).'
    );

    if (!confirmed) return;

    setUpdating(true);
    try {
      await systemService.triggerUpdate({
        version: updateCheck.latestVersion,
        branch: 'main'
      });
      
      toast.success('Update started! The system will restart shortly.');
      
      // Poll for completion
      setTimeout(() => {
        window.location.reload();
      }, 10000);
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to start update';
      toast.error(message);
      setUpdating(false);
    }
  };

  const getStatusBadge = (status: string) => {
    const config = {
      success: { icon: CheckCircle2, color: 'bg-green-500', label: 'Success' },
      failed: { icon: XCircle, color: 'bg-red-500', label: 'Failed' },
      pending: { icon: Clock, color: 'bg-yellow-500', label: 'Pending' },
      downloading: { icon: Download, color: 'bg-blue-500', label: 'Downloading' },
      installing: { icon: RefreshCw, color: 'bg-blue-500', label: 'Installing' },
    };

    const { icon: Icon, color, label } = config[status as keyof typeof config] || config.pending;
    
    return (
      <Badge className={`${color} text-white`}>
        <Icon className="h-3 w-3 mr-1" />
        {label}
      </Badge>
    );
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const parseMarkdown = (text: string) => {
    // Simple markdown parsing for changelog
    return text
      .replace(/^### (.+)$/gm, '<h3 class="text-lg font-semibold mt-4 mb-2">$1</h3>')
      .replace(/^## (.+)$/gm, '<h2 class="text-xl font-bold mt-6 mb-3">$1</h2>')
      .replace(/^# (.+)$/gm, '<h1 class="text-2xl font-bold mt-8 mb-4">$1</h1>')
      .replace(/^\* (.+)$/gm, '<li class="ml-4">$1</li>')
      .replace(/^- (.+)$/gm, '<li class="ml-4">$1</li>')
      .replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>')
      .replace(/\*(.+?)\*/g, '<em>$1</em>')
      .replace(/`(.+?)`/g, '<code class="bg-muted px-1 py-0.5 rounded text-sm">$1</code>')
      .replace(/\n\n/g, '</p><p class="mb-2">');
  };

  return (
    <div className="space-y-6">
      {/* Current Version */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle className="flex items-center gap-2">
                <ArrowUpCircle className="h-5 w-5 text-primary" />
                System Version
              </CardTitle>
              <CardDescription>
                Current version and update status
              </CardDescription>
            </div>
            <Button 
              onClick={handleCheckForUpdates}
              disabled={checking}
              variant="outline"
              size="sm"
            >
              {checking ? (
                <>
                  <RefreshCw className="h-4 w-4 mr-2 animate-spin" />
                  Checking...
                </>
              ) : (
                <>
                  <RefreshCw className="h-4 w-4 mr-2" />
                  Check for Updates
                </>
              )}
            </Button>
          </div>
        </CardHeader>
        <CardContent className="space-y-4">
          {systemInfo && (
            <div className="flex items-center justify-between p-4 bg-muted rounded-lg">
              <div>
                <p className="text-sm text-muted-foreground">Current Version</p>
                <p className="text-2xl font-bold">{systemInfo.version}</p>
                {systemInfo.buildDate && (
                  <p className="text-xs text-muted-foreground mt-1">
                    Built: {formatDate(systemInfo.buildDate)}
                  </p>
                )}
              </div>
              {!updateCheck?.hasUpdate && updateCheck !== null && (
                <Badge variant="default" className="bg-green-500">
                  <CheckCircle2 className="h-3 w-3 mr-1" />
                  Up to date
                </Badge>
              )}
            </div>
          )}

          {/* Update Available */}
          {updateCheck?.hasUpdate && updateCheck.release && (
            <Alert>
              <Download className="h-4 w-4" />
              <AlertDescription>
                <div className="space-y-3">
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="font-semibold">
                        New version available: {updateCheck.latestVersion}
                      </p>
                      <p className="text-sm text-muted-foreground">
                        Released: {formatDate(updateCheck.release.published_at)}
                      </p>
                    </div>
                    <Button
                      onClick={handleTriggerUpdate}
                      disabled={updating}
                      size="sm"
                    >
                      {updating ? (
                        <>
                          <RefreshCw className="h-4 w-4 mr-2 animate-spin" />
                          Updating...
                        </>
                      ) : (
                        <>
                          <Download className="h-4 w-4 mr-2" />
                          Update Now
                        </>
                      )}
                    </Button>
                  </div>

                  {/* Changelog Toggle */}
                  <div>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => setShowChangelog(!showChangelog)}
                      className="px-0"
                    >
                      {showChangelog ? (
                        <>
                          <ChevronUp className="h-4 w-4 mr-1" />
                          Hide Changelog
                        </>
                      ) : (
                        <>
                          <ChevronDown className="h-4 w-4 mr-1" />
                          Show Changelog
                        </>
                      )}
                    </Button>

                    {showChangelog && updateCheck.release.body && (
                      <div 
                        className="mt-3 p-4 bg-background rounded border text-sm space-y-2"
                        dangerouslySetInnerHTML={{ 
                          __html: parseMarkdown(updateCheck.release.body) 
                        }}
                      />
                    )}
                  </div>
                </div>
              </AlertDescription>
            </Alert>
          )}

          {updating && (
            <Alert>
              <AlertCircle className="h-4 w-4" />
              <AlertDescription>
                Update in progress... The page will reload automatically when complete.
                Please do not close this window.
              </AlertDescription>
            </Alert>
          )}
        </CardContent>
      </Card>

      {/* Update History */}
      <Card>
        <CardHeader>
          <Button
            variant="ghost"
            onClick={() => setShowHistory(!showHistory)}
            className="w-full justify-between p-0 hover:bg-transparent"
          >
            <CardTitle className="flex items-center gap-2">
              <Clock className="h-5 w-5 text-primary" />
              Update History
            </CardTitle>
            {showHistory ? <ChevronUp /> : <ChevronDown />}
          </Button>
        </CardHeader>
        {showHistory && (
          <CardContent>
            {updateHistory.length === 0 ? (
              <p className="text-sm text-muted-foreground text-center py-4">
                No update history available
              </p>
            ) : (
              <div className="space-y-3">
                {updateHistory.map((update) => (
                  <div
                    key={update.id}
                    className="flex items-center justify-between p-3 border rounded-lg"
                  >
                    <div className="flex-1">
                      <div className="flex items-center gap-2 mb-1">
                        <span className="font-medium">
                          {update.fromVersion} → {update.toVersion}
                        </span>
                        {getStatusBadge(update.status)}
                      </div>
                      <p className="text-sm text-muted-foreground">
                        {formatDate(update.startedAt)}
                      </p>
                      {update.errorMessage && (
                        <p className="text-sm text-destructive mt-1">
                          {update.errorMessage}
                        </p>
                      )}
                    </div>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        )}
      </Card>
    </div>
  );
};
