import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Trash2, Plus, Pencil, X, Check, ExternalLink } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { toast } from "sonner";
import { applicationsService } from "@/services";
import type { Domain } from "@/types";

interface DomainsProps {
  appId: number;
}

export const Domains = ({ appId }: DomainsProps) => {
  const [domains, setDomains] = useState<Domain[]>([]);
  const [loading, setLoading] = useState(true);
  const [newDomain, setNewDomain] = useState("");
  const [editingId, setEditingId] = useState<number | null>(null);
  const [editDomain, setEditDomain] = useState("");
  const [showAddForm, setShowAddForm] = useState(false);

  const fetchDomains = async () => {
    try {
      setLoading(true);
      const data = await applicationsService.getDomains(appId);
      setDomains(data);
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to fetch domains");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchDomains();
  }, [appId]);

  const handleAdd = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newDomain.trim()) {
      toast.error("Domain is required");
      return;
    }

    try {
      await applicationsService.createDomain({
        appId,
        domain: newDomain.trim(),
      });
      toast.success("Domain added");
      setNewDomain("");
      setShowAddForm(false);
      await fetchDomains();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to add domain");
    }
  };

  const handleUpdate = async (id: number) => {
    if (!editDomain.trim()) {
      toast.error("Domain is required");
      return;
    }

    try {
      await applicationsService.updateDomain({
        id,
        domain: editDomain.trim(),
      });
      toast.success("Domain updated");
      setEditingId(null);
      await fetchDomains();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to update domain");
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm("Are you sure you want to delete this domain?")) {
      return;
    }

    try {
      await applicationsService.deleteDomain(id);
      toast.success("Domain deleted");
      await fetchDomains();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to delete domain");
    }
  };

  const startEdit = (domain: Domain) => {
    setEditingId(domain.id);
    setEditDomain(domain.domain);
  };

  const cancelEdit = () => {
    setEditingId(null);
    setEditDomain("");
  };

  const getSslStatusColor = (status: string) => {
    switch (status) {
      case "active":
        return "default";
      case "pending":
        return "secondary";
      case "failed":
        return "destructive";
      default:
        return "outline";
    }
  };

  if (loading) {
    return <div className="text-muted-foreground">Loading domains...</div>;
  }

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle>Domains</CardTitle>
            <CardDescription>
              Manage custom domains for your application. Changes will be applied on next deployment.
            </CardDescription>
          </div>
          {!showAddForm && (
            <Button onClick={() => setShowAddForm(true)} size="sm">
              <Plus className="h-4 w-4 mr-2" />
              Add Domain
            </Button>
          )}
        </div>
      </CardHeader>
      <CardContent className="space-y-4">
        {showAddForm && (
          <form onSubmit={handleAdd} className="space-y-4 p-4 border rounded-lg bg-muted/50">
            <div className="space-y-2">
              <Label htmlFor="new-domain">Domain</Label>
              <Input
                id="new-domain"
                placeholder="example.com"
                value={newDomain}
                onChange={(e) => setNewDomain(e.target.value)}
                autoFocus
              />
              <p className="text-sm text-muted-foreground">
                Enter the domain without http:// or https://
              </p>
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
                  setNewDomain("");
                }}
              >
                <X className="h-4 w-4 mr-2" />
                Cancel
              </Button>
            </div>
          </form>
        )}

        {domains.length === 0 && !showAddForm ? (
          <p className="text-muted-foreground text-center py-8">
            No domains added yet. Click "Add Domain" to get started.
          </p>
        ) : (
          <div className="space-y-2">
            {domains.map((domain) => (
              <div key={domain.id} className="p-4 border rounded-lg bg-card">
                {editingId === domain.id ? (
                  <div className="space-y-4">
                    <div className="space-y-2">
                      <Label htmlFor={`edit-domain-${domain.id}`}>Domain</Label>
                      <Input
                        id={`edit-domain-${domain.id}`}
                        value={editDomain}
                        onChange={(e) => setEditDomain(e.target.value)}
                      />
                    </div>
                    <div className="flex gap-2">
                      <Button size="sm" onClick={() => handleUpdate(domain.id)}>
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
                    <div className="flex items-center gap-3 flex-1">
                      <a
                        href={`https://${domain.domain}`}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="font-mono font-semibold hover:underline flex items-center gap-1"
                      >
                        {domain.domain}
                        <ExternalLink className="h-3 w-3" />
                      </a>
                      <Badge variant={getSslStatusColor(domain.sslStatus)}>
                        {domain.sslStatus}
                      </Badge>
                    </div>
                    <div className="flex gap-2">
                      <Button
                        size="sm"
                        variant="ghost"
                        onClick={() => startEdit(domain)}
                      >
                        <Pencil className="h-4 w-4" />
                      </Button>
                      <Button
                        size="sm"
                        variant="ghost"
                        onClick={() => handleDelete(domain.id)}
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
