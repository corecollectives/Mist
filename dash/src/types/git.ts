
export interface GitHubApp {
  id: number;
  name: string;
  app_id: number;
  client_id: string;
  slug: string;
  created_at: string;
}

export interface Repository {
  id: string;
  name: string;
  fullName: string;
  provider: string;
  url: string;
  isPrivate: boolean;
  description?: string;
  defaultBranch: string;
  createdAt?: string;
  updatedAt?: string;
}
