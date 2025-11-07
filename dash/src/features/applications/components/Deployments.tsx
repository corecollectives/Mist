import { useEffect, useState } from "react"
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { toast } from "sonner"
import { DeploymentLogsOverlay } from "../components/DeploymentLogsOverlay"

interface Deployment {
  id: number
  app_id: number
  commit_hash: string
  commit_message: string
  status: string
  created_at: string
}

export const DeploymentsTab = ({ appId }: { appId: number }) => {
  const [deployments, setDeployments] = useState<Deployment[]>([])
  const [loading, setLoading] = useState(true)
  const [selectedDeployment, setSelectedDeployment] = useState<number | null>(null)

  const fetchDeployments = async () => {
    try {
      setLoading(true)
      const res = await fetch("/api/deployments/getByAppId", {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({ appId })
      })

      const data = await res.json()
      if (!data.success) throw new Error(data.error)

      setDeployments(data.data || [])
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to fetch deployments")
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchDeployments()
  }, [appId])

  return (
    <>
      {/* ✅ Logs Overlay */}
      {selectedDeployment && (
        <DeploymentLogsOverlay
          deploymentId={selectedDeployment}
          open={!!selectedDeployment}
          onClose={() => setSelectedDeployment(null)}
        />
      )}

      {/* ✅ Deployments List */}
      <Card>
        <CardHeader>
          <CardTitle>Deployments</CardTitle>
        </CardHeader>

        <CardContent className="space-y-4">
          {loading && <p className="text-muted-foreground">Loading deployments...</p>}

          {!loading && deployments.length === 0 && (
            <p className="text-muted-foreground">No deployments yet.</p>
          )}

          {!loading && deployments.length > 0 && (
            <div className="space-y-4">
              {deployments.map((d) => (
                <div
                  key={d.id}
                  className="flex items-start justify-between bg-muted/20 p-4 rounded-lg border"
                >
                  <div className="space-y-1">
                    <div className="flex items-center gap-2">
                      <Badge
                        variant={
                          d.status === "success"
                            ? "default"
                            : d.status === "failed"
                              ? "destructive"
                              : "outline"
                        }
                      >
                        {d.status}
                      </Badge>

                      <span className="text-xs text-muted-foreground">
                        #{d.id}
                      </span>
                    </div>

                    <p className="font-mono text-sm">
                      {d.commit_hash.slice(0, 7)} – {d.commit_message}
                    </p>

                    <p className="text-xs text-muted-foreground">
                      {new Date(d.created_at).toLocaleString()}
                    </p>
                  </div>

                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setSelectedDeployment(d.id)}
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
