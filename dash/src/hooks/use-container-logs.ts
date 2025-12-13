import { useEffect, useRef, useState, useCallback } from 'react';

interface ContainerLogEvent {
  type: 'log' | 'status' | 'error' | 'end';
  timestamp: string;
  data: {
    line?: string;
    message?: string;
    container?: string;
    state?: string;
    status?: string;
  };
}

interface UseContainerLogsOptions {
  appId: number;
  enabled: boolean;
  onError?: (error: string) => void;
}

export const useContainerLogs = ({
  appId,
  enabled,
  onError,
}: UseContainerLogsOptions) => {
  const [logs, setLogs] = useState<string[]>([]);
  const [containerState, setContainerState] = useState<string>('');
  const [error, setError] = useState<string | null>(null);
  const [isConnected, setIsConnected] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<number>(0);
  const reconnectAttemptsRef = useRef(0);

  const connectWebSocket = useCallback(() => {
    if (!enabled) return;

    // Prevent duplicate connections
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      console.log('[ContainerLogs] WebSocket already connected, skipping');
      return;
    }

    try {
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const host = window.location.host;
      const ws = new WebSocket(`${protocol}//${host}/api/ws/container/logs?appId=${appId}`);
      wsRef.current = ws;

      ws.onopen = () => {
        console.log('[ContainerLogs] WebSocket connected');
        setIsConnected(true);
        setError(null);
        setIsLoading(false);
        reconnectAttemptsRef.current = 0;
      };

      ws.onmessage = (event) => {
        if (event.data instanceof Blob) {
          return;
        }

        try {
          const logEvent: ContainerLogEvent = JSON.parse(event.data);

          switch (logEvent.type) {
            case 'log':
              if (logEvent.data.line && logEvent.data.line.trim()) {
                setLogs((prev) => [...prev, logEvent.data.line!]);
              }
              break;

            case 'status':
              if (logEvent.data.state) {
                setContainerState(logEvent.data.state);
              }
              break;

            case 'error':
              const errorMsg = logEvent.data.message || 'Unknown error';
              console.error('[ContainerLogs] Error:', errorMsg);
              setError(errorMsg);
              onError?.(errorMsg);
              break;

            case 'end':
              console.log('[ContainerLogs] Log stream ended');
              setError('Log stream ended');
              break;
          }
        } catch (err) {
          console.error('[ContainerLogs] Error parsing message:', err);
        }
      };

      ws.onerror = (event) => {
        console.error('[ContainerLogs] WebSocket error:', event);
        setError('Connection error occurred');
        setIsConnected(false);
      };

      ws.onclose = () => {
        console.log('[ContainerLogs] WebSocket closed');
        setIsConnected(false);
        wsRef.current = null;

        if (enabled && reconnectAttemptsRef.current < 5) {
          const delay = Math.min(1000 * Math.pow(2, reconnectAttemptsRef.current), 10000);
          reconnectAttemptsRef.current++;

          console.log(`[ContainerLogs] Reconnecting in ${delay}ms (attempt ${reconnectAttemptsRef.current})`);
          reconnectTimeoutRef.current = window.setTimeout(() => {
            connectWebSocket();
          }, delay);
        } else if (reconnectAttemptsRef.current >= 5) {
          setError('Failed to connect after multiple attempts');
          setIsLoading(false);
        }
      };
    } catch (err) {
      console.error('[ContainerLogs] Error creating WebSocket:', err);
      setError('Failed to establish connection');
      setIsLoading(false);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [appId, enabled]);

  useEffect(() => {
    if (!enabled) {
      return;
    }

    if (!wsRef.current || wsRef.current.readyState === WebSocket.CLOSED) {
      connectWebSocket();
    } else {
      // Close connection when disabled
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }
      if (wsRef.current) {
        wsRef.current.close();
        wsRef.current = null;
      }
      setIsConnected(false);
      setIsLoading(false);
    }

    return () => {
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }
      if (wsRef.current) {
        wsRef.current.close();
        wsRef.current = null;
      }
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [enabled, appId]);

  const clearLogs = () => {
    setLogs([]);
  };

  const disconnect = () => {
    if (wsRef.current) {
      wsRef.current.close();
    }
    setIsConnected(false);
  };

  return {
    logs,
    containerState,
    error,
    isConnected,
    isLoading,
    clearLogs,
    disconnect,
  };
};
