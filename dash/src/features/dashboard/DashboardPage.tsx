import { useEffect } from 'react';
import { toast } from 'sonner';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { FullScreenLoading } from '@/shared/components';
import { useDashboardStore } from './store';
import { SystemOverview, ChartCard, MetricCard } from './components';
import { formatPercentage, formatMemory } from './utils';

export default function DashboardPage() {
  const {
    stats,
    isConnected,
    isLoading,
    error,
    connectWebSocket,
    disconnectWebSocket,
    getLatestStats,
    getAverageCpuUsage,
    getMemoryUsagePercentage
  } = useDashboardStore();

  useEffect(() => {
    connectWebSocket();
    
    return () => {
      disconnectWebSocket();
    };
  }, [connectWebSocket, disconnectWebSocket]);

  useEffect(() => {
    if (error) {
      toast.error(error);
    }
  }, [error]);

  const latestStats = getLatestStats();
  const averageCpuUsage = getAverageCpuUsage();
  const memoryUsagePercentage = getMemoryUsagePercentage();

  if (isLoading) {
    return <FullScreenLoading />;
  }

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <div className="flex items-center justify-between py-6 border-b border-border shrink-0">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-foreground">
            Dashboard
          </h1>
          <p className="text-muted-foreground mt-1">
            System monitoring and performance overview
          </p>
        </div>
        <div className="flex items-center gap-2">
          <div className={`w-3 h-3 rounded-full ${isConnected ? 'bg-green-500' : 'bg-red-500'}`} />
          <span className="text-sm text-muted-foreground">
            {isConnected ? 'Connected' : 'Disconnected'}
          </span>
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

      {/* Content */}
      <div className="py-6 space-y-6">
        {/* System Overview */}
        <SystemOverview stats={latestStats} />

        {/* Charts */}
        {stats.length > 0 && (
          <div className="grid gap-6 md:grid-cols-2">
            <ChartCard
              title="CPU Usage"
              data={stats}
              dataKey="cpuUsage"
              color="#8B5CF6"
              formatter={formatPercentage}
            />
            
            <ChartCard
              title="Memory Usage"
              data={stats}
              dataKey="memory.used"
              color="#A371F7"
              formatter={formatMemory}
            />
          </div>
        )}

        {/* Summary Stats */}
        {stats.length > 0 && (
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
            <MetricCard
              title="Current CPU"
              value={latestStats?.cpuUsage || 0}
              showUsageColor
            />
            
            <MetricCard
              title="Average CPU"
              value={averageCpuUsage}
              showUsageColor
            />
            
            <MetricCard
              title="Memory Usage"
              value={memoryUsagePercentage}
              showUsageColor
            />
            
            <MetricCard
              title="Load Average"
              value={latestStats?.loadAverage.oneMinute || 0}
              formatter={(value: number) => value.toFixed(2)}
            />
          </div>
        )}

        {/* No Data State */}
        {stats.length === 0 && !isLoading && (
          <div className="flex flex-col items-center justify-center py-12">
            <p className="text-muted-foreground text-lg mb-4">
              {isConnected ? 'Waiting for system data...' : 'Unable to connect to system monitoring'}
            </p>
          </div>
        )}
      </div>
    </div>
  );
}
