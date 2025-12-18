import type { SystemInfo, UpdateCheck, UpdateHistory, SystemHealth, TriggerUpdateRequest, TriggerUpdateResponse } from '@/types/system';

const API_BASE = '/api';

export const systemService = {
  async getVersion(): Promise<SystemInfo> {
    const response = await fetch(`${API_BASE}/system/version`, {
      credentials: 'include',
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to get version');
    }

    return data.data;
  },

  async checkForUpdates(): Promise<UpdateCheck> {
    const response = await fetch(`${API_BASE}/system/updates/check`, {
      credentials: 'include',
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to check for updates');
    }

    return data.data;
  },

  async triggerUpdate(request: TriggerUpdateRequest): Promise<TriggerUpdateResponse> {
    const response = await fetch(`${API_BASE}/system/updates/trigger`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(request),
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to trigger update');
    }

    return data.data;
  },

  async getUpdateHistory(): Promise<UpdateHistory[]> {
    const response = await fetch(`${API_BASE}/system/updates/history`, {
      credentials: 'include',
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to get update history');
    }

    return data.data || [];
  },

  async getUpdateStatus(): Promise<UpdateHistory | null> {
    const response = await fetch(`${API_BASE}/system/updates/status`, {
      credentials: 'include',
    });

    const data = await response.json();
    if (!data.success) {
      return null;
    }

    return data.data;
  },

  async getSystemHealth(): Promise<SystemHealth> {
    const response = await fetch(`${API_BASE}/system/health`, {
      credentials: 'include',
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to get system health');
    }

    return data.data;
  },
};
