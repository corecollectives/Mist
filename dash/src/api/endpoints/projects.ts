import { apiClient, type ApiResponse } from '../client';
import type { Project, ProjectCreateInput, ProjectUpdateInput } from '../../types/project';

export const projectsApi = {
  async getAll(): Promise<ApiResponse<Project[]>> {
    return apiClient.get<Project[]>('/projects/getAll');
  },

  async getById(projectId: string | number): Promise<ApiResponse<Project>> {
    return apiClient.get<Project>(`/projects/getFromId?id=${projectId}`);
  },

  async create(projectData: ProjectCreateInput): Promise<ApiResponse<Project>> {
    return apiClient.post<Project>('/projects/create', projectData);
  },

  async update(
    projectId: string | number,
    projectData: ProjectUpdateInput
  ): Promise<ApiResponse<Project>> {
    return apiClient.put<Project>(`/projects/update?id=${projectId}`, projectData);
  },

  async delete(projectId: string | number): Promise<ApiResponse<void>> {
    return apiClient.delete<void>(`/projects/delete?id=${projectId}`);
  },
} as const;
