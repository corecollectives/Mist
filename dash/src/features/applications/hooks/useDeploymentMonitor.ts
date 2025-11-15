import { useEffect, useRef, useState, useCallback } from 'react';
import type { DeploymentEvent, StatusUpdate, LogUpdate, Deployment } from '@/types/deployment';

interface UseDeploymentMonitorOptions {
  deploymentId: number;
  enabled: boolean;
  onComplete?: () => void;
  onError?: (error: string) => void;
}

export const useDeploymentMonitor = ({
  deploymentId,
  enabled,
  onComplete,
  onError,
}: UseDeploymentMonitorOptions) => {
  const [logs, setLogs] = useState<string[]>([]);
  const [status, setStatus] = useState<StatusUpdate | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isConnected, setIsConnected] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [isLive, setIsLive] = useState(false);

  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<number>(0);
  const reconnectAttemptsRef = useRef(0);
  const hasFetchedRef = useRef(false);

  const fetchCompletedDeployment = useCallback(async () => {
    if (hasFetchedRef.current) return;

    try {
      setIsLoading(true);
      const response = await fetch(`/api/deployments/logs?id=${deploymentId}`, {
        credentials: 'include',
      });

      if (!response.ok) {
        if (response.status === 400) {
          setIsLive(true);
          setIsLoading(false);
          hasFetchedRef.current = true;
          return;
        }
        throw new Error('Failed to fetch deployment logs');
      }

      const result = await response.json();
      const deployment: Deployment = result.data.deployment;
      const logsContent: string = result.data.logs;

      if (logsContent) {
        setLogs(logsContent.split('\n').filter(line => line.length > 0));
      }

      setStatus({
        deployment_id: deployment.id,
        status: deployment.status,
        stage: deployment.stage,
        progress: deployment.progress,
        message: deployment.status === 'success'
          ? 'Deployment completed successfully'
          : 'Deployment failed',
        error_message: deployment.error_message,
        duration: deployment.duration,
      });

      if (deployment.status === 'failed' && deployment.error_message) {
        setError(deployment.error_message);
      } else if (deployment.status === 'success') {
      }

      setIsLoading(false);
      hasFetchedRef.current = true;
    } catch (err) {
      console.error('[DeploymentMonitor] Error fetching completed deployment:', err);
      setError('Failed to load deployment logs');
      setIsLoading(false);
      hasFetchedRef.current = true;
    }
  }, [deploymentId]);

  const connectWebSocket = useCallback(() => {
    if (!enabled || !isLive) return;

    try {
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const host = window.location.host;
      const ws = new WebSocket(`${protocol}//${host}/api/deployments/logs/stream?id=${deploymentId}`)
      wsRef.current = ws;

      ws.onopen = () => {
        console.log('[DeploymentMonitor] WebSocket connected');
        setIsConnected(true);
        setError(null);
        setIsLoading(false);
        reconnectAttemptsRef.current = 0;
      };

      ws.onmessage = (event) => {
        try {
          const deploymentEvent: DeploymentEvent = JSON.parse(event.data);

          switch (deploymentEvent.type) {
            case 'log':
              const logData = deploymentEvent.data as LogUpdate;
              setLogs((prev) => [...prev, logData.line]);
              break;

            case 'status':
              const statusData = deploymentEvent.data as StatusUpdate;
              setStatus(statusData);

              // Handle completion
              if (statusData.status === 'success') {
                onComplete?.();
              }

              // Handle errors
              if (statusData.status === 'failed' && statusData.error_message) {
                setError(statusData.error_message);
                onError?.(statusData.error_message);
              }
              break;

            case 'error':
              const errorMsg = (deploymentEvent.data as any).message || 'Unknown error';
              setError(errorMsg);
              onError?.(errorMsg);
              break;
          }
        } catch (err) {
          console.error('[DeploymentMonitor] Error parsing message:', err);
        }
      };

      ws.onerror = (event) => {
        console.error('[DeploymentMonitor] WebSocket error:', event);
        setError('Connection error occurred');
        setIsConnected(false);
      };

      ws.onclose = (event) => {
        console.log('[DeploymentMonitor] WebSocket closed:', event.code, event.reason);
        setIsConnected(false);
        wsRef.current = null;

        if (status?.status === 'success' || status?.status === 'failed') {
          return;
        }

        if (reconnectAttemptsRef.current < 10) {
          const delay = Math.min(1000 * Math.pow(2, reconnectAttemptsRef.current), 30000);
          reconnectAttemptsRef.current++;

          console.log(`[DeploymentMonitor] Reconnecting in ${delay}ms (attempt ${reconnectAttemptsRef.current})...`);
          reconnectTimeoutRef.current = window.setTimeout(() => {
            connectWebSocket();
          }, delay);
        } else {
          setError('Failed to connect after multiple attempts');
          setIsLoading(false);
        }
      };
    } catch (err) {
      console.error('[DeploymentMonitor] Error creating WebSocket:', err);
      setError('Failed to establish connection');
      setIsLoading(false);
    }
  }, [deploymentId, enabled, isLive, onComplete, onError, status?.status]);

  useEffect(() => {
    if (enabled) {
      fetchCompletedDeployment();
    }

    return () => {
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, [enabled, fetchCompletedDeployment]);

  useEffect(() => {
    if (isLive && enabled) {
      connectWebSocket();
    }

    return () => {
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, [isLive, enabled, connectWebSocket]);

  const reset = () => {
    setLogs([]);
    setStatus(null);
    setError(null);
    setIsConnected(false);
    setIsLoading(true);
    setIsLive(false);
    hasFetchedRef.current = false;
  };

  return {
    logs,
    status,
    error,
    isConnected,
    isLoading,
    isLive,
    reset,
  };
};
