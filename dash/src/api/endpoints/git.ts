
import type { GitHubApp } from '@/types';
import { apiClient, type ApiResponse } from '../client';

export const gitApi = {
  async getApp(): Promise<ApiResponse<{ app: GitHubApp; isInstalled: boolean }>> {
    return apiClient.get('/github/app');
  },

  async getRepositories(): Promise<ApiResponse<any[]>> {
    return apiClient.get('/github/repositories');
  },

  async createApp(): Promise<ApiResponse<GitHubApp>> {
    return apiClient.post('/github/app/create');
  },

  async installApp(appId: number, userId: number): Promise<ApiResponse<void>> {
    return apiClient.post('/github/app/install', { appId, userId });
  },
} as const;
