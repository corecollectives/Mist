export type App = {
  id: number;
  projectId: number;
  createdBy: number;
  name: string;
  description: string | null;
  gitProviderId: number | null;
  gitRepository: string | null;
  gitBranch: string;
  gitCloneUrl: string | null;
  deploymentStrategy: string;
  port: number | null;
  rootDirectory: string;
  buildCommand: string | null;
  startCommand: string | null;
  dockerfilePath: string | null;
  healthcheckPath: string | null;
  healthcheckInterval: number;
  status: string;
  createdAt: string;
  updatedAt: string;
};

export type CreateAppRequest = {
  projectId: number;
  name: string;
  description?: string;
  gitRepository?: string;
  gitBranch?: string;
  port?: number;
  rootDirectory?: string;
  buildCommand?: string;
  startCommand?: string;
};

export type UpdateAppRequest = Partial<Omit<App, 'id' | 'createdAt' | 'updatedAt'>>;

export type EnvVariable = {
  id: number;
  appId: number;
  key: string;
  value: string;
  createdAt: string;
  updatedAt: string;
};

export type CreateEnvVariableRequest = {
  appId: number;
  key: string;
  value: string;
};

export type UpdateEnvVariableRequest = {
  id: number;
  key: string;
  value: string;
};

export type Domain = {
  id: number;
  appId: number;
  domain: string;
  sslStatus: string;
  createdAt: string;
  updatedAt: string;
};

export type CreateDomainRequest = {
  appId: number;
  domain: string;
};

export type UpdateDomainRequest = {
  id: number;
  domain: string;
};

export type ContainerStatus = {
  name: string;
  status: string;
  state: string;
  uptime: string;
  healthy: boolean;
};
