import { create } from 'zustand';
import { usersApi } from '@/api/endpoints/users';
import type { User, CreateUserData } from '@/types';

interface UsersState {
  users: User[];
  isLoading: boolean;
  error: string | null;
  lastFetch: number | null;

  fetchUsers: () => Promise<void>;
  createUser: (userData: CreateUserData) => Promise<boolean>;
  clearError: () => void;
  
  getAdminUsers: () => User[];
  getUsersByRole: (role: string) => User[];
}

export const useUsersStore = create<UsersState>((set, get) => ({
  users: [],
  isLoading: false,
  error: null,
  lastFetch: null,

  fetchUsers: async () => {
    const now = Date.now();
    const { lastFetch } = get();
    
    if (lastFetch && now - lastFetch < 5 * 60 * 1000) {
      return;
    }

    set({ isLoading: true, error: null });
    
    try {
      const response = await usersApi.getUsers();
      set({ 
        users: response.data,
        isLoading: false,
        lastFetch: now
      });
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to fetch users';
      set({ 
        error: errorMessage, 
        isLoading: false 
      });
    }
  },

  createUser: async (userData: CreateUserData) => {
    set({ isLoading: true, error: null });
    
    try {
      const response = await usersApi.createUser(userData);
      
      set(state => ({
        users: [...state.users, response.data],
        isLoading: false
      }));
      
      return true;
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to create user';
      set({ 
        error: errorMessage, 
        isLoading: false 
      });
      return false;
    }
  },

  clearError: () => {
    set({ error: null });
  },

  getAdminUsers: () => {
    const { users } = get();
    return users.filter(user => user.role === 'admin' || user.role === 'owner');
  },

  getUsersByRole: (role: string) => {
    const { users } = get();
    return users.filter(user => user.role === role);
  },
}));
