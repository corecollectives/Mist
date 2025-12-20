import { apiClient } from '@/api';

export interface SystemSettings {
  wildcardDomain: string | null;
  mistAppName: string;
}

export const settingsService = {
  async getSystemSettings(): Promise<SystemSettings> {
    const response = await apiClient.get<SystemSettings>('/settings/system');
    return response.data;
  },

  async updateSystemSettings(
    wildcardDomain: string | null,
    mistAppName: string
  ): Promise<SystemSettings> {
    const response = await apiClient.put<SystemSettings>('/settings/system', {
      wildcardDomain,
      mistAppName,
    });
    return response.data;
  },
};
