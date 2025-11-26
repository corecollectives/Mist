import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Trash2, Plus, Pencil, X, Check } from "lucide-react";
import { toast } from "sonner";
import { applicationsService } from "@/services";
import type { EnvVariable } from "@/types";

interface EnvironmentVariablesProps {
  appId: number;
}

export const EnvironmentVariables = ({ appId }: EnvironmentVariablesProps) => {
  const [envVars, setEnvVars] = useState<EnvVariable[]>([]);
  const [loading, setLoading] = useState(true);
  const [newKey, setNewKey] = useState("");
  const [newValue, setNewValue] = useState("");
  const [editingId, setEditingId] = useState<number | null>(null);
  const [editKey, setEditKey] = useState("");
  const [editValue, setEditValue] = useState("");
  const [showAddForm, setShowAddForm] = useState(false);

  const fetchEnvVars = async () => {
    try {
      setLoading(true);
      const data = await applicationsService.getEnvVariables(appId);
      setEnvVars(data);
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to fetch environment variables");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchEnvVars();
  }, [appId]);

  const handleAdd = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newKey.trim()) {
      toast.error("Key is required");
      return;
    }

    try {
      await applicationsService.createEnvVariable({
        appId,
        key: newKey.trim(),
        value: newValue,
      });
      toast.success("Environment variable added");
      setNewKey("");
      setNewValue("");
      setShowAddForm(false);
      await fetchEnvVars();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to add environment variable");
    }
  };

  const handleUpdate = async (id: number) => {
    if (!editKey.trim()) {
      toast.error("Key is required");
      return;
    }

    try {
      await applicationsService.updateEnvVariable({
        id,
        key: editKey.trim(),
        value: editValue,
      });
      toast.success("Environment variable updated");
      setEditingId(null);
      await fetchEnvVars();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to update environment variable");
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm("Are you sure you want to delete this environment variable?")) {
      return;
    }

    try {
      await applicationsService.deleteEnvVariable(id);
      toast.success("Environment variable deleted");
      await fetchEnvVars();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to delete environment variable");
    }
  };

  const startEdit = (env: EnvVariable) => {
    setEditingId(env.id);
    setEditKey(env.key);
    setEditValue(env.value);
  };

  const cancelEdit = () => {
    setEditingId(null);
    setEditKey("");
    setEditValue("");
  };

  if (loading) {
    return <div className="text-muted-foreground">Loading environment variables...</div>;
  }

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle>Environment Variables</CardTitle>
            <CardDescription>
              Manage environment variables for your application. Changes will be applied on next deployment.
            </CardDescription>
          </div>
          {!showAddForm && (
            <Button onClick={() => setShowAddForm(true)} size="sm">
              <Plus className="h-4 w-4 mr-2" />
              Add Variable
            </Button>
          )}
        </div>
      </CardHeader>
      <CardContent className="space-y-4">
        {showAddForm && (
          <form onSubmit={handleAdd} className="space-y-4 p-4 border rounded-lg bg-muted/50">
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="new-key">Key</Label>
                <Input
                  id="new-key"
                  placeholder="API_KEY"
                  value={newKey}
                  onChange={(e) => setNewKey(e.target.value)}
                  autoFocus
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="new-value">Value</Label>
                <Input
                  id="new-value"
                  placeholder="your-api-key-value"
                  value={newValue}
                  onChange={(e) => setNewValue(e.target.value)}
                  type="password"
                />
              </div>
            </div>
            <div className="flex gap-2">
              <Button type="submit" size="sm">
                <Check className="h-4 w-4 mr-2" />
                Add
              </Button>
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={() => {
                  setShowAddForm(false);
                  setNewKey("");
                  setNewValue("");
                }}
              >
                <X className="h-4 w-4 mr-2" />
                Cancel
              </Button>
            </div>
          </form>
        )}

        {envVars.length === 0 && !showAddForm ? (
          <p className="text-muted-foreground text-center py-8">
            No environment variables added yet. Click "Add Variable" to get started.
          </p>
        ) : (
          <div className="space-y-2">
            {envVars.map((env) => (
              <div key={env.id} className="p-4 border rounded-lg bg-card">
                {editingId === env.id ? (
                  <div className="space-y-4">
                    <div className="grid grid-cols-2 gap-4">
                      <div className="space-y-2">
                        <Label htmlFor={`edit-key-${env.id}`}>Key</Label>
                        <Input
                          id={`edit-key-${env.id}`}
                          value={editKey}
                          onChange={(e) => setEditKey(e.target.value)}
                        />
                      </div>
                      <div className="space-y-2">
                        <Label htmlFor={`edit-value-${env.id}`}>Value</Label>
                        <Input
                          id={`edit-value-${env.id}`}
                          value={editValue}
                          onChange={(e) => setEditValue(e.target.value)}
                          type="password"
                        />
                      </div>
                    </div>
                    <div className="flex gap-2">
                      <Button size="sm" onClick={() => handleUpdate(env.id)}>
                        <Check className="h-4 w-4 mr-2" />
                        Save
                      </Button>
                      <Button size="sm" variant="outline" onClick={cancelEdit}>
                        <X className="h-4 w-4 mr-2" />
                        Cancel
                      </Button>
                    </div>
                  </div>
                ) : (
                  <div className="flex items-center justify-between">
                    <div className="flex-1 font-mono">
                      <span className="font-semibold">{env.key}</span>
                      <span className="text-muted-foreground ml-2">=</span>
                      <span className="ml-2 text-muted-foreground">••••••••</span>
                    </div>
                    <div className="flex gap-2">
                      <Button
                        size="sm"
                        variant="ghost"
                        onClick={() => startEdit(env)}
                      >
                        <Pencil className="h-4 w-4" />
                      </Button>
                      <Button
                        size="sm"
                        variant="ghost"
                        onClick={() => handleDelete(env.id)}
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </div>
                  </div>
                )}
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
};
