/**
 * Frontend Deployment Monitor Component
 * 
 * This component provides real-time deployment monitoring with:
 * - Live log streaming
 * - Status updates
 * - Progress tracking
 * - Error handling and display
 */

import React, { useEffect, useRef, useState } from 'react';

// Types matching backend models
interface DeploymentEvent {
  type: 'log' | 'status' | 'progress' | 'error';
  timestamp: string;
  data: LogUpdate | StatusUpdate;
}

interface LogUpdate {
  line: string;
  timestamp: string;
}

interface StatusUpdate {
  deployment_id: number;
  status: string;
  stage: string;
  progress: number;
  message: string;
  error_message?: string;
}

interface Deployment {
  id: number;
  app_id: number;
  commit_hash: string;
  commit_message: string;
  status: string;
  stage: string;
  progress: number;
  error_message?: string;
  created_at: string;
  started_at?: string;
  finished_at?: string;
  duration?: number;
}

interface DeploymentMonitorProps {
  deploymentId: number;
  onComplete?: (deployment: Deployment) => void;
  onError?: (error: string) => void;
}

const DeploymentMonitor: React.FC<DeploymentMonitorProps> = ({
  deploymentId,
  onComplete,
  onError,
}) => {
  const [logs, setLogs] = useState<string[]>([]);
  const [status, setStatus] = useState<StatusUpdate | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isConnected, setIsConnected] = useState(false);
  const wsRef = useRef<WebSocket | null>(null);
  const logsEndRef = useRef<HTMLDivElement>(null);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout>();
  const reconnectAttemptsRef = useRef(0);

  // Auto-scroll to bottom of logs
  const scrollToBottom = () => {
    logsEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [logs]);

  // WebSocket connection with reconnection logic
  const connect = () => {
    try {
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const wsUrl = `${protocol}//${window.location.host}/api/deployments/logs?id=${deploymentId}`;
      
      const ws = new WebSocket(wsUrl);
      wsRef.current = ws;

      ws.onopen = () => {
        console.log('WebSocket connected');
        setIsConnected(true);
        setError(null);
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
                onComplete?.({
                  id: statusData.deployment_id,
                  status: statusData.status,
                  stage: statusData.stage,
                  progress: statusData.progress,
                } as Deployment);
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
          console.error('Error parsing WebSocket message:', err);
        }
      };

      ws.onerror = (event) => {
        console.error('WebSocket error:', event);
        setError('Connection error occurred');
        setIsConnected(false);
      };

      ws.onclose = (event) => {
        console.log('WebSocket closed:', event.code, event.reason);
        setIsConnected(false);
        wsRef.current = null;

        // Don't reconnect if deployment is finished
        if (status?.status === 'success' || status?.status === 'failed') {
          return;
        }

        // Reconnect with exponential backoff
        if (reconnectAttemptsRef.current < 10) {
          const delay = Math.min(1000 * Math.pow(2, reconnectAttemptsRef.current), 30000);
          reconnectAttemptsRef.current++;
          
          console.log(`Reconnecting in ${delay}ms (attempt ${reconnectAttemptsRef.current})...`);
          reconnectTimeoutRef.current = setTimeout(() => {
            connect();
          }, delay);
        } else {
          setError('Failed to connect after multiple attempts');
        }
      };
    } catch (err) {
      console.error('Error creating WebSocket:', err);
      setError('Failed to establish connection');
    }
  };

  useEffect(() => {
    connect();

    // Cleanup
    return () => {
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, [deploymentId]);

  // Status badge color
  const getStatusColor = (status?: string) => {
    switch (status) {
      case 'success':
        return 'bg-green-500';
      case 'failed':
        return 'bg-red-500';
      case 'building':
      case 'deploying':
      case 'cloning':
        return 'bg-blue-500 animate-pulse';
      case 'pending':
        return 'bg-yellow-500';
      default:
        return 'bg-gray-500';
    }
  };

  return (
    <div className="deployment-monitor flex flex-col h-full">
      {/* Header with status */}
      <div className="bg-gray-800 text-white p-4 flex items-center justify-between">
        <div className="flex items-center space-x-4">
          <h2 className="text-xl font-bold">Deployment #{deploymentId}</h2>
          {status && (
            <div className="flex items-center space-x-2">
              <span className={`px-3 py-1 rounded-full text-sm font-medium ${getStatusColor(status.status)}`}>
                {status.status.toUpperCase()}
              </span>
              <span className="text-gray-400">{status.message}</span>
            </div>
          )}
        </div>
        
        {/* Connection indicator */}
        <div className="flex items-center space-x-2">
          <div className={`w-2 h-2 rounded-full ${isConnected ? 'bg-green-500' : 'bg-red-500'}`} />
          <span className="text-sm text-gray-400">
            {isConnected ? 'Connected' : 'Disconnected'}
          </span>
        </div>
      </div>

      {/* Progress bar */}
      {status && status.status !== 'success' && status.status !== 'failed' && (
        <div className="bg-gray-700 px-4 py-2">
          <div className="flex items-center justify-between mb-1">
            <span className="text-sm text-gray-300">{status.stage}</span>
            <span className="text-sm text-gray-300">{status.progress}%</span>
          </div>
          <div className="w-full bg-gray-600 rounded-full h-2">
            <div
              className="bg-blue-500 h-2 rounded-full transition-all duration-300"
              style={{ width: `${status.progress}%` }}
            />
          </div>
        </div>
      )}

      {/* Error banner */}
      {error && (
        <div className="bg-red-500 text-white px-4 py-3 flex items-start">
          <svg className="w-6 h-6 mr-2 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
            <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
          </svg>
          <div>
            <h3 className="font-bold">Deployment Failed</h3>
            <p className="text-sm mt-1">{error}</p>
          </div>
        </div>
      )}

      {/* Success banner */}
      {status?.status === 'success' && (
        <div className="bg-green-500 text-white px-4 py-3 flex items-center">
          <svg className="w-6 h-6 mr-2" fill="currentColor" viewBox="0 0 20 20">
            <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
          </svg>
          <span className="font-bold">Deployment Successful!</span>
        </div>
      )}

      {/* Logs viewer */}
      <div className="flex-1 bg-black text-green-400 p-4 overflow-auto font-mono text-sm">
        {logs.length === 0 ? (
          <div className="text-gray-500 text-center py-8">
            Waiting for logs...
          </div>
        ) : (
          logs.map((log, index) => (
            <div key={index} className="whitespace-pre-wrap break-words">
              {log}
            </div>
          ))
        )}
        <div ref={logsEndRef} />
      </div>
    </div>
  );
};

export default DeploymentMonitor;


// ============================================
// Usage Example
// ============================================

const DeploymentPage: React.FC = () => {
  const [deploymentId, setDeploymentId] = useState<number | null>(null);
  const [isDeploying, setIsDeploying] = useState(false);

  const handleDeploy = async (appId: number) => {
    try {
      setIsDeploying(true);
      
      const response = await fetch('/api/deployments', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ appId }),
      });

      if (!response.ok) {
        throw new Error('Failed to create deployment');
      }

      const deployment = await response.json();
      setDeploymentId(deployment.id);
    } catch (error) {
      console.error('Deployment error:', error);
      alert('Failed to start deployment');
      setIsDeploying(false);
    }
  };

  const handleDeploymentComplete = (deployment: Deployment) => {
    console.log('Deployment completed:', deployment);
    setIsDeploying(false);
    // Refresh app status, show success notification, etc.
  };

  const handleDeploymentError = (error: string) => {
    console.error('Deployment error:', error);
    setIsDeploying(false);
    // Show error notification
  };

  return (
    <div className="container mx-auto p-4">
      <button
        onClick={() => handleDeploy(123)}
        disabled={isDeploying}
        className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded disabled:opacity-50"
      >
        {isDeploying ? 'Deploying...' : 'Deploy'}
      </button>

      {deploymentId && (
        <div className="mt-4 h-screen">
          <DeploymentMonitor
            deploymentId={deploymentId}
            onComplete={handleDeploymentComplete}
            onError={handleDeploymentError}
          />
        </div>
      )}
    </div>
  );
};
