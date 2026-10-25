"use client";

import { useEffect, useState, useCallback, useRef } from "react";
import {
  Activity,
  Cpu,
  Database,
  HardDrive,
  RefreshCw,
  Pause,
  Play,
  Server,
  Zap,
  CheckCircle2,
  AlertTriangle,
} from "lucide-react";
import {
  AreaChart,
  Area,
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";

export interface SystemHealthData {
  status: string;
  uptime: string;
  runtime: {
    cpus: number;
    goroutines: number;
    alloc_mb: number;
    total_alloc_mb: number;
    sys_mb: number;
    heap_inuse_mb: number;
    num_gc: number;
  };
  db_pool?: {
    acquired_conns: number;
    idle_conns: number;
    total_conns: number;
    max_conns: number;
  };
  redis_pool?: {
    total_conns: number;
    idle_conns: number;
    stale_conns: number;
    hits: number;
    misses: number;
    timeouts: number;
  };
}

interface HistoryPoint {
  time: string;
  alloc_mb: number;
  sys_mb: number;
  goroutines: number;
  db_acquired: number;
}

const INTERVAL_OPTIONS = [
  { label: "2s (Fast)", value: 2000 },
  { label: "5s (Normal)", value: 5000 },
  { label: "10s", value: 10000 },
  { label: "Paused", value: 0 },
];

const RENDER_STARTER_MEM_LIMIT_MB = 400; // Render Starter container GOMEMLIMIT

export function SystemHealthDashboard() {
  const [data, setData] = useState<SystemHealthData | null>(null);
  const [history, setHistory] = useState<HistoryPoint[]>([]);
  const [intervalMs, setIntervalMs] = useState<number>(5000);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [lastUpdated, setLastUpdated] = useState<Date | null>(null);
  const timerRef = useRef<NodeJS.Timeout | null>(null);

  const fetchMetrics = useCallback(async () => {
    try {
      const res = await fetch("/api/admin/system-health");
      if (!res.ok) {
        throw new Error(`HTTP ${res.status}`);
      }
      const json: SystemHealthData = await res.json();
      setData(json);
      setError(null);
      setLastUpdated(new Date());

      const now = new Date();
      const timeStr = now.toLocaleTimeString([], { hour12: false, hour: "2-digit", minute: "2-digit", second: "2-digit" });

      setHistory((prev) => {
        const point: HistoryPoint = {
          time: timeStr,
          alloc_mb: parseFloat(json.runtime.alloc_mb.toFixed(2)),
          sys_mb: parseFloat(json.runtime.sys_mb.toFixed(2)),
          goroutines: json.runtime.goroutines,
          db_acquired: json.db_pool ? json.db_pool.acquired_conns : 0,
        };
        const next = [...prev, point];
        return next.slice(-25); // retain last 25 ticks
      });
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : "Failed to connect";
      setError(msg);
    } finally {
      setLoading(false);
    }
  }, []);

  // Polling loop
  useEffect(() => {
    fetchMetrics();

    if (intervalMs > 0) {
      timerRef.current = setInterval(fetchMetrics, intervalMs);
    }

    return () => {
      if (timerRef.current) clearInterval(timerRef.current);
    };
  }, [fetchMetrics, intervalMs]);

  const memPercent = data
    ? Math.min(Math.round((data.runtime.alloc_mb / RENDER_STARTER_MEM_LIMIT_MB) * 100), 100)
    : 0;

  const dbConnsAcquired = data?.db_pool?.acquired_conns ?? 0;
  const dbConnsMax = data?.db_pool?.max_conns ?? 10;
  const dbPercent = dbConnsMax > 0 ? Math.round((dbConnsAcquired / dbConnsMax) * 100) : 0;

  const redisHits = data?.redis_pool?.hits ?? 0;
  const redisMisses = data?.redis_pool?.misses ?? 0;
  const redisTotal = redisHits + redisMisses;
  const redisHitRate = redisTotal > 0 ? ((redisHits / redisTotal) * 100).toFixed(1) : "100";

  return (
    <div className="flex flex-col gap-6">
      {/* Top Controls Bar */}
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-3">
          <div className="flex items-center gap-2">
            <span className="relative flex h-3 w-3">
              {error ? (
                <span className="relative inline-flex rounded-full h-3 w-3 bg-red-500" />
              ) : (
                <>
                  <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75" />
                  <span className="relative inline-flex rounded-full h-3 w-3 bg-emerald-500" />
                </>
              )}
            </span>
            <span className="text-sm font-medium">
              {error ? (
                <span className="text-destructive">Disconnected ({error})</span>
              ) : (
                <span className="text-muted-foreground">
                  Backend Online • Uptime: <strong className="text-foreground">{data?.uptime ?? "..."}</strong>
                </span>
              )}
            </span>
          </div>
          {lastUpdated && (
            <span className="text-xs text-muted-foreground hidden md:inline">
              (Updated {lastUpdated.toLocaleTimeString()})
            </span>
          )}
        </div>

        {/* Polling controls */}
        <div className="flex items-center gap-2 flex-wrap">
          <div className="flex items-center rounded-lg border bg-card p-1 text-xs">
            {INTERVAL_OPTIONS.map((opt) => (
              <button
                key={opt.value}
                onClick={() => setIntervalMs(opt.value)}
                className={`px-2.5 py-1 rounded font-medium transition-colors ${
                  intervalMs === opt.value
                    ? "bg-primary text-primary-foreground shadow-xs"
                    : "text-muted-foreground hover:text-foreground"
                }`}
              >
                {opt.label}
              </button>
            ))}
          </div>

          <Button
            variant="outline"
            size="sm"
            onClick={() => fetchMetrics()}
            disabled={loading}
            className="gap-1.5 h-8 text-xs"
          >
            <RefreshCw className={`h-3.5 w-3.5 ${loading ? "animate-spin" : ""}`} />
            Refresh
          </Button>

          <Button
            variant="ghost"
            size="sm"
            onClick={() => setIntervalMs((prev) => (prev === 0 ? 5000 : 0))}
            className="h-8 text-xs gap-1"
          >
            {intervalMs === 0 ? (
              <>
                <Play className="h-3 w-3 text-emerald-500" /> Resume
              </>
            ) : (
              <>
                <Pause className="h-3 w-3 text-amber-500" /> Pause
              </>
            )}
          </Button>
        </div>
      </div>

      {/* KPI Cards Grid */}
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {/* Memory Card */}
        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
            <CardTitle className="text-xs font-medium text-muted-foreground uppercase tracking-wider">
              Memory (RAM)
            </CardTitle>
            <HardDrive className="h-4 w-4 text-emerald-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold tracking-tight">
              {data ? `${data.runtime.alloc_mb.toFixed(1)} MB` : "..."}
            </div>
            <div className="mt-2">
              <div className="flex justify-between text-xs text-muted-foreground mb-1">
                <span>Heap In-Use</span>
                <span>{memPercent}% of 400MB</span>
              </div>
              <div className="h-2 w-full rounded-full bg-secondary overflow-hidden">
                <div
                  className={`h-full transition-all duration-300 ${
                    memPercent > 80 ? "bg-destructive" : memPercent > 60 ? "bg-amber-500" : "bg-emerald-500"
                  }`}
                  style={{ width: `${memPercent}%` }}
                />
              </div>
            </div>
            <div className="mt-3 flex items-center justify-between text-xs text-muted-foreground border-t pt-2">
              <span>OS Sys: {data?.runtime.sys_mb.toFixed(1)} MB</span>
              <span>GC: {data?.runtime.num_gc}</span>
            </div>
          </CardContent>
        </Card>

        {/* Goroutines & CPU */}
        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
            <CardTitle className="text-xs font-medium text-muted-foreground uppercase tracking-wider">
              Goroutines & CPU
            </CardTitle>
            <Cpu className="h-4 w-4 text-blue-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold tracking-tight flex items-baseline gap-2">
              <span>{data ? data.runtime.goroutines : "..."}</span>
              <span className="text-xs font-normal text-muted-foreground">active tasks</span>
            </div>
            <div className="mt-2 flex items-center gap-2">
              <Badge
                variant={
                  (data?.runtime.goroutines ?? 0) > 1000
                    ? "destructive"
                    : (data?.runtime.goroutines ?? 0) > 300
                    ? "secondary"
                    : "outline"
                }
                className="text-[11px]"
              >
                {(data?.runtime.goroutines ?? 0) > 1000
                  ? "High Load"
                  : (data?.runtime.goroutines ?? 0) > 300
                  ? "Elevated"
                  : "Healthy"}
              </Badge>
              <span className="text-xs text-muted-foreground">
                {data ? `${data.runtime.cpus} CPU cores` : ""}
              </span>
            </div>
            <div className="mt-3 text-xs text-muted-foreground border-t pt-2">
              Lightweight Go scheduler threads
            </div>
          </CardContent>
        </Card>

        {/* PostgreSQL Connection Pool */}
        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
            <CardTitle className="text-xs font-medium text-muted-foreground uppercase tracking-wider">
              PostgreSQL Pool
            </CardTitle>
            <Database className="h-4 w-4 text-indigo-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold tracking-tight flex items-baseline gap-1.5">
              <span>{dbConnsAcquired}</span>
              <span className="text-sm font-normal text-muted-foreground">/ {dbConnsMax} max</span>
            </div>
            <div className="mt-2">
              <div className="flex justify-between text-xs text-muted-foreground mb-1">
                <span>Pool Saturation</span>
                <span>{dbPercent}%</span>
              </div>
              <div className="h-2 w-full rounded-full bg-secondary overflow-hidden">
                <div
                  className={`h-full transition-all duration-300 ${
                    dbPercent > 80 ? "bg-destructive" : dbPercent > 50 ? "bg-amber-500" : "bg-indigo-500"
                  }`}
                  style={{ width: `${dbPercent}%` }}
                />
              </div>
            </div>
            <div className="mt-3 flex items-center justify-between text-xs text-muted-foreground border-t pt-2">
              <span>Idle: {data?.db_pool?.idle_conns ?? 0}</span>
              <span>Total: {data?.db_pool?.total_conns ?? 0}</span>
            </div>
          </CardContent>
        </Card>

        {/* Redis Cache Pool */}
        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
            <CardTitle className="text-xs font-medium text-muted-foreground uppercase tracking-wider">
              Redis Cache
            </CardTitle>
            <Zap className="h-4 w-4 text-amber-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold tracking-tight flex items-baseline gap-1.5">
              <span>{redisHitRate}%</span>
              <span className="text-xs font-normal text-muted-foreground">hit rate</span>
            </div>
            <div className="mt-2 flex items-center justify-between text-xs text-muted-foreground">
              <span>Hits: {redisHits}</span>
              <span>Misses: {redisMisses}</span>
            </div>
            <div className="mt-3 flex items-center justify-between text-xs text-muted-foreground border-t pt-2">
              <span>Conns: {data?.redis_pool?.total_conns ?? 0}</span>
              <span>Timeouts: {data?.redis_pool?.timeouts ?? 0}</span>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Real-time Charts */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Memory Timeline Chart */}
        <Card>
          <CardHeader className="pb-2">
            <div className="flex items-center justify-between">
              <div>
                <CardTitle className="text-sm font-semibold">Memory Usage Over Time (MB)</CardTitle>
                <CardDescription className="text-xs">
                  Heap allocation vs. OS allocated memory
                </CardDescription>
              </div>
              <Badge variant="outline" className="text-xs font-mono">
                {history.length} ticks
              </Badge>
            </div>
          </CardHeader>
          <CardContent>
            {history.length < 2 ? (
              <div className="flex h-56 items-center justify-center text-xs text-muted-foreground">
                Collecting metrics stream...
              </div>
            ) : (
              <ResponsiveContainer width="100%" height={220}>
                <AreaChart data={history} margin={{ top: 8, right: 12, left: -16, bottom: 0 }}>
                  <defs>
                    <linearGradient id="allocGrad" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="5%" stopColor="hsl(var(--primary))" stopOpacity={0.4} />
                      <stop offset="95%" stopColor="hsl(var(--primary))" stopOpacity={0.0} />
                    </linearGradient>
                  </defs>
                  <CartesianGrid strokeDasharray="3 3" stroke="hsl(var(--border))" vertical={false} />
                  <XAxis
                    dataKey="time"
                    tick={{ fontSize: 10, fill: "hsl(var(--muted-foreground))" }}
                    axisLine={false}
                    tickLine={false}
                  />
                  <YAxis
                    tick={{ fontSize: 10, fill: "hsl(var(--muted-foreground))" }}
                    axisLine={false}
                    tickLine={false}
                    unit="MB"
                  />
                  <Tooltip
                    contentStyle={{
                      backgroundColor: "hsl(var(--card))",
                      borderColor: "hsl(var(--border))",
                      borderRadius: "8px",
                      fontSize: "12px",
                    }}
                  />
                  <Area
                    type="monotone"
                    dataKey="alloc_mb"
                    name="Heap Alloc (MB)"
                    stroke="hsl(var(--primary))"
                    fill="url(#allocGrad)"
                    strokeWidth={2}
                  />
                  <Line
                    type="monotone"
                    dataKey="sys_mb"
                    name="OS Sys (MB)"
                    stroke="hsl(var(--muted-foreground))"
                    strokeDasharray="4 4"
                    strokeWidth={1.5}
                    dot={false}
                  />
                </AreaChart>
              </ResponsiveContainer>
            )}
          </CardContent>
        </Card>

        {/* Goroutines & DB Connections Timeline */}
        <Card>
          <CardHeader className="pb-2">
            <div className="flex items-center justify-between">
              <div>
                <CardTitle className="text-sm font-semibold">Concurrency & DB Connections</CardTitle>
                <CardDescription className="text-xs">
                  Active Goroutines & PostgreSQL pool connections
                </CardDescription>
              </div>
              <Badge variant="outline" className="text-xs font-mono">
                Live
              </Badge>
            </div>
          </CardHeader>
          <CardContent>
            {history.length < 2 ? (
              <div className="flex h-56 items-center justify-center text-xs text-muted-foreground">
                Collecting metrics stream...
              </div>
            ) : (
              <ResponsiveContainer width="100%" height={220}>
                <LineChart data={history} margin={{ top: 8, right: 12, left: -16, bottom: 0 }}>
                  <CartesianGrid strokeDasharray="3 3" stroke="hsl(var(--border))" vertical={false} />
                  <XAxis
                    dataKey="time"
                    tick={{ fontSize: 10, fill: "hsl(var(--muted-foreground))" }}
                    axisLine={false}
                    tickLine={false}
                  />
                  <YAxis
                    tick={{ fontSize: 10, fill: "hsl(var(--muted-foreground))" }}
                    axisLine={false}
                    tickLine={false}
                  />
                  <Tooltip
                    contentStyle={{
                      backgroundColor: "hsl(var(--card))",
                      borderColor: "hsl(var(--border))",
                      borderRadius: "8px",
                      fontSize: "12px",
                    }}
                  />
                  <Line
                    type="monotone"
                    dataKey="goroutines"
                    name="Active Goroutines"
                    stroke="#3b82f6"
                    strokeWidth={2}
                    dot={false}
                  />
                  <Line
                    type="monotone"
                    dataKey="db_acquired"
                    name="DB Conns Acquired"
                    stroke="#8b5cf6"
                    strokeWidth={2}
                    dot={false}
                  />
                </LineChart>
              </ResponsiveContainer>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Observability Guide Card */}
      <Card className="bg-card/50">
        <CardHeader className="pb-3">
          <CardTitle className="text-sm font-semibold flex items-center gap-2">
            <Activity className="h-4 w-4 text-primary" />
            Interpreting Metrics During Load Tests
          </CardTitle>
        </CardHeader>
        <CardContent className="text-xs text-muted-foreground space-y-2 leading-relaxed">
          <p>
            • <strong>RAM Ceiling:</strong> On Render Starter, your container has 512MB RAM and a{" "}
            <code>GOMEMLIMIT=400MiB</code>. If <em>Heap Alloc</em> approaches 350-400MB, the Go garbage collector will throttle to avoid OOM termination.
          </p>
          <p>
            • <strong>DB Saturation:</strong> Your backend pool allows up to 10 connections. If <em>DB Conns Acquired</em> stays near 10 during a test, queries are queuing up.
          </p>
          <p>
            • <strong>Goroutine Stability:</strong> When a load test ends, active goroutines should drop back to baseline (~25–50). If goroutines steadily climb and do not return, it indicates connection or context leaks.
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
