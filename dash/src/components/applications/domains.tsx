import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Trash2, Plus, Pencil, X, Check, ExternalLink, ChevronDown, ChevronUp } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { toast } from "sonner";
import { useDomains } from "@/hooks";
import { DNSValidation } from "./dns-validation";
import type { Domain } from "@/types";

interface DomainsProps {
  appId: number;
}

export const Domains = ({ appId }: DomainsProps) => {
  const { domains, loading, createDomain, updateDomain, deleteDomain, updateDomainInState } = useDomains({
    appId,
    autoFetch: true
  });

  const [newDomain, setNewDomain] = useState("");
  const [editingId, setEditingId] = useState<number | null>(null);
  const [editDomain, setEditDomain] = useState("");
  const [showAddForm, setShowAddForm] = useState(false);
  const [expandedDomain, setExpandedDomain] = useState<number | null>(null);

  const handleAdd = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newDomain.trim()) {
      toast.error("Domain is required");
      return;
    }

    const result = await createDomain(newDomain.trim());
    if (result) {
      setNewDomain("");
      setShowAddForm(false);
    }
  };

  const handleUpdate = async (id: number) => {
    if (!editDomain.trim()) {
      toast.error("Domain is required");
      return;
    }

    const result = await updateDomain(id, editDomain.trim());
    if (result) {
      setEditingId(null);
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm("Are you sure you want to delete this domain?")) {
      return;
    }
    await deleteDomain(id);
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

  const getDnsStatusBadge = (domain: Domain) => {
    if (domain.dnsConfigured) {
      return (
        <Badge variant="default" className="text-xs">
          DNS Configured
        </Badge>
      );
    }
    return (
      <Badge variant="secondary" className="text-xs">
        DNS Pending
      </Badge>
    );
  };

  const toggleDomainExpansion = (domainId: number) => {
    setExpandedDomain(expandedDomain === domainId ? null : domainId);
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
              <div key={domain.id} className="border rounded-lg bg-card">
                {editingId === domain.id ? (
                  <div className="p-4 space-y-4">
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
                  <>
                    <div className="p-4 flex items-center justify-between">
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
                        <div className="flex items-center gap-2">
                          <Badge variant={getSslStatusColor(domain.sslStatus)}>
                            SSL: {domain.sslStatus}
                          </Badge>
                          {getDnsStatusBadge(domain)}
                        </div>
                      </div>
                      <div className="flex gap-2">
                        <Button
                          size="sm"
                          variant="ghost"
                          onClick={() => toggleDomainExpansion(domain.id)}
                        >
                          {expandedDomain === domain.id ? (
                            <ChevronUp className="h-4 w-4" />
                          ) : (
                            <ChevronDown className="h-4 w-4" />
                          )}
                        </Button>
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
                    {expandedDomain === domain.id && (
                      <div className="px-4 pb-4">
                        <DNSValidation
                          domain={domain}
                          onVerified={(updatedDomain) => {
                            updateDomainInState(updatedDomain);
                          }}
                        />
                      </div>
                    )}
                  </>
                )}
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
};
