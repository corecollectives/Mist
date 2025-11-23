import { useState, useEffect } from "react"
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { applicationsService } from "@/services"
import type { ContainerStatus } from "@/types"
import {
  Activity,
  CheckCircle2,
  Clock,
  RefreshCw,
  Server,
  AlertCircle,
  Loader2,
  Play,
  ExternalLink,
} from "lucide-react"

interface AppStatsProps {
  appId: number
  appStatus: string
  previewUrl?: string
}

export const AppStats = ({ appId, previewUrl }: AppStatsProps) => {
  const [containerStatus, setContainerStatus] = useState<ContainerStatus | null>(null)
  const [loading, setLoading] = useState(false)
  const [lastUpdated, setLastUpdated] = useState<Date>(new Date())

  const fetchContainerStatus = async () => {
    try {
      setLoading(true)
      const data = await applicationsService.getContainerStatus(appId)
      setContainerStatus(data)
      setLastUpdated(new Date())
    } catch (error) {
      console.error("Failed to fetch container status:", error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchContainerStatus()
    
    // Refresh every 15 seconds
    const interval = setInterval(fetchContainerStatus, 15000)
    return () => clearInterval(interval)
  }, [appId])

  const getStateBadge = () => {
    if (!containerStatus) return null

    switch (containerStatus.state) {
      case "running":
        return (
          <Badge className="bg-green-500 text-white flex items-center gap-1.5">
            <span className="relative flex h-2 w-2">
              <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-300 opacity-75"></span>
              <span className="relative inline-flex rounded-full h-2 w-2 bg-green-100"></span>
            </span>
            Running
          </Badge>
        )
      case "stopped":
        return (
          <Badge variant="secondary" className="flex items-center gap-1.5">
            <div className="h-2 w-2 rounded-full bg-gray-400" />
            Stopped
          </Badge>
        )
      default:
        return (
          <Badge variant="outline" className="flex items-center gap-1.5">
            {containerStatus.state}
          </Badge>
        )
    }
  }

  const formatUptime = (uptime: string) => {
    try {
      const startTime = new Date(uptime)
      const now = new Date()
      const diffMs = now.getTime() - startTime.getTime()
      
      const seconds = Math.floor(diffMs / 1000)
      const minutes = Math.floor(seconds / 60)
      const hours = Math.floor(minutes / 60)
      const days = Math.floor(hours / 24)

      if (days > 0) return `${days}d ${hours % 24}h`
      if (hours > 0) return `${hours}h ${minutes % 60}m`
      if (minutes > 0) return `${minutes}m`
      return `${seconds}s`
    } catch {
      return "N/A"
    }
  }

  return (
    <Card className="border-border/50">
      <CardHeader className="border-b border-border/50 bg-muted/30">
        <div className="flex items-center justify-between">
          <CardTitle className="text-lg font-semibold flex items-center gap-2">
            <Activity className="h-5 w-5 text-primary" />
            Container Status
          </CardTitle>
          <Button
            variant="outline"
            size="sm"
            onClick={fetchContainerStatus}
            disabled={loading}
          >
            <RefreshCw className={`h-4 w-4 ${loading ? "animate-spin" : ""}`} />
          </Button>
        </div>
      </CardHeader>

      <CardContent className="p-6">
        {!containerStatus && loading && (
          <div className="flex items-center justify-center py-8 text-muted-foreground">
            <Loader2 className="h-6 w-6 animate-spin mr-2" />
            Loading container status...
          </div>
        )}

        {containerStatus && (
          <div className="space-y-4">
            {/* Status Row */}
            <div className="flex items-center justify-between py-3 border-b border-border/50">
              <div className="flex items-center gap-2 text-sm text-muted-foreground">
                <Server className="h-4 w-4" />
                <span>Status</span>
              </div>
              <div className="flex items-center gap-2">
                {getStateBadge()}
              </div>
            </div>

            {/* Container Name */}
            <div className="flex items-center justify-between py-3 border-b border-border/50">
              <div className="flex items-center gap-2 text-sm text-muted-foreground">
                <Server className="h-4 w-4" />
                <span>Container</span>
              </div>
              <code className="text-sm font-mono">{containerStatus.name}</code>
            </div>

            {/* Health Status */}
            {containerStatus.state === "running" && (
              <>
                <div className="flex items-center justify-between py-3 border-b border-border/50">
                  <div className="flex items-center gap-2 text-sm text-muted-foreground">
                    <Activity className="h-4 w-4" />
                    <span>Health</span>
                  </div>
                  {containerStatus.healthy ? (
                    <Badge variant="outline" className="flex items-center gap-1.5 text-green-600 border-green-600">
                      <CheckCircle2 className="h-3 w-3" />
                      Healthy
                    </Badge>
                  ) : (
                    <Badge variant="outline" className="flex items-center gap-1.5 text-yellow-600 border-yellow-600">
                      <AlertCircle className="h-3 w-3" />
                      Unknown
                    </Badge>
                  )}
                </div>

                {/* Uptime */}
                <div className="flex items-center justify-between py-3 border-b border-border/50">
                  <div className="flex items-center gap-2 text-sm text-muted-foreground">
                    <Clock className="h-4 w-4" />
                    <span>Uptime</span>
                  </div>
                  <div className="text-sm font-medium">
                    {formatUptime(containerStatus.uptime)}
                  </div>
                </div>

                {/* Started At */}
                <div className="flex items-center justify-between py-3 border-b border-border/50">
                  <div className="flex items-center gap-2 text-sm text-muted-foreground">
                    <Play className="h-4 w-4" />
                    <span>Started</span>
                  </div>
                  <div className="text-sm">
                    {new Date(containerStatus.uptime).toLocaleString()}
                  </div>
                </div>
              </>
            )}

            {/* Preview URL */}
            {previewUrl && containerStatus.state === "running" && (
              <div className="pt-2">
                <a
                  href={previewUrl}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="flex items-center justify-center gap-2 w-full py-2 px-4 rounded-lg bg-primary text-primary-foreground hover:bg-primary/90 transition-colors font-medium text-sm"
                >
                  <ExternalLink className="h-4 w-4" />
                  View Live App
                </a>
              </div>
            )}

            {/* Last Updated */}
            <div className="pt-2 text-xs text-muted-foreground text-center">
              Last updated: {lastUpdated.toLocaleTimeString()}
            </div>
          </div>
        )}

        {!containerStatus && !loading && (
          <div className="text-center py-8 text-muted-foreground">
            <AlertCircle className="h-8 w-8 mx-auto mb-2 opacity-50" />
            <p className="text-sm">No container information available</p>
          </div>
        )}
      </CardContent>
    </Card>
  )
}
