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
  createdAt: string;   // or Date if you parse it
  updatedAt: string;   // or Date if you parse it
};
