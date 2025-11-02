import { create } from 'zustand';

export interface DiskInfo {
  name: string;
  totalSpace: number;
  availableSpace: number;
  usedSpace: number;
}

export interface SystemStats {
  cpuUsage: number;
  memory: {
    total: number;
    used: number;
  };
  disk: DiskInfo[];
  loadAverage: {
    oneMinute: number;
    fiveMinutes: number;
    fifteenMinutes: number;
  };
  timestamp: number;
  uptime: number;
  cpuTemperature: number;
}

interface DashboardState {
  stats: SystemStats[];
  isConnected: boolean;
  isLoading: boolean;
  error: string | null;
  wsConnection: WebSocket | null;

  connectWebSocket: () => void;
  disconnectWebSocket: () => void;
  addStats: (stats: SystemStats) => void;
  setError: (error: string | null) => void;
  clearStats: () => void;
  
  getLatestStats: () => SystemStats | null;
  getAverageCpuUsage: () => number;
  getMemoryUsagePercentage: () => number;
}

export const useDashboardStore = create<DashboardState>((set, get) => ({
  stats: [],
  isConnected: false,
  isLoading: true,
  error: null,
  wsConnection: null,

  connectWebSocket: () => {
    const { wsConnection, disconnectWebSocket } = get();
    
    if (wsConnection) {
      disconnectWebSocket();
    }

    try {
      const ws = new WebSocket('/api/ws/stats');
      
      ws.onopen = () => {
        set({ 
          isConnected: true, 
          isLoading: false, 
          error: null,
          wsConnection: ws 
        });
      };

      ws.onmessage = (event) => {
        try {
          const data: SystemStats = JSON.parse(event.data);
          get().addStats(data);
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error);
        }
      };

      ws.onerror = (error) => {
        console.error('WebSocket error:', error);
        set({ 
          error: 'WebSocket connection error',
          isConnected: false 
        });
      };

      ws.onclose = () => {
        set({ 
          isConnected: false,
          wsConnection: null 
        });
      };

    } catch (error) {
      set({ 
        error: 'Failed to create WebSocket connection',
        isLoading: false 
      });
    }
  },

  disconnectWebSocket: () => {
    const { wsConnection } = get();
    if (wsConnection) {
      wsConnection.close();
      set({ 
        wsConnection: null, 
        isConnected: false 
      });
    }
  },

  addStats: (newStats: SystemStats) => {
    set(state => ({
      stats: [...state.stats.slice(-30), newStats] // Keep last 30 entries
    }));
  },

  setError: (error: string | null) => {
    set({ error });
  },

  clearStats: () => {
    set({ stats: [] });
  },

  getLatestStats: () => {
    const { stats } = get();
    return stats.length > 0 ? stats[stats.length - 1] : null;
  },

  getAverageCpuUsage: () => {
    const { stats } = get();
    if (stats.length === 0) return 0;
    
    const total = stats.reduce((acc, curr) => acc + curr.cpuUsage, 0);
    return total / stats.length;
  },

  getMemoryUsagePercentage: () => {
    const latest = get().getLatestStats();
    if (!latest) return 0;
    
    return (latest.memory.used / latest.memory.total) * 100;
  },
}));
