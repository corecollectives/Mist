import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { toast } from "sonner";
import { applicationsService } from "@/services";
import type { App } from "@/types";

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
  const [saving, setSaving] = useState(false);

  const handleSave = async () => {
    try {
      setSaving(true);

      const updates: any = {
        rootDirectory,
        dockerfilePath: dockerfilePath || null,
        buildCommand: buildCommand || null,
        startCommand: startCommand || null,
        healthcheckPath: healthcheckPath || null,
      };

      if (port) {
        const portNum = parseInt(port);
        if (isNaN(portNum) || portNum < 1 || portNum > 65535) {
          toast.error("Port must be a number between 1 and 65535");
          return;
        }
        updates.port = portNum;
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
      <CardHeader>
        <CardTitle>Application Settings</CardTitle>
        <CardDescription>
          Configure your application settings. Changes will be applied on next deployment.
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
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
            />
            <p className="text-sm text-muted-foreground">
              The port your application runs on
            </p>
          </div>

          <div className="space-y-2">
            <Label htmlFor="rootDirectory">Root Directory</Label>
            <Input
              id="rootDirectory"
              placeholder="/"
              value={rootDirectory}
              onChange={(e) => setRootDirectory(e.target.value)}
            />
            <p className="text-sm text-muted-foreground">
              The root directory of your application
            </p>
          </div>

          <div className="space-y-2">
            <Label htmlFor="dockerfilePath">Dockerfile Path</Label>
            <Input
              id="dockerfilePath"
              placeholder="Dockerfile"
              value={dockerfilePath}
              onChange={(e) => setDockerfilePath(e.target.value)}
            />
            <p className="text-sm text-muted-foreground">
              Path to your Dockerfile (optional)
            </p>
          </div>

          <div className="space-y-2">
            <Label htmlFor="healthcheckPath">Health Check Path</Label>
            <Input
              id="healthcheckPath"
              placeholder="/health"
              value={healthcheckPath}
              onChange={(e) => setHealthcheckPath(e.target.value)}
            />
            <p className="text-sm text-muted-foreground">
              Path for health checks (optional)
            </p>
          </div>

          <div className="space-y-2 md:col-span-2">
            <Label htmlFor="buildCommand">Build Command</Label>
            <Input
              id="buildCommand"
              placeholder="npm run build"
              value={buildCommand}
              onChange={(e) => setBuildCommand(e.target.value)}
            />
            <p className="text-sm text-muted-foreground">
              Command to build your application (optional)
            </p>
          </div>

          <div className="space-y-2 md:col-span-2">
            <Label htmlFor="startCommand">Start Command</Label>
            <Input
              id="startCommand"
              placeholder="npm start"
              value={startCommand}
              onChange={(e) => setStartCommand(e.target.value)}
            />
            <p className="text-sm text-muted-foreground">
              Command to start your application (optional)
            </p>
          </div>
        </div>

        <div className="flex justify-end">
          <Button onClick={handleSave} disabled={saving}>
            {saving ? "Saving..." : "Save Settings"}
          </Button>
        </div>
      </CardContent>
    </Card>
  );
};
