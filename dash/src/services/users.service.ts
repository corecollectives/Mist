import type { User, CreateUserData } from '@/types';

const API_BASE = '/api';

export const usersService = {
  /**
   * Get all users
   */
  async getAll(): Promise<User[]> {
    const response = await fetch(`${API_BASE}/users/getAll`, {
      credentials: 'include',
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to fetch users');
    }

    // Transform users to include isAdmin flag
    const users: User[] = (data.data || []).map((u: User) => ({
      ...u,
      isAdmin: u.role === 'admin' || u.role === 'owner',
    }));

    return users;
  },

  /**
   * Create new user
   */
  async create(userData: CreateUserData): Promise<User> {
    const response = await fetch(`${API_BASE}/users/create`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(userData),
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to create user');
    }

    return data.data;
  },

  /**
   * Update user
   */
  async update(userId: number, updates: Partial<User>): Promise<User> {
    const response = await fetch(`${API_BASE}/users/update`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ id: userId, ...updates }),
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to update user');
    }

    return data.data;
  },

  /**
   * Delete user
   */
  async delete(userId: number): Promise<void> {
    const response = await fetch(`${API_BASE}/users/delete?id=${userId}`, {
      method: 'DELETE',
      credentials: 'include',
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to delete user');
    }
  },
};
