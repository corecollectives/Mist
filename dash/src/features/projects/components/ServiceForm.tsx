import { useState } from "react";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Alert } from "@/components/ui/alert";
import { Info } from "lucide-react";
import type { CreateAppRequest } from "@/types/app";

interface ServiceFormProps {
  projectId: number;
  onSubmit: (data: CreateAppRequest) => void;
  onBack: () => void;
}

export function ServiceForm({ projectId, onSubmit, onBack }: ServiceFormProps) {
  const [formData, setFormData] = useState({
    name: "",
    description: "",
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit({
      projectId,
      appType: "service",
      name: formData.name,
      description: formData.description || undefined,
      port: 3000, // Internal port, doesn't matter for services
    });
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div>
        <h3 className="text-lg font-semibold mb-2">Create Background Service</h3>
        <p className="text-sm text-muted-foreground">
          Workers, bots, and processes that run internally without external access
        </p>
      </div>

      <Alert>
        <Info className="h-4 w-4" />
        <div className="ml-2">
          <p className="text-sm">
            Background services don't need port configuration - they run internally without
            external HTTP access.
          </p>
        </div>
      </Alert>

      <div className="space-y-4">
        <div>
          <Label htmlFor="name">Service Name *</Label>
          <Input
            id="name"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            placeholder="my-discord-bot"
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
            placeholder="Brief description of your service"
            className="mt-1"
          />
        </div>
      </div>

      <div className="flex justify-between pt-4">
        <Button type="button" variant="outline" onClick={onBack}>
          Back
        </Button>
        <Button type="submit">Create Service</Button>
      </div>
    </form>
  );
}
