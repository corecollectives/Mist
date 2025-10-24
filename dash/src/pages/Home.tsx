
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
        // toast.error("WebSocket connection error. Please try refreshing the page.");
        // setLoading(false);
      };

    }

    return () => {
      if (wsRef.current) {
        console.log("Closing WebSocket connection");
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
        <div className="bg-[#161B22] p-3 border border-[#30363D] rounded-md">
          <p className="text-[#C9D1D9]">
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
    return <div className="flex h-screen w-full items-center justify-center">
      <Loading />
    </div>;
  }

  return (
    <div className="min-h-screen bg-[#0D1117] p-6 space-y-6 w-full">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div className="bg-[#161B22] p-6 rounded-lg border border-[#30363D]">
          <h2 className="text-[#C9D1D9] text-xl font-semibold mb-4">
            CPU Usage
          </h2>
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
          {stats.length > 0 && (
            <div className="mt-4 grid grid-cols-2 gap-4">
              <div className="bg-[#0D1117] p-4 rounded-lg border border-[#30363D]">
                <p className="text-[#8B949E] text-sm">Current Usage</p>
                <p className="text-[#C9D1D9] text-lg font-semibold">
                  {stats[stats.length - 1].cpuUsage.toFixed(1)}%
                </p>
              </div>
              <div className="bg-[#0D1117] p-4 rounded-lg border border-[#30363D]">
                <p className="text-[#8B949E] text-sm">Average Usage</p>
                <p className="text-[#C9D1D9] text-lg font-semibold">
                  {(
                    stats.reduce((acc, curr) => acc + curr.cpuUsage, 0) /
                    stats.length
                  ).toFixed(1)}
                  %
                </p>
              </div>
            </div>
          )}
        </div>

        <div className="bg-[#161B22] p-6 rounded-lg border border-[#30363D]">
          <h2 className="text-[#C9D1D9] text-xl font-semibold mb-4">
            Memory Usage
          </h2>
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
          {stats.length > 0 && (
            <div className="mt-4 grid grid-cols-2 gap-4">
              <div className="bg-[#0D1117] p-4 rounded-lg border border-[#30363D]">
                <p className="text-[#8B949E] text-sm">Used Memory</p>
                <p className="text-[#C9D1D9] text-lg font-semibold">
                  {formatMemory(stats[stats.length - 1].memory.used)}
                </p>
              </div>
              <div className="bg-[#0D1117] p-4 rounded-lg border border-[#30363D]">
                <p className="text-[#8B949E] text-sm">Total Memory</p>
                <p className="text-[#C9D1D9] text-lg font-semibold">
                  {formatMemory(stats[stats.length - 1].memory.total)}
                </p>
              </div>
            </div>
          )}
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="bg-[#161B22] p-6 rounded-lg border border-[#30363D]">
          <h2 className="text-[#C9D1D9] text-xl font-semibold mb-4">
            System Load Average
          </h2>
          <div className="grid grid-cols-1 gap-4">
            {stats.length > 0 && (
              <>
                <div className="bg-[#0D1117] p-4 rounded-lg border border-[#30363D]">
                  <div className="flex items-center justify-between mb-2">
                    <p className="text-[#8B949E]">1 Minute</p>
                    <span
                      className={`px-2 py-1 rounded text-xs ${stats[stats.length - 1].loadAverage.oneMinute > 2
                        ? "bg-[#F8514933] text-[#F85149]"
                        : "bg-[#1F6FEB33] text-[#1F6FEB]"
                        }`}
                    >
                      {stats[stats.length - 1].loadAverage.oneMinute.toFixed(2)}
                    </span>
                  </div>
                  <div className="w-full h-2 bg-[#303D3D] rounded-full">
                    <div
                      className="h-full rounded-full transition-all"
                      style={{
                        width: `${Math.min(
                          stats[stats.length - 1].loadAverage.oneMinute * 50,
                          100
                        )}%`,
                        backgroundColor:
                          stats[stats.length - 1].loadAverage.oneMinute > 2
                            ? "#F85149"
                            : "#1F6FEB",
                      }}
                    />
                  </div>
                </div>

                <div className="bg-[#0D1117] p-4 rounded-lg border border-[#30363D]">
                  <div className="flex items-center justify-between mb-2">
                    <p className="text-[#8B949E]">5 Minutes</p>
                    <span
                      className={`px-2 py-1 rounded text-xs ${stats[stats.length - 1].loadAverage.fiveMinutes > 2
                        ? "bg-[#F8514933] text-[#F85149]"
                        : "bg-[#1F6FEB33] text-[#1F6FEB]"
                        }`}
                    >
                      {stats[stats.length - 1].loadAverage.fiveMinutes.toFixed(2)}
                    </span>
                  </div>
                  <div className="w-full h-2 bg-[#303D3D] rounded-full">
                    <div
                      className="h-full rounded-full transition-all"
                      style={{
                        width: `${Math.min(
                          stats[stats.length - 1].loadAverage.fiveMinutes * 50,
                          100
                        )}%`,
                        backgroundColor:
                          stats[stats.length - 1].loadAverage.fiveMinutes > 2
                            ? "#F85149"
                            : "#1F6FEB",
                      }}
                    />
                  </div>
                </div>

                <div className="bg-[#0D1117] p-4 rounded-lg border border-[#30363D]">
                  <div className="flex items-center justify-between mb-2">
                    <p className="text-[#8B949E]">15 Minutes</p>
                    <span
                      className={`px-2 py-1 rounded text-xs ${stats[stats.length - 1].loadAverage.fifteenMinutes > 2
                        ? "bg-[#F8514933] text-[#F85149]"
                        : "bg-[#1F6FEB33] text-[#1F6FEB]"
                        }`}
                    >
                      {stats[stats.length - 1].loadAverage.fifteenMinutes.toFixed(
                        2
                      )}
                    </span>
                  </div>
                  <div className="w-full h-2 bg-[#303D3D] rounded-full">
                    <div
                      className="h-full rounded-full transition-all"
                      style={{
                        width: `${Math.min(
                          stats[stats.length - 1].loadAverage.fifteenMinutes * 50,
                          100
                        )}%`,
                        backgroundColor:
                          stats[stats.length - 1].loadAverage.fifteenMinutes > 2
                            ? "#F85149"
                            : "#1F6FEB",
                      }}
                    />
                  </div>
                </div>
              </>
            )}
          </div>
        </div>

        <div className="bg-[#161B22] p-6 rounded-lg border border-[#30363D]">
          <h2 className="text-[#C9D1D9] text-xl font-semibold mb-4">
            System Uptime
          </h2>
          {stats.length > 0 && (
            <div className="bg-[#0D1117] p-4 rounded-lg border border-[#30363D]">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-[#8B949E] text-sm">Current Uptime</p>
                  <p className="text-[#C9D1D9] text-2xl font-semibold mt-1">
                    {formatUptime(stats[stats.length - 1].uptime)}
                  </p>
                </div>
                <div className="h-12 w-12 rounded-full bg-[#1F6FEB33] flex items-center justify-center">
                  <svg
                    className="w-6 h-6 text-[#1F6FEB]"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                  </svg>
                </div>
              </div>
              <div className="mt-4 pt-4 border-t border-[#303D3D]">
                <p className="text-[#8B949E] text-sm">
                  System has been running since{" "}
                  {new Date(
                    Date.now() - stats[stats.length - 1].uptime * 1000
                  ).toLocaleString("en-US", {
                    dateStyle: "medium",
                    timeStyle: "short",
                  })}
                </p>
              </div>
            </div>
          )}
        </div>

        <div className="bg-[#161B22] p-6 rounded-lg border border-[#30363D]">
          <h2 className="text-[#C9D1D9] text-xl font-semibold mb-4">
            CPU Temperature
          </h2>
          {stats.length > 0 && (
            <div className="bg-[#0D1117] p-4 rounded-lg border border-[#30363D]">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-[#8B949E] text-sm">Current Temperature</p>
                  <p
                    className={`text-2xl font-semibold mt-1 ${stats[stats.length - 1].cpuTemperature > 80
                      ? "text-[#F85149]"
                      : stats[stats.length - 1].cpuTemperature > 60
                        ? "text-[#F0883E]"
                        : "text-[#1F6FEB]"
                      }`}
                  >
                    {stats[stats.length - 1].cpuTemperature !== -1 ? stats[stats.length - 1].cpuTemperature.toFixed(1) + "°C" : "N/A"}
                  </p>
                </div>
                <div
                  className={`h-12 w-12 rounded-full flex items-center justify-center ${stats[stats.length - 1].cpuTemperature > 80
                    ? "bg-[#F8514933]"
                    : stats[stats.length - 1].cpuTemperature > 60
                      ? "bg-[#F0883E33]"
                      : "bg-[#1F6FEB33]"
                    }`}
                >
                  <svg
                    className={`w-6 h-6 ${stats[stats.length - 1].cpuTemperature > 80
                      ? "text-[#F85149]"
                      : stats[stats.length - 1].cpuTemperature > 60
                        ? "text-[#F0883E]"
                        : "text-[#1F6FEB]"
                      }`}
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"
                    />
                  </svg>
                </div>
              </div>
              <div className="mt-4">
                <div className="w-full h-2 bg-[#303D3D] rounded-full">
                  <div
                    className="h-full rounded-full transition-all"
                    style={{
                      width: `${Math.min(stats[stats.length - 1].cpuTemperature, 100)}%`,
                      backgroundColor:
                        stats[stats.length - 1].cpuTemperature > 80
                          ? "#F85149"
                          : stats[stats.length - 1].cpuTemperature > 60
                            ? "#F0883E"
                            : "#1F6FEB",
                    }}
                  />
                </div>
              </div>
              <div className="mt-4 pt-4 border-t border-[#303D3D]">
                <p className="text-[#8B949E] text-sm">
                  Status:{" "}
                  <span
                    className={
                      stats[stats.length - 1].cpuTemperature > 80
                        ? "text-[#F85149]"
                        : stats[stats.length - 1].cpuTemperature > 60
                          ? "text-[#F0883E]"
                          : "text-[#1F6FEB]"
                    }
                  >
                    {stats[stats.length - 1].cpuTemperature > 80
                      ? "Critical"
                      : stats[stats.length - 1].cpuTemperature > 60
                        ? "Warning"
                        : "Normal"}
                  </span>
                </p>
              </div>
            </div>
          )}
        </div>
      </div>

      <div className="bg-[#161B22] p-6 rounded-lg border border-[#30363D]">
        <h2 className="text-[#C9D1D9] text-xl font-semibold mb-4">
          Disk Usage
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {stats.length > 0 &&
            stats[stats.length - 1].disk.map((disk) => (
              <div
                key={disk.name}
                className="bg-[#0D1117] p-4 rounded-lg border border-[#30363D]"
              >
                <div className="flex justify-between items-start mb-2">
                  <p className="text-[#C9D1D9] font-medium">{disk.name}</p>
                  <p className="text-[#8B949E] text-sm">
                    {((disk.usedSpace / disk.totalSpace) * 100).toFixed(1)}%
                  </p>
                </div>

                <div className="w-full h-2 bg-[#303D3D] rounded-full mb-3">
                  <div
                    className="h-full bg-[#1F6FEB] rounded-full"
                    style={{
                      width: `${(disk.usedSpace / disk.totalSpace) * 100}%`,
                      backgroundColor:
                        (disk.usedSpace / disk.totalSpace) * 100 > 90
                          ? "#F85149"
                          : "#1F6FEB",
                    }}
                  />
                </div>

                <div className="grid grid-cols-2 gap-2 text-sm">
                  <div>
                    <p className="text-[#8B949E]">Used</p>
                    <p className="text-[#C9D1D9]">
                      {formatMemory(disk.usedSpace)}
                    </p>
                  </div>
                  <div>
                    <p className="text-[#8B949E]">Available</p>
                    <p className="text-[#C9D1D9]">
                      {formatMemory(disk.availableSpace)}
                    </p>
                  </div>
                  <div className="col-span-2">
                    <p className="text-[#8B949E]">Total</p>
                    <p className="text-[#C9D1D9]">
                      {formatMemory(disk.totalSpace)}
                    </p>
                  </div>
                </div>
              </div>
            ))}
        </div>
      </div>

    </div>
  );
};

