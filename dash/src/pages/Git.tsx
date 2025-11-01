import { useEffect, useState } from "react"
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
  CardFooter,
} from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { ExternalLink, Github, Gitlab, GitFork, GitMerge } from "lucide-react"
import { useAuth } from "@/context/AuthContext"
import type { GitHubApp } from "@/lib/types"
import { cn } from "@/lib/utils"
import { Alert, AlertDescription } from "@/components/ui/alert"
import Loading from "@/components/Loading"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/components/ui/dialog"

export function GitPage() {
  const [loading, setLoading] = useState(true)
  const [app, setApp] = useState<GitHubApp | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [isInstalled, setIsInstalled] = useState(false)
  const [open, setOpen] = useState(false) // 👈 new state for dialog
  const { user } = useAuth()

  const generateState = (appId: number, userId: number) => {
    const payload = { appId, userId }
    return btoa(JSON.stringify(payload))
  }

  const fetchApp = async () => {
    try {
      const res = await fetch("/api/github/app")
      const data = await res.json()
      if (data.success) {
        setApp(data.data.app)
        setIsInstalled(data.data.isInstalled)
      } else {
        setError(data.error || "Failed to load GitHub App info")
      }
    } catch (err) {
      console.error(err)
      setError("Failed to load GitHub App info")
    } finally {
      setLoading(false)
    }
  }

  const fetchRepos = async () => {
    try {
      const res = await fetch("/api/github/app/repositories")
      const data = await res.json()
      console.log("Fetched repos:", data)
    } catch (err) {
      console.error("Failed to fetch repos:", err)
    }
  }

  useEffect(() => {
    fetchApp()
    fetchRepos()
  }, [])

  const handleCreateApp = () => {
    // Close modal and redirect
    setOpen(false)
    window.open("/api/github/app/create", "_blank")
  }

  if (loading)
    return (
      <div className="flex h-screen w-full items-center justify-center">
        <Loading />
      </div>
    )

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <div className="flex items-center justify-between py-6 border-b border-border flex-shrink-0">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-foreground">
            Git Integrations
          </h1>
          <p className="text-muted-foreground mt-1">
            Connect your Git providers to enable automatic deployments.
          </p>
        </div>
      </div>

      {/* Error */}
      {error && (
        <div className="mt-4">
          <Alert variant="destructive">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        </div>
      )}

      {/* Main Grid */}
      <div className="grid gap-4 py-6 md:grid-cols-2 lg:grid-cols-3">
        {/* GitHub Card */}
        <Card className="border-border bg-card hover:border-primary transition-colors">
          <CardHeader className="pb-2 flex flex-row items-center gap-2">
            <Github className="w-5 h-5 text-muted-foreground" />
            <CardTitle className="text-lg font-semibold">GitHub</CardTitle>
          </CardHeader>

          <CardContent>
            {app ? (
              <div className="space-y-3">
                <div>
                  <p className="font-medium text-foreground">{app.name}</p>
                  <p className="text-sm text-muted-foreground mt-1">
                    App ID: {app.app_id}
                  </p>
                </div>

                <div className="flex flex-wrap gap-2 pt-2">
                  <Button
                    onClick={() => {
                      const state = generateState(
                        app.app_id,
                        parseInt(user!.id.toString())
                      )
                      window.open(
                        `https://github.com/apps/${app.slug}/installations/new?state=${state}`
                      )
                    }}
                    disabled={isInstalled}
                    className="transition-colors disabled:cursor-not-allowed"
                  >
                    {isInstalled ? "Installed" : "Install App"}
                    <ExternalLink className="w-4 h-4 ml-2" />
                  </Button>

                  <Button
                    variant="outline"
                    onClick={() =>
                      window.open(
                        `https://github.com/settings/apps/${app.slug}`,
                        "_blank"
                      )
                    }
                  >
                    Manage
                  </Button>
                </div>
              </div>
            ) : (
              <div className="space-y-3">
                <p className="text-sm text-muted-foreground">
                  No GitHub App connected yet. Create one to enable Git
                  deployments.
                </p>
                <Button
                  onClick={() => setOpen(true)} // 👈 open popup first
                  size="default"
                  className="transition-colors"
                >
                  Create GitHub App
                </Button>
              </div>
            )}
          </CardContent>

          {app && (
            <CardFooter className="text-xs text-muted-foreground border-t border-border pt-2">
              Created at: {new Date(app.created_at).toLocaleString()}
            </CardFooter>
          )}
        </Card>

        {/* Other Providers */}
        {[
          { name: "GitLab", icon: Gitlab },
          { name: "Bitbucket", icon: GitFork },
          { name: "Gitea", icon: GitMerge },
        ].map(({ name, icon: Icon }) => (
          <Card
            key={name}
            className={cn(
              "h-full flex flex-col items-center justify-between border border-dashed border-border bg-card hover:border-primary/30 cursor-not-allowed opacity-60 transition-colors"
            )}
          >
            <CardHeader className="flex flex-col items-center space-y-3 pb-4">
              <Icon className="w-6 h-6 text-muted-foreground" />
              <CardTitle className="text-base font-medium">{name}</CardTitle>
            </CardHeader>
            <CardContent className="pb-6">
              <Badge variant="secondary">Coming Soon</Badge>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Confirmation Popup */}
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Create GitHub App</DialogTitle>
            <DialogDescription>
              This will create a new GitHub App in your account with permissions
              for:
              <ul className="list-disc list-inside mt-2 text-muted-foreground">
                <li>Accessing your repositories</li>
                <li>Receiving push & deployment events</li>
                <li>Managing webhooks for automation</li>
                <li>Other users will be able to use this app for deployments</li>
              </ul>
              <p className="mt-2">
                You’ll be redirected to GitHub to complete the process.
              </p>
            </DialogDescription>
          </DialogHeader>

          <DialogFooter className="flex justify-end space-x-2">
            <Button variant="outline" onClick={() => setOpen(false)}>
              Cancel
            </Button>
            <Button onClick={handleCreateApp}>Continue</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
