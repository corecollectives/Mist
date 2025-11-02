import { create } from 'zustand';
import { gitApi } from '@/api/endpoints/git';
import type { GitHubApp, Repository } from '@/types';

interface GitState {
  app: GitHubApp | null;
  repositories: Repository[];
  isInstalled: boolean;
  isLoading: boolean;
  error: string | null;
  lastFetch: number | null;

  fetchApp: () => Promise<void>;
  fetchRepositories: () => Promise<void>;
  createApp: () => Promise<void>;
  installApp: (appId: number, userId: number) => Promise<void>;
  clearError: () => void;
  
  getRepositoryById: (id: string) => Repository | undefined;
  getRepositoriesByProvider: (provider: string) => Repository[];
}

export const useGitStore = create<GitState>((set, get) => ({
  app: null,
  repositories: [],
  isInstalled: false,
  isLoading: false,
  error: null,
  lastFetch: null,

  fetchApp: async () => {
    const now = Date.now();
    const { lastFetch } = get();
    
    if (lastFetch && now - lastFetch < 10 * 60 * 1000) {
      return;
    }

    set({ isLoading: true, error: null });
    
    try {
      const response = await gitApi.getApp();
      set({ 
        app: response.data.app,
        isInstalled: response.data.isInstalled,
        isLoading: false,
        lastFetch: now
      });
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to fetch GitHub app';
      set({ 
        error: errorMessage, 
        isLoading: false,
        app: null,
        isInstalled: false
      });
    }
  },

  fetchRepositories: async () => {
    set({ isLoading: true, error: null });
    
    try {
      const response = await gitApi.getRepositories();
      set({ 
        repositories: response.data,
        isLoading: false
      });
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to fetch repositories';
      set({ 
        error: errorMessage, 
        isLoading: false 
      });
    }
  },

  createApp: async () => {
    set({ isLoading: true, error: null });
    
    try {
      const response = await gitApi.createApp();
      set({ 
        app: response.data,
        isLoading: false
      });
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to create GitHub app';
      set({ 
        error: errorMessage, 
        isLoading: false 
      });
      throw error;
    }
  },

  installApp: async (appId: number, userId: number) => {
    set({ isLoading: true, error: null });
    
    try {
      await gitApi.installApp(appId, userId);
      set({ 
        isInstalled: true,
        isLoading: false
      });
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to install GitHub app';
      set({ 
        error: errorMessage, 
        isLoading: false 
      });
      throw error;
    }
  },

  clearError: () => {
    set({ error: null });
  },

  getRepositoryById: (id: string) => {
    const { repositories } = get();
    return repositories.find(repo => repo.id === id);
  },

  getRepositoriesByProvider: (provider: string) => {
    const { repositories } = get();
    return repositories.filter(repo => repo.provider === provider);
  },
}));
