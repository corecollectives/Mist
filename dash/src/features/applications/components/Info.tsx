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
  GitCommit,
  Rocket,
  Server,
  Activity,
  Clock,
  Loader2,
} from "lucide-react"
import { Button } from "@/components/ui/button"
import { toast } from "sonner"
import { useState } from "react"
import type { App } from "@/types/app"
import { DeploymentMonitor } from "./DeploymentMonitor"

interface Props {
  app: App
  latestCommit?: {
    sha: string
    html_url: string
    author?: string
    timestamp?: string
  } | null
}

const InfoItem = ({
  icon: Icon,
  label,
  children,
  className = ""
}: {
  icon: any
  label: string
  children: React.ReactNode
  className?: string
}) => (
  <div className={`space-y-2 ${className}`}>
    <div className="flex items-center gap-2 text-sm font-medium text-muted-foreground">
      <Icon className="h-4 w-4" />
      <span>{label}</span>
    </div>
    <div className="pl-6">{children}</div>
  </div>
)

const getStatusConfig = (status: string) => {
  switch (status.toLowerCase()) {
    case "running":
      return {
        color: "text-green-600 dark:text-green-400",
        bgColor: "bg-green-500/10 border-green-500/20",
        icon: CheckCircle,
      }
    case "error":
    case "failed":
      return {
        color: "text-red-600 dark:text-red-400",
        bgColor: "bg-red-500/10 border-red-500/20",
        icon: XCircle,
      }
    case "deploying":
    case "building":
      return {
        color: "text-blue-600 dark:text-blue-400",
        bgColor: "bg-blue-500/10 border-blue-500/20",
        icon: Activity,
      }
    default:
      return {
        color: "text-gray-600 dark:text-gray-400",
        bgColor: "bg-gray-500/10 border-gray-500/20",
        icon: XCircle,
      }
  }
}

