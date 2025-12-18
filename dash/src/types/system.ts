export interface SystemInfo {
  version: string;
  buildDate: string;
}

export interface GithubRelease {
  tag_name: string;
  name: string;
  body: string;
  draft: boolean;
  prerelease: boolean;
  created_at: string;
  published_at: string;
  tarball_url: string;
  zipball_url: string;
  html_url: string;
}

export interface UpdateCheck {
  hasUpdate: boolean;
  currentVersion: string;
  latestVersion?: string;
  release?: GithubRelease;
  message?: string;
}

export interface UpdateHistory {
  id: number;
  fromVersion: string;
  toVersion: string;
  status: 'pending' | 'downloading' | 'building' | 'installing' | 'success' | 'failed' | 'rolled_back';
  startedAt: string;
  completedAt?: string;
  errorMessage?: string;
  rollbackAvailable: boolean;
  initiatedBy?: number;
}

export interface SystemHealth {
  serviceActive: boolean;
  diskFree: number;
  diskTotal: number;
  uptime: string;
}

export interface TriggerUpdateRequest {
  version?: string;
  branch?: string;
}

export interface TriggerUpdateResponse {
  updateId: number;
  message: string;
}
