import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { toast } from "sonner";
import { applicationsService } from "@/services";
import type { App, RestartPolicy } from "@/types";

interface AppSettingsProps {
  app: App;
  onUpdate: () => void;
}

export const AppSettings = ({ app, onUpdate }: AppSettingsProps) => {
  const [port, setPort] = useState(app.port?.toString() || "");
  const [buildCommand, setBuildCommand] = useState(app.buildCommand || "");
  const [startCommand, setStartCommand] = useState(app.startCommand || "");
  const [rootDirectory, setRootDirectory] = useState(app.rootDirectory || "");
  const [dockerfilePath, setDockerfilePath] = useState(app.dockerfilePath || "");
  const [healthcheckPath, setHealthcheckPath] = useState(app.healthcheckPath || "");
  const [cpuLimit, setCpuLimit] = useState(app.cpuLimit?.toString() || "");
  const [memoryLimit, setMemoryLimit] = useState(app.memoryLimit?.toString() || "");
  const [restartPolicy, setRestartPolicy] = useState(app.restartPolicy || "unless-stopped");
  const [deploymentStrategy, setDeploymentStrategy] = useState(app.deploymentStrategy || "auto");
  const [saving, setSaving] = useState(false);

  const handleSave = async () => {
    try {
      setSaving(true);

      const updates: Partial<{
        rootDirectory: string;
        dockerfilePath: string | null;
        buildCommand: string | null;
        startCommand: string | null;
        healthcheckPath: string | null;
        restartPolicy: RestartPolicy;
        deploymentStrategy: string;
        port: number;
        cpuLimit: number | null;
        memoryLimit: number | null;
      }> = {
        rootDirectory,
        dockerfilePath: dockerfilePath || null,
        buildCommand: buildCommand || null,
        startCommand: startCommand || null,
        healthcheckPath: healthcheckPath || null,
        restartPolicy: restartPolicy as RestartPolicy,
        deploymentStrategy,
      };

      if (port) {
        const portNum = parseInt(port);
        if (isNaN(portNum) || portNum < 1 || portNum > 65535) {
          toast.error("Port must be a number between 1 and 65535");
          return;
        }
        updates.port = portNum;
      }

      if (cpuLimit) {
        const cpu = parseFloat(cpuLimit);
        if (isNaN(cpu) || cpu <= 0) {
          toast.error("CPU limit must be a positive number");
          return;
        }
        updates.cpuLimit = cpu;
      } else {
        updates.cpuLimit = null;
      }

      if (memoryLimit) {
        const memory = parseInt(memoryLimit);
        if (isNaN(memory) || memory <= 0) {
          toast.error("Memory limit must be a positive number");
          return;
        }
        updates.memoryLimit = memory;
      } else {
        updates.memoryLimit = null;
      }

      await applicationsService.update(app.id, updates);
      toast.success("Settings updated successfully");
      onUpdate();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to update settings");
    } finally {
      setSaving(false);
    }
  };

  return (
    <Card>
      <CardHeader className="flex justify-between items-center">
        <div className="flex flex-col gap-2">
          <CardTitle>Application Settings</CardTitle>
          <CardDescription>
            Configure your application settings. Changes will be applied on next deployment.
          </CardDescription>
        </div>
        <div className="flex justify-end">
          <Button onClick={handleSave} disabled={saving}>
            {saving ? "Saving..." : "Save Settings"}
          </Button>
        </div>
      </CardHeader>
      <CardContent className="space-y-6">
        {/* App Type Badge */}
        <div className="flex items-center gap-2 pb-4 border-b">
          <Label>Application Type:</Label>
          <Badge variant="secondary" className="capitalize">
            {app.appType || 'web'}
          </Badge>
          {app.templateName && (
            <>
              <span className="text-muted-foreground">•</span>
              <Badge variant="outline">{app.templateName}</Badge>
            </>
          )}
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="space-y-2">
            <Label htmlFor="port">Port</Label>
            <Input
              id="port"
              type="number"
              placeholder="3000"
              value={port}
              onChange={(e) => setPort(e.target.value)}
              min="1"
              max="65535"
              disabled={app.appType === 'database'}
            />
            <p className="text-sm text-muted-foreground">
              {app.appType === 'database'
                ? 'Port is managed by the template'
                : 'The port your application runs on'}
            </p>
          </div>

          <div className="space-y-2">
            <Label htmlFor="rootDirectory">Root Directory</Label>
            <Input
              id="rootDirectory"
              placeholder="/"
              value={rootDirectory}
              onChange={(e) => setRootDirectory(e.target.value)}
              disabled={app.appType === 'database'}
            />
            <p className="text-sm text-muted-foreground">
              {app.appType === 'database'
                ? 'Not applicable for database apps'
                : 'The root directory of your application'}
            </p>
          </div>

          <div className="space-y-2">
            <Label htmlFor="cpuLimit">CPU Limit (cores)</Label>
            <Input
              id="cpuLimit"
              type="number"
              placeholder="e.g., 0.5 or 2"
              value={cpuLimit}
              onChange={(e) => setCpuLimit(e.target.value)}
              min="0"
              step="0.1"
            />
            <p className="text-sm text-muted-foreground">
              Maximum CPU cores (leave empty for no limit)
            </p>
          </div>

          <div className="space-y-2">
            <Label htmlFor="memoryLimit">Memory Limit (MB)</Label>
            <Input
              id="memoryLimit"
              type="number"
              placeholder="e.g., 512 or 1024"
              value={memoryLimit}
              onChange={(e) => setMemoryLimit(e.target.value)}
              min="0"
            />
            <p className="text-sm text-muted-foreground">
              Maximum memory in MB (leave empty for no limit)
            </p>
          </div>

          <div className="space-y-2">
            <Label htmlFor="restartPolicy">Restart Policy</Label>
            <select
              id="restartPolicy"
              value={restartPolicy}
              onChange={(e) => setRestartPolicy(e.target.value as RestartPolicy)}
              className="w-full bg-background border rounded-md px-3 py-2"
            >
              <option value="no">No</option>
              <option value="always">Always</option>
              <option value="on-failure">On Failure</option>
              <option value="unless-stopped">Unless Stopped</option>
            </select>
            <p className="text-sm text-muted-foreground">
              When to restart the container
            </p>
          </div>

          <div className="space-y-2">
            <Label htmlFor="deploymentStrategy">Deployment Strategy</Label>
            <select
              id="deploymentStrategy"
              value={deploymentStrategy}
              onChange={(e) => setDeploymentStrategy(e.target.value)}
              className="w-full bg-background border rounded-md px-3 py-2"
              disabled={app.appType === 'database'}
            >
              <option value="auto">Automatic</option>
              <option value="manual">Manual</option>
            </select>
            <p className="text-sm text-muted-foreground">
              {app.appType === 'database'
                ? 'Deployment strategy managed by template'
                : 'Auto: Deploy on every push. Manual: Deploy only when triggered manually'}
            </p>
          </div>

          <div className="space-y-2">
            <Label htmlFor="dockerfilePath">Dockerfile Path</Label>
            <Input
              id="dockerfilePath"
              placeholder="Dockerfile"
              value={dockerfilePath}
              onChange={(e) => setDockerfilePath(e.target.value)}
              disabled={app.appType === 'database'}
            />
            <p className="text-sm text-muted-foreground">
              {app.appType === 'database'
                ? 'Not applicable for database apps'
                : 'Path to your Dockerfile (optional)'}
            </p>
          </div>

          <div className="space-y-2">
            <Label htmlFor="healthcheckPath">Health Check Path</Label>
            <Input
              id="healthcheckPath"
              placeholder="/health"
              value={healthcheckPath}
              onChange={(e) => setHealthcheckPath(e.target.value)}
              disabled={app.appType === 'database'}
            />
            <p className="text-sm text-muted-foreground">
              {app.appType === 'database'
                ? 'Health checks managed by template'
                : 'Path for health checks (optional)'}
            </p>
          </div>

          <div className="space-y-2 md:col-span-2">
            <Label htmlFor="buildCommand">Build Command</Label>
            <Input
              id="buildCommand"
              placeholder="npm run build"
              value={buildCommand}
              onChange={(e) => setBuildCommand(e.target.value)}
              disabled={app.appType === 'database'}
            />
            <p className="text-sm text-muted-foreground">
              {app.appType === 'database'
                ? 'Not applicable for database apps'
                : 'Command to build your application (optional)'}
            </p>
          </div>

          <div className="space-y-2 md:col-span-2">
            <Label htmlFor="startCommand">Start Command</Label>
            <Input
              id="startCommand"
              placeholder="npm start"
              value={startCommand}
              onChange={(e) => setStartCommand(e.target.value)}
              disabled={app.appType === 'database'}
            />
            <p className="text-sm text-muted-foreground">
              {app.appType === 'database'
                ? 'Start command managed by template'
                : 'Command to start your application (optional)'}
            </p>
          </div>
        </div>

      </CardContent>
    </Card>
  );
};
