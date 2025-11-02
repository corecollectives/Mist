
export interface User {
  id: number;
  username: string;
  email: string;
  role: 'owner' | 'admin' | 'user';
  isAdmin?: boolean;
  createdAt?: string;
  updatedAt?: string;
}

export interface LoginCredentials {
  email: string;
  password: string;
}

export interface SignupCredentials {
  username: string;
  email: string;
  password: string;
}

export interface UserCreateInput {
  username: string;
  email: string;
  password: string;
  role: 'admin' | 'user';
}

export type CreateUserData = UserCreateInput;
