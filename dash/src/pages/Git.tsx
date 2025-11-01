import { useEffect, useState } from "react";
import { Card, CardHeader, CardTitle, CardContent, CardFooter } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { Badge } from "@/components/ui/badge";
import { ExternalLink } from "lucide-react";
import type { GitHubApp } from "@/lib/types";
import { toast } from "sonner";


export function GitPage() {
  const [loading, setLoading] = useState(true);
  const [app, setApp] = useState<GitHubApp | null>(null);
  const [error, setError] = useState<string | null>(null);

  const fetchApp = async () => {
    setLoading(true);
    try {
      const res = await fetch("/api/github/app");
      const data = await res.json();
      if (data.success) {
        if (data.data) {
          setApp(data.data);
        }
      }
      else {
        setError(data.error || null);
        toast.error(data.error || "An error occurred while fetching the app info");
      }

    } catch (err) {
      console.error(err);
      setError("Failed to load app info");
    }
    setLoading(false);
  }

  useEffect(() => {
    fetchApp()
  }, []);

  const manifest = {
    name: "Mist PaaS",
    url: window.location.origin,
    hook_attributes: {
      url: "https://api.mist.local/api/github/webhook",
    },
    redirect_url: "http://localhost:8080/api/github/callback",
    callback_urls: [window.location.origin],
    public: false,
    default_permissions: {
      contents: "read",
      metadata: "read",
      pull_requests: "read",
      deployments: "read",
      administration: "write",
      repository_hooks: "write",
    },
    default_events: ["push", "pull_request", "deployment_status"],
  };

  const handleButtonClick = () => {
    const githubUrl = `https://github.com/settings/apps/new?state=abc123`;
    const form = document.createElement("form")
    form.method = "POST"
    form.action = githubUrl
    form.style.display = "none"

    const input = document.createElement("input")
    input.type = "hidden"
    input.name = "manifest"
    input.value = JSON.stringify(manifest)
    form.appendChild(input)

    document.body.appendChild(form)
    form.submit()

  }
  if (loading) {
    return (
      <div className="max-w-2xl mx-auto mt-16">
        <Card>
          <CardHeader>
            <CardTitle>GitHub Integration</CardTitle>
          </CardHeader>
          <CardContent>
            <Skeleton className="h-8 w-3/4 mb-4" />
            <Skeleton className="h-6 w-1/2" />
          </CardContent>
        </Card>
      </div>
    );
  }

  if (error) {
    return (
      <div className="max-w-2xl mx-auto mt-16">
        <Card>
          <CardHeader>
            <CardTitle>GitHub Integration</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-red-500">{error}</p>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="max-w-2xl mx-auto mt-16 space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>GitHub Integration</CardTitle>
        </CardHeader>
        <CardContent>
          {app ? (
            <div className="space-y-4">
              <div>
                <p className="text-sm text-muted-foreground">Connected GitHub App</p>
                <p className="text-lg font-semibold">{app.name}</p>
                <Badge variant="secondary" className="mt-2">
                  App ID: {app.app_id}
                </Badge>
              </div>

              <div className="pt-2">
                <p className="text-sm text-muted-foreground">
                  Installed App Slug:
                </p>
                <p className="font-medium">{app.slug}</p>
              </div>

              <div className="flex gap-3 pt-4">
                <Button
                  onClick={() =>
                    window.open(
                      `https://github.com/apps/${app.slug}/installations/new`,
                      "_blank"
                    )
                  }
                >
                  Install App
                  <ExternalLink className="w-4 h-4 ml-2" />
                </Button>

                <Button
                  variant="outline"
                  onClick={() =>
                    window.open(`https://github.com/settings/apps/${app.slug}`, "_blank")
                  }
                >
                  Manage App
                </Button>
              </div>
            </div>
          ) : (
            <div className="flex flex-col items-center gap-4">
              <p className="text-muted-foreground text-center">
                No GitHub App connected yet. Create one to enable Git deployments.
              </p>
              <Button
                size="lg"
                onClick={() => {
                  handleButtonClick()
                }}
              >
                Create GitHub App
              </Button>
            </div>
          )}
        </CardContent>

        {app && (
          <CardFooter className="text-sm text-muted-foreground">
            Created at: {new Date(app.created_at).toLocaleString()}
          </CardFooter>
        )}
      </Card>
    </div>
  );
}
