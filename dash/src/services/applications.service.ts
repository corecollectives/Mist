import type { 
  App, 
  CreateAppRequest, 
  UpdateAppRequest 
} from '@/types';

const API_BASE = '/api';

export const applicationsService = {
  /**
   * Get application by ID
   */
  async getById(appId: number): Promise<App> {
    const response = await fetch(`${API_BASE}/apps/getById`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ appId }),
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to fetch application');
    }

    return data.data;
  },

  /**
   * Get applications by project ID
   */
  async getByProjectId(projectId: number): Promise<App[]> {
    const response = await fetch(`${API_BASE}/apps/getByProjectId`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ projectId }),
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to fetch applications');
    }

    return data.data || [];
  },

  /**
   * Create new application
   */
  async create(request: CreateAppRequest): Promise<App> {
    const response = await fetch(`${API_BASE}/apps/create`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(request),
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to create application');
    }

    return data.data;
  },

  /**
   * Update application
   */
  async update(appId: number, updates: UpdateAppRequest): Promise<App> {
    const response = await fetch(`${API_BASE}/apps/update`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ appId, ...updates }),
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to update application');
    }

    return data.data;
  },

  /**
   * Get latest commit for application
   */
  async getLatestCommit(appId: number): Promise<any> {
    const response = await fetch(`${API_BASE}/apps/getLatestCommit`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ appId }),
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to fetch latest commit');
    }

    return data.data;
  },
};