export const AppInfo = ({ app, latestCommit }: Props) => {
  const [deploying, setDeploying] = useState(false)
  const [logsOpen, setLogsOpen] = useState(false)
  const [deploymentId, setDeploymentId] = useState<number | null>(null)

  const statusConfig = getStatusConfig(app.status)

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

      if (!res.ok) {
        throw new Error(data.error || "Failed to create deployment")
      }

      toast.success("Deployment started!")

      setDeploymentId(data.id)
      setLogsOpen(true)

    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Deployment failed")
    } finally {
      setDeploying(false)
    }
  }

  return (
    <Card className="border-border/50">
      <CardHeader className="border-b border-border/50 bg-muted/30">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <CardTitle className="text-xl font-semibold">Application Overview</CardTitle>
          <Button
            onClick={handleDeploy}
            disabled={deploying}
            className="flex items-center gap-2 shadow-sm"
            size="sm"
          >
            {deploying ? (
              <>
                <Loader2 className="h-4 w-4 animate-spin" />
                Deploying...
              </>
            ) : (
              <>
                <Rocket className="h-4 w-4" />
                Deploy
              </>
            )}
          </Button>
        </div>
      </CardHeader>

      <CardContent className="p-6">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* Status */}
          <InfoItem icon={Activity} label="Status">
            <div className="flex items-center gap-2">
              <Badge
                variant="outline"
                className={`${statusConfig.bgColor} ${statusConfig.color} font-medium border`}
              >
                <span className="relative flex h-2 w-2 mr-2">
                  {app.status === "running" && (
                    <>
                      <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
                      <span className="relative inline-flex rounded-full h-2 w-2 bg-green-500"></span>
                    </>
                  )}
                  {app.status !== "running" && (
                    <span className="relative inline-flex rounded-full h-2 w-2 bg-current opacity-75"></span>
                  )}
                </span>
                {app.status.charAt(0).toUpperCase() + app.status.slice(1)}
              </Badge>
            </div>
          </InfoItem>

          {/* Deployment Strategy */}
          <InfoItem icon={Server} label="Deployment Strategy">
            <Badge variant="secondary" className="font-mono text-xs">
              {app.deploymentStrategy || "Not specified"}
            </Badge>
          </InfoItem>

          {/* Git Repository */}
          <InfoItem icon={Github} label="Repository">
            {app.gitRepository ? (
              <a
                href={`https://github.com/${app.gitRepository}`}
                target="_blank"
                rel="noopener noreferrer"
                className="font-mono text-sm text-primary hover:underline flex items-center gap-2 group"
              >
                <span className="truncate">{app.gitRepository}</span>
                <svg className="w-3 h-3 opacity-0 group-hover:opacity-100 transition-opacity" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                </svg>
              </a>
            ) : (
              <p className="text-muted-foreground text-sm">Not connected</p>
            )}
          </InfoItem>

          {/* Branch */}
          <InfoItem icon={GitBranch} label="Branch">
            <div className="flex items-center gap-2">
              <Badge variant="outline" className="font-mono text-xs">
                {app.gitBranch || "Not specified"}
              </Badge>
            </div>
          </InfoItem>

          {app.port !== 0 && (
            <InfoItem icon={Server} label="Port">
              <div className="font-mono text-sm font-medium">
                {app.port || <s className="text-red-500">"Not configured"</s>}
              </div>
            </InfoItem>
          )}

          {/* Root Directory */}
          {app.rootDirectory && (
            <InfoItem icon={Server} label="Root Directory">
              <div className="font-mono text-sm text-muted-foreground">
                {app.rootDirectory}
              </div>
            </InfoItem>
          )}

          {/* Build Command */}
          {app.buildCommand && (
            <InfoItem icon={Server} label="Build Command" className="md:col-span-2">
              <div className="font-mono text-sm bg-muted/50 border border-border/50 rounded px-3 py-2">
                {app.buildCommand}
              </div>
            </InfoItem>
          )}

          {/* Start Command */}
          {app.startCommand && (
            <InfoItem icon={Rocket} label="Start Command" className="md:col-span-2">
              <div className="font-mono text-sm bg-muted/50 border border-border/50 rounded px-3 py-2">
                {app.startCommand}
              </div>
            </InfoItem>
          )}

          {/* Latest Commit */}
          <InfoItem icon={GitCommit} label="Latest Commit" className="md:col-span-2">
            {latestCommit ? (
              <div className="space-y-2 p-3 rounded-lg bg-muted/50 border border-border/50">
                <a
                  href={latestCommit.html_url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="font-mono text-sm text-primary hover:underline inline-flex items-center gap-2 group"
                >
                  <span>{latestCommit.sha.slice(0, 7)}</span>
                  <svg className="w-3 h-3 opacity-0 group-hover:opacity-100 transition-opacity" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                  </svg>
                </a>
                <div className="flex items-center gap-4 text-xs text-muted-foreground">
                  {latestCommit.author && (
                    <span>by {latestCommit.author}</span>
                  )}
                  {latestCommit.timestamp && (
                    <span className="flex items-center gap-1">
                      <Clock className="w-3 h-3" />
                      {new Date(latestCommit.timestamp).toLocaleString()}
                    </span>
                  )}
                </div>
              </div>
            ) : (
              <p className="text-muted-foreground text-sm">
                No commit information available
              </p>
            )}
          </InfoItem>

          {/* Healthcheck */}
          {app.healthcheckPath && (
            <InfoItem icon={Activity} label="Health Check" className="md:col-span-2">
              <div className="flex items-center gap-4">
                <code className="font-mono text-sm bg-muted/50 border border-border/50 rounded px-3 py-1">
                  {app.healthcheckPath}
                </code>
                <span className="text-xs text-muted-foreground">
                  Interval: {app.healthcheckInterval}s
                </span>
              </div>
            </InfoItem>
          )}

          {/* Timestamps */}
          <InfoItem icon={Clock} label="Created" className="md:col-span-2">
            <div className="flex flex-wrap gap-4 text-sm text-muted-foreground">
              <div>
                <span className="font-medium text-foreground">Created:</span>{" "}
                {new Date(app.createdAt).toLocaleString()}
              </div>
              <div>
                <span className="font-medium text-foreground">Updated:</span>{" "}
                {new Date(app.updatedAt).toLocaleString()}
              </div>
            </div>
          </InfoItem>
        </div>
      </CardContent>

      {deploymentId !== null && (
        <DeploymentMonitor
          deploymentId={deploymentId}
          open={logsOpen}
          onClose={() => {
            setLogsOpen(false)
            setDeploymentId(null)
          }}
          onComplete={() => {
            toast.success("Deployment completed successfully!")
          }}
        />
      )}
    </Card>
  )
}
