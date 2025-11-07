import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import {
  CheckCircle,
  XCircle,
  Github,
  GitBranch,
  Globe,
  GitCommit,
  Rocket,
} from "lucide-react"
import { Button } from "@/components/ui/button"
import { toast } from "sonner"
import { useState } from "react"
import type { App } from "@/types/app"
import { DeploymentLogsOverlay } from "./components/DeploymentLogsOverlay"

interface Props {
  app: App
  latestCommit?: {
    sha: string
    html_url: string
    author?: string
    timestamp?: string
  } | null
}

export const AppInfo = ({ app, latestCommit }: Props) => {
  const [deploying, setDeploying] = useState(false)

  const [logsOpen, setLogsOpen] = useState(false)
  const [deploymentId, setDeploymentId] = useState<number | null>(null)

  const handleDeploy = async () => {
    try {
      setDeploying(true)

      const res = await fetch("/api/deployments/create", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ appId: app.id }),
      })

      const data = await res.json()

      toast.success("Deployment started!")

      // ✅ Open logs overlay
      setDeploymentId(data.id)
      setLogsOpen(true)

    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Deployment failed")
    } finally {
      setDeploying(false)
    }
  }

  return (
    <Card>
      <CardHeader>
        <div className="flex justify-between items-center">
          <CardTitle className="text-xl">Application Overview</CardTitle>

          {/* ✅ Deploy Button */}
          <Button
            onClick={handleDeploy}
            disabled={deploying}
            className="flex items-center gap-2"
          >
            <Rocket className="h-4 w-4" />
            {deploying ? "Deploying..." : "Deploy"}
          </Button>
        </div>
      </CardHeader>

      <CardContent className="grid grid-cols-1 md:grid-cols-2 gap-8">

        {/* ✅ STATUS */}
        <div className="space-y-1">
          <p className="text-sm text-muted-foreground flex items-center gap-2">
            Status
          </p>

          <div className="flex items-center gap-2 mt-1">
            {app.status === "running" ? (
              <CheckCircle className="h-4 w-4 text-green-500" />
            ) : (
              <XCircle className="h-4 w-4 text-red-500" />
            )}

            <Badge
              variant={app.status === "running" ? "default" : "destructive"}
              className="px-3 py-1 text-sm"
            >
              {app.status}
            </Badge>
          </div>
        </div>

        {/* ✅ GIT REPO */}
        <div className="space-y-1">
          <p className="text-sm text-muted-foreground flex items-center gap-2">
            <Github className="h-4 w-4" />
            Connected Repository
          </p>

          {app.gitRepository ? (
            <a
              href={`https://github.com/${app.gitRepository}`}
              target="_blank"
              className="font-mono text-sm mt-1 text-blue-500 underline block"
            >
              {app.gitRepository}
            </a>
          ) : (
            <p className="text-muted-foreground text-sm mt-1">Not connected</p>
          )}
        </div>

        {/* ✅ BRANCH */}
        <div className="space-y-1">
          <p className="text-sm text-muted-foreground flex items-center gap-2">
            <GitBranch className="h-4 w-4" />
            Branch
          </p>

          <p className="font-mono text-sm mt-1">
            {app.gitBranch || "Not specified"}
          </p>
        </div>

        {/* ✅ DOMAINS */}
        <div className="space-y-1">
          <p className="text-sm text-muted-foreground flex items-center gap-2">
            <Globe className="h-4 w-4" />
            Domains
          </p>

          <div className="flex flex-col mt-1 gap-1">

            <div className="flex items-center gap-2">
              <Badge variant="secondary" className="font-mono">
                app.example.com
              </Badge>
              <CheckCircle className="h-4 w-4 text-green-500" />
            </div>

            <div className="flex items-center gap-2">
              <Badge variant="secondary" className="font-mono">
                www.app.example.com
              </Badge>
              <XCircle className="h-4 w-4 text-red-500" />
            </div>

          </div>
        </div>

        {/* ✅ LATEST COMMIT */}
        <div className="md:col-span-2 space-y-1">
          <p className="text-sm text-muted-foreground flex items-center gap-2">
            <GitCommit className="h-4 w-4" />
            Latest Commit
          </p>

          {latestCommit ? (
            <div className="space-y-1 mt-1">
              <a
                href={latestCommit.html_url}
                target="_blank"
                className="font-mono text-sm text-blue-500 underline inline-block"
              >
                {latestCommit.sha.slice(0, 7)}
              </a>

              {latestCommit.author && (
                <p className="text-sm text-muted-foreground">
                  by {latestCommit.author}
                </p>
              )}

              {latestCommit.timestamp && (
                <p className="text-xs text-muted-foreground">
                  {new Date(latestCommit.timestamp).toLocaleString()}
                </p>
              )}
            </div>
          ) : (
            <p className="text-muted-foreground text-sm mt-1">
              No commit information available.
            </p>
          )}
        </div>
      </CardContent>
      {deploymentId !== null && (
        <DeploymentLogsOverlay
          deploymentId={deploymentId}
          open={logsOpen}
          onClose={() => setLogsOpen(false)}
        />
      )}

    </Card>
  )
}
