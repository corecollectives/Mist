
import { apiClient, type ApiResponse } from '../client';
import type { User, UserCreateInput } from '../../types/user';

export const usersApi = {
  async getAll(): Promise<ApiResponse<User[]>> {
    return apiClient.get<User[]>('/users/getAll');
  },

  async getUsers(): Promise<ApiResponse<User[]>> {
    return this.getAll();
  },

  async create(userData: UserCreateInput): Promise<ApiResponse<User>> {
    return apiClient.post<User>('/users/create', userData);
  },

  async createUser(userData: UserCreateInput): Promise<ApiResponse<User>> {
    return this.create(userData);
  },
} as const;
