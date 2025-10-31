import { useEffect, useState, useRef } from "react";
import {
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Area,
  AreaChart,
  BarChart,
  Bar,
} from "recharts";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import Loading from "../components/Loading";

interface DiskInfo {
  name: string;
  totalSpace: number;
  availableSpace: number;
  usedSpace: number;
}

interface SystemStats {
  cpuUsage: number;
  memory: {
    total: number;
    used: number;
  };
  disk: DiskInfo[];
  loadAverage: {
    oneMinute: number;
    fiveMinutes: number;
    fifteenMinutes: number;
  };
  timestamp: number;
  uptime: number;
  cpuTemperature: number;
}

export const HomePage: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const wsRef = useRef<WebSocket | null>(null);
  const [stats, setStats] = useState<SystemStats[]>([]);

  useEffect(() => {
    if (!wsRef.current) {
      wsRef.current = new WebSocket("/api/ws/stats");

      wsRef.current.onopen = () => {
        setLoading(false);
      };

      wsRef.current.onmessage = (event) => {
        const data: SystemStats = JSON.parse(event.data);
        setStats((prev) => [...prev.slice(-30), data]);
      };

      wsRef.current.onerror = (error) => {
        console.error("WebSocket error:", error);
      };
    }

    return () => {
      if (wsRef.current) {
        wsRef.current.close();
        wsRef.current = null;
      }
    };
  }, []);

  const formatMemory = (bytes: number) => {
    const gb = bytes / (1024 * 1024 * 1024);
    return `${gb.toFixed(2)} GB`;
  };

  const formatUptime = (seconds: number) => {
    const days = Math.floor(seconds / (24 * 60 * 60));
    const hours = Math.floor((seconds % (24 * 60 * 60)) / (60 * 60));
    const minutes = Math.floor((seconds % (60 * 60)) / 60);
    const parts = [];
    if (days > 0) parts.push(`${days}d`);
    if (hours > 0) parts.push(`${hours}h`);
    if (minutes > 0) parts.push(`${minutes}m`);
    return parts.join(" ");
  };

  const customTooltip = ({ active, payload, label }: any) => {
    if (active && payload && payload.length) {
      return (
        <div className="bg-popover p-3 border border-border rounded-md">
          <p className="text-foreground">
            {new Date(label * 1000).toLocaleTimeString()}
          </p>
          {payload.map((pld: any, index: number) => (
            <p key={index} style={{ color: pld.color }}>
              {pld.dataKey === "cpuUsage"
                ? `CPU: ${pld.value.toFixed(1)}%`
                : `Memory: ${formatMemory(pld.value)} / ${formatMemory(
                    stats[stats.length - 1]?.memory.total
                  )}`}
            </p>
          ))}
        </div>
      );
    }
    return null;
  };

  if (loading || stats.length === 0) {
    return (
      <div className="flex h-full w-full items-center justify-center">
        <Loading />
      </div>
    );
  }

  const latest = stats[stats.length - 1];

  return (
    <div className="min-h-screen w-full space-y-6 bg-background">
      {/* CPU + Memory charts */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card>
          <CardHeader>
            <CardTitle>CPU Usage</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="h-[300px]">
              <ResponsiveContainer width="100%" height="100%">
<AreaChart data={stats}>
                <CartesianGrid strokeDasharray="3 3" stroke="#30363D" />
                <XAxis
                  dataKey="timestamp"
                  stroke="#C9D1D9"
                  tickFormatter={(timestamp) =>
                    new Date(timestamp * 1000).toLocaleTimeString()
                  }
                />
                <YAxis
                  stroke="#C9D1D9"
                  domain={[0, 100]}
                  tickFormatter={(value) => `${value}%`}
                />
                <Tooltip content={customTooltip} />
                <Area
                  type="monotone"
                  dataKey="cpuUsage"
                  stroke="#1F6FEB"
                  fill="#1F6FEB33"
                  strokeWidth={2}
                />
              </AreaChart>
              </ResponsiveContainer>
            </div>

            <div className="mt-4 grid grid-cols-2 gap-4">
              <Card className="p-4">
                <p className="text-sm text-muted-foreground">Current Usage</p>
                <p className="text-lg font-semibold">
                  {latest.cpuUsage.toFixed(1)}%
                </p>
              </Card>
              <Card className="p-4">
                <p className="text-sm text-muted-foreground">Average Usage</p>
                <p className="text-lg font-semibold">
                  {(
                    stats.reduce((acc, curr) => acc + curr.cpuUsage, 0) /
                    stats.length
                  ).toFixed(1)}
                  %
                </p>
              </Card>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Memory Usage</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="h-[300px]">
              <ResponsiveContainer width="100%" height="100%">
<BarChart data={stats}>
                <CartesianGrid strokeDasharray="3 3" stroke="#30363D" />
                <XAxis
                  dataKey="timestamp"
                  stroke="#C9D1D9"
                  tickFormatter={(timestamp) =>
                    new Date(timestamp * 1000).toLocaleTimeString()
                  }
                />
                <YAxis
                  stroke="#C9D1D9"
                  domain={[0, stats[0]?.memory.total || "dataMax"]}
                  tickFormatter={(value) => `${formatMemory(value)}`}
                />
                <Tooltip content={customTooltip} />
                <Bar
                  dataKey="memory.used"
                  fill="#A371F7"
                  radius={[4, 4, 0, 0]}
                />
              </BarChart>
              </ResponsiveContainer>
            </div>

            <div className="mt-4 grid grid-cols-2 gap-4">
              <Card className="p-4">
                <p className="text-sm text-muted-foreground">Used Memory</p>
                <p className="text-lg font-semibold">
                  {formatMemory(latest.memory.used)}
                </p>
              </Card>
              <Card className="p-4">
                <p className="text-sm text-muted-foreground">Total Memory</p>
                <p className="text-lg font-semibold">
                  {formatMemory(latest.memory.total)}
                </p>
              </Card>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* System Load, Uptime, Temperature */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <Card>
          <CardHeader>
            <CardTitle>System Load Average</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            {(["oneMinute", "fiveMinutes", "fifteenMinutes"] as const).map(
              (key, i) => (
                <div key={key}>
                  <div className="flex justify-between mb-2">
                    <p className="text-sm text-muted-foreground">
                      {["1 Minute", "5 Minutes", "15 Minutes"][i]}
                    </p>
                    <span
                      className={`px-2 py-1 rounded text-xs ${
                        latest.loadAverage[key] > 2
                          ? "bg-destructive/20 text-destructive"
                          : "bg-primary/20 text-primary"
                      }`}
                    >
                      {latest.loadAverage[key].toFixed(2)}
                    </span>
                  </div>
                  <div className="w-full h-2 bg-muted rounded-full">
                    <div
                      className={`h-full rounded-full transition-all ${
                        latest.loadAverage[key] > 2
                          ? "bg-destructive"
                          : "bg-primary"
                      }`}
                      style={{
                        width: `${Math.min(latest.loadAverage[key] * 50, 100)}%`,
                      }}
                    />
                  </div>
                </div>
              )
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>System Uptime</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground mb-1">
              Current Uptime
            </p>
            <p className="text-2xl font-semibold">
              {formatUptime(latest.uptime)}
            </p>
            <p className="mt-3 text-sm text-muted-foreground border-t pt-3">
              Running since{" "}
              {new Date(Date.now() - latest.uptime * 1000).toLocaleString()}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>CPU Temperature</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground mb-1">
              Current Temperature
            </p>
            <p
              className={`text-2xl font-semibold ${
                latest.cpuTemperature > 80
                  ? "text-destructive"
                  : latest.cpuTemperature > 60
                  ? "text-yellow-500"
                  : "text-primary"
              }`}
            >
              {latest.cpuTemperature !== -1
                ? `${latest.cpuTemperature.toFixed(1)}°C`
                : "N/A"}
            </p>
            <div className="mt-4 h-2 bg-muted rounded-full">
              <div
                className="h-full rounded-full transition-all"
                style={{
                  width: `${Math.min(latest.cpuTemperature, 100)}%`,
                  backgroundColor:
                    latest.cpuTemperature > 80
                      ? "hsl(var(--destructive))"
                      : latest.cpuTemperature > 60
                      ? "#F59E0B"
                      : "hsl(var(--primary))",
                }}
              />
            </div>
            <p className="mt-3 text-sm text-muted-foreground">
              Status:{" "}
              <span
                className={
                  latest.cpuTemperature > 80
                    ? "text-destructive"
                    : latest.cpuTemperature > 60
                    ? "text-yellow-500"
                    : "text-primary"
                }
              >
                {latest.cpuTemperature > 80
                  ? "Critical"
                  : latest.cpuTemperature > 60
                  ? "Warning"
                  : "Normal"}
              </span>
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Disk usage */}
      <Card>
        <CardHeader>
          <CardTitle>Disk Usage</CardTitle>
        </CardHeader>
        <CardContent className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {latest.disk.map((disk) => (
            <Card key={disk.name} className="p-4">
              <div className="flex justify-between mb-2">
                <p className="font-medium">{disk.name}</p>
                <p className="text-sm text-muted-foreground">
                  {((disk.usedSpace / disk.totalSpace) * 100).toFixed(1)}%
                </p>
              </div>
              <div className="w-full h-2 bg-muted rounded-full mb-3">
                <div
                  className={`h-full rounded-full ${
                    (disk.usedSpace / disk.totalSpace) * 100 > 90
                      ? "bg-destructive"
                      : "bg-primary"
                  }`}
                  style={{
                    width: `${(disk.usedSpace / disk.totalSpace) * 100}%`,
                  }}
                />
              </div>
              <div className="grid grid-cols-2 gap-2 text-sm">
                <div>
                  <p className="text-muted-foreground">Used</p>
                  <p>{formatMemory(disk.usedSpace)}</p>
                </div>
                <div>
                  <p className="text-muted-foreground">Available</p>
                  <p>{formatMemory(disk.availableSpace)}</p>
                </div>
                <div className="col-span-2">
                  <p className="text-muted-foreground">Total</p>
                  <p>{formatMemory(disk.totalSpace)}</p>
                </div>
              </div>
            </Card>
          ))}
        </CardContent>
      </Card>
    </div>
  );
};
