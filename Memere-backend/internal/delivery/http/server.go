package http

import (
	"math"
	"net/http"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	redis "github.com/redis/go-redis/v9"

	"github.com/Bereke1t2/Memere/memere-backend/internal/version"
)

var (
	startTime      = time.Now()
	lastSampleTime = time.Now()
	lastCPUTime    = cpuDuration()
	cpuPercent     float64
	cpuMu          sync.RWMutex
)

func cpuDuration() time.Duration {
	var r syscall.Rusage
	if err := syscall.Getrusage(syscall.RUSAGE_SELF, &r); err != nil {
		return 0
	}
	return time.Duration(r.Utime.Sec+r.Stime.Sec)*time.Second +
		time.Duration(r.Utime.Usec+r.Stime.Usec)*time.Microsecond
}

func init() {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for now := range ticker.C {
			curCPUTime := cpuDuration()
			cpuMu.Lock()
			deltaWall := now.Sub(lastSampleTime).Seconds()
			deltaCPU := (curCPUTime - lastCPUTime).Seconds()
			if deltaWall > 0 {
				numCPU := float64(runtime.NumCPU())
				if numCPU < 1 {
					numCPU = 1
				}
				pct := (deltaCPU / (deltaWall * numCPU)) * 100
				if pct < 0 {
					pct = 0
				}
				if pct > 100 {
					pct = 100
				}
				cpuPercent = math.Round(pct*10) / 10
			}
			lastSampleTime = now
			lastCPUTime = curCPUTime
			cpuMu.Unlock()
		}
	}()
}

func getCPUPercent() float64 {
	cpuMu.RLock()
	defer cpuMu.RUnlock()
	return cpuPercent
}

// healthHandler is the liveness probe (/healthz). Returns 200 if the process
// is alive regardless of dependency state. Kubernetes restarts the pod only
// when this returns non-200 or times out — i.e., the process is hung/dead.
func healthHandler(_ *pgxpool.Pool, _ *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

// readyzHandler is the readiness probe (/readyz). Returns 200 only when the
// process can serve traffic: both PostgreSQL and Redis must be reachable.
// Kubernetes stops routing traffic to the pod while this returns non-200.
func readyzHandler(db *pgxpool.Pool, cache *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		errDB := db.Ping(c.Request.Context())
		errRedis := cache.Ping(c.Request.Context()).Err()

		if errDB != nil || errRedis != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unavailable",
				"db":     errDB == nil,
				"redis":  errRedis == nil,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	}
}

// systemStatsHandler reports internal runtime statistics (Go memory, goroutines,
// GC stats, DB connection pool, and Redis pool). It provides full observability
// during load tests without requiring paid APM or cloud metrics dashboards.
func systemStatsHandler(db *pgxpool.Pool, cache *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		stats := gin.H{
			"status": "ok",
			"uptime": time.Since(startTime).Round(time.Second).String(),
			"runtime": gin.H{
				"goroutines":     runtime.NumGoroutine(),
				"cpus":           runtime.NumCPU(),
				"cpu_percent":    getCPUPercent(),
				"alloc_mb":       float64(m.Alloc) / (1024 * 1024),
				"total_alloc_mb": float64(m.TotalAlloc) / (1024 * 1024),
				"sys_mb":         float64(m.Sys) / (1024 * 1024),
				"heap_inuse_mb":  float64(m.HeapInuse) / (1024 * 1024),
				"num_gc":         m.NumGC,
			},
		}

		if db != nil {
			dbStat := db.Stat()
			stats["db_pool"] = gin.H{
				"acquired_conns": dbStat.AcquiredConns(),
				"idle_conns":     dbStat.IdleConns(),
				"total_conns":    dbStat.TotalConns(),
				"max_conns":      dbStat.MaxConns(),
			}
		}

		if cache != nil {
			cacheStat := cache.PoolStats()
			stats["redis_pool"] = gin.H{
				"total_conns": cacheStat.TotalConns,
				"idle_conns":  cacheStat.IdleConns,
				"stale_conns": cacheStat.StaleConns,
				"hits":        cacheStat.Hits,
				"misses":      cacheStat.Misses,
				"timeouts":    cacheStat.Timeouts,
			}
		}

		c.JSON(http.StatusOK, stats)
	}
}

// versionHandler reports the build version, commit SHA, and build time injected
// via ldflags at image build time. Useful for verifying which release is running.
func versionHandler() gin.HandlerFunc {
	body := gin.H{
		"version":    version.Version,
		"commit":     version.Commit,
		"build_time": version.BuildTime,
	}
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, body)
	}
}
