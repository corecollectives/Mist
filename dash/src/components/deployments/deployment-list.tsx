import { useEffect, useState } from "react"
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { toast } from "sonner"
import { DeploymentMonitor } from "@/components/deployments"
import type { Deployment, App } from "@/types"
import { Loader2, Clock, CheckCircle2, XCircle, PlayCircle, AlertCircle } from "lucide-react"
import { deploymentsService } from "@/services"

export const DeploymentsTab = ({ appId, app }: { appId: number; app?: App }) => {
  const [deployments, setDeployments] = useState<Deployment[]>([])
  const [loading, setLoading] = useState(true)
  const [deploying, setDeploying] = useState(false)
  const [selectedDeployment, setSelectedDeployment] = useState<number | null>(null)

  const fetchDeployments = async () => {
    try {
      setLoading(true)
      const data = await deploymentsService.getByAppId(appId)
      setDeployments(data || [])
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to fetch deployments")
    } finally {
      setLoading(false)
    }
  }

  const handleDeploy = async () => {
    try {
      setDeploying(true)

      const deployment = await deploymentsService.create(appId)
      toast.success('Deployment started successfully')

      // Open the monitor immediately
      setSelectedDeployment(deployment.id)

      // Refresh deployments list
      await fetchDeployments()
    } catch (error) {
      console.error('Deployment error:', error)
      toast.error(error instanceof Error ? error.message : 'Failed to start deployment')
    } finally {
      setDeploying(false)
    }
  }

  const handleDeploymentComplete = () => {
    toast.success('Deployment completed successfully!')
    fetchDeployments()
  }

  useEffect(() => {
    fetchDeployments()

    // Auto-refresh deployments every 10 seconds to catch updates
    const interval = setInterval(fetchDeployments, 10000)
    return () => clearInterval(interval)
  }, [appId])

  // Helper to get status icon and color
  const getStatusBadge = (deployment: Deployment) => {
    const { status, stage } = deployment

    switch (status) {
      case 'success':
        return (
          <Badge className="bg-green-500 text-white flex items-center gap-1.5">
            <CheckCircle2 className="h-3 w-3" />
            Success
          </Badge>
        )
      case 'failed':
        return (
          <Badge variant="destructive" className="flex items-center gap-1.5">
            <XCircle className="h-3 w-3" />
            Failed
          </Badge>
        )
      case 'building':
      case 'deploying':
      case 'cloning':
        return (
          <Badge className="bg-blue-500 text-white flex items-center gap-1.5 animate-pulse">
            <Loader2 className="h-3 w-3 animate-spin" />
            {stage.charAt(0).toUpperCase() + stage.slice(1)}
          </Badge>
        )
      case 'pending':
        return (
          <Badge variant="outline" className="flex items-center gap-1.5">
            <AlertCircle className="h-3 w-3" />
            Pending
          </Badge>
        )
      default:
        return <Badge variant="outline">{status}</Badge>
    }
  }

  return (
    <>
      {/* Deployment Monitor */}
      {selectedDeployment && (
        <DeploymentMonitor
          deploymentId={selectedDeployment}
          open={!!selectedDeployment}
          onClose={() => setSelectedDeployment(null)}
          onComplete={handleDeploymentComplete}
        />
      )}

      {/* Deployments Card */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <CardTitle>Deployments</CardTitle>
          <Button
            onClick={handleDeploy}
            disabled={deploying}
            className="flex items-center gap-2"
          >
            {deploying ? (
              <>
                <Loader2 className="h-4 w-4 animate-spin" />
                Deploying...
              </>
            ) : (
              <>
                <PlayCircle className="h-4 w-4" />
                Deploy Now
              </>
            )}
          </Button>
        </CardHeader>

        <CardContent className="space-y-4">
          {loading && (
            <div className="flex items-center justify-center py-8 text-muted-foreground">
              <Loader2 className="h-6 w-6 animate-spin mr-2" />
              Loading deployments...
            </div>
          )}

          {!loading && deployments.length === 0 && (
            <div className="text-center py-12">
              <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-muted mb-4">
                <PlayCircle className="h-8 w-8 text-muted-foreground" />
              </div>
              <p className="text-muted-foreground mb-2">No deployments yet</p>
              <p className="text-sm text-muted-foreground">
                Click "Deploy Now" to create your first deployment
              </p>
            </div>
          )}

          {!loading && deployments.length > 0 && (
            <div className="space-y-3">
              {deployments.map((d) => (
                <div
                  key={d.id}
                  className="flex items-start justify-between bg-muted/20 p-4 rounded-lg border hover:bg-muted/30 transition-colors"
                >
                  <div className="flex-1 space-y-2">
                    <div className="flex items-center gap-3 flex-wrap">
                      {getStatusBadge(d)}

                      <span className="text-xs text-muted-foreground font-mono">
                        #{d.id}
                      </span>

                      {/* Progress indicator for in-progress deployments */}
                      {d.status !== 'success' && d.status !== 'failed' && d.progress > 0 && (
                        <div className="flex items-center gap-2">
                          <div className="w-24 bg-muted rounded-full h-1.5 overflow-hidden">
                            <div
                              className="bg-primary h-full transition-all duration-300"
                              style={{ width: `${d.progress}%` }}
                            />
                          </div>
                          <span className="text-xs text-muted-foreground">
                            {d.progress}%
                          </span>
                        </div>
                      )}
                    </div>

                    <div className="space-y-1">
                      {app?.appType === 'database' ? (
                        <p className="font-mono text-sm">
                          <span className="text-primary">Version: {d.commit_hash}</span>
                          {d.commit_message && (
                            <>
                              {' – '}
                              {d.commit_message}
                            </>
                          )}
                        </p>
                      ) : (
                        <p className="font-mono text-sm">
                          <span className="text-primary">{d.commit_hash.slice(0, 7)}</span>
                          {' – '}
                          {d.commit_message}
                        </p>
                      )}

                      {d.error_message && (
                        <p className="text-xs text-red-500 flex items-start gap-1">
                          <XCircle className="h-3 w-3 mt-0.5 shrink-0" />
                          {d.error_message}
                        </p>
                      )}
                    </div>

                    <div className="flex items-center gap-4 text-xs text-muted-foreground">
                      <span className="flex items-center gap-1">
                        <Clock className="h-3 w-3" />
                        {new Date(d.created_at).toLocaleString()}
                      </span>

                      {d.duration && (
                        <span>
                          Duration: {d.duration}s
                        </span>
                      )}
                    </div>
                  </div>

                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setSelectedDeployment(d.id)}
                    className="ml-4"
                  >
                    View Logs
                  </Button>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </>
  )
}
