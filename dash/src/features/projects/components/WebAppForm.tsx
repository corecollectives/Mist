import { useState } from "react";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import type { CreateAppRequest } from "@/types/app";

interface WebAppFormProps {
  projectId: number;
  onSubmit: (data: CreateAppRequest) => void;
  onBack: () => void;
}

export function WebAppForm({ projectId, onSubmit, onBack }: WebAppFormProps) {
  const [formData, setFormData] = useState({
    name: "",
    description: "",
    port: 3000,
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit({
      projectId,
      appType: "web",
      name: formData.name,
      description: formData.description || undefined,
      port: formData.port,
    });
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div>
        <h3 className="text-lg font-semibold mb-2">Create Web Application</h3>
        <p className="text-sm text-muted-foreground">
          HTTP servers and web applications that need external access
        </p>
      </div>

      <div className="space-y-4">
        <div>
          <Label htmlFor="name">Application Name *</Label>
          <Input
            id="name"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            placeholder="my-web-app"
            required
            className="mt-1"
          />
          <p className="text-xs text-muted-foreground mt-1">
            Lowercase letters, numbers, and hyphens only
          </p>
        </div>

        <div>
          <Label htmlFor="description">Description</Label>
          <Textarea
            id="description"
            value={formData.description}
            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
            placeholder="Brief description of your application"
            className="mt-1"
          />
        </div>

        <div>
          <Label htmlFor="port">Port *</Label>
          <Input
            id="port"
            type="number"
            value={formData.port}
            onChange={(e) => setFormData({ ...formData, port: parseInt(e.target.value) })}
            placeholder="3000"
            required
            min={1}
            max={65535}
            className="mt-1"
          />
          <p className="text-xs text-muted-foreground mt-1">
            The port your application listens on (default: 3000)
          </p>
        </div>
      </div>

      <div className="flex justify-between pt-4">
        <Button type="button" variant="outline" onClick={onBack}>
          Back
        </Button>
        <Button type="submit">Create Application</Button>
      </div>
    </form>
  );
}
