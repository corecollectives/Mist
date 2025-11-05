export type DeploymentStrategy = 'auto' | 'manual';
export type AppStatus = 'stopped' | 'running' | 'error' | 'building';
export interface App {
  id: number;
  project_id: number;
  created_by: number | null;
  name: string;
  description?: string | null;
  git_provider_id?: number | null;
  git_repository?: string | null;
  git_branch?: string | null; // default: 'main'
  deployment_strategy: DeploymentStrategy; // default: 'manual'
  port?: number | null;
  root_directory?: string | null;
  build_command?: string | null;
  start_command?: string | null;
  dockerfile_path?: string | null;
  healthcheck_path?: string | null;
  healthcheck_interval: number; // default: 30
  status: AppStatus; // default: 'stopped'
  created_at: string; // ISO datetime
  updated_at: string; // ISO datetime
}
