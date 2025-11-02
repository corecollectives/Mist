
import { apiClient, type ApiResponse } from '../client';
import type { User, LoginCredentials, SignupCredentials } from '../../types';

export const authApi = {
  async getMe(): Promise<ApiResponse<{ setupRequired?: boolean } & User>> {
    return apiClient.get('/auth/me');
  },

  async me(): Promise<ApiResponse<{ setupRequired?: boolean } & User>> {
    return this.getMe();
  },

  async login(credentials: LoginCredentials): Promise<ApiResponse<User>> {
    return apiClient.post<User>('/auth/login', credentials);
  },

  async signup(credentials: SignupCredentials): Promise<ApiResponse<User>> {
    return apiClient.post<User>('/auth/signup', credentials);
  },

  async logout(): Promise<ApiResponse<void>> {
    return apiClient.post<void>('/auth/logout');
  },
} as const;
