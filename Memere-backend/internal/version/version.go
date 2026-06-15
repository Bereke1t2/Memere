// Package version holds build-time variables injected via -ldflags. They are
// set to meaningful values by the Dockerfile and CI pipeline; in local dev
// builds they fall back to the sentinel values below.
package version

// Variables are overridden at link time:
//
//	-X 'github.com/Bereke1t2/Memere/memere-backend/internal/version.Version=v1.2.3'
//	-X 'github.com/Bereke1t2/Memere/memere-backend/internal/version.Commit=abc1234'
//	-X 'github.com/Bereke1t2/Memere/memere-backend/internal/version.BuildTime=2026-06-15T12:00:00Z'
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
)
