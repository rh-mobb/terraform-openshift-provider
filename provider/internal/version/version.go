package version

import (
	"fmt"
	"runtime"
)

// These variables are set via ldflags during build
var (
	// Version is the provider version (set via ldflags)
	Version = "dev"

	// GitCommit is the git commit hash (set via ldflags)
	GitCommit = "unknown"

	// BuildDate is the build date (set via ldflags)
	BuildDate = "unknown"
)

// String returns a human-readable version string
func String() string {
	return fmt.Sprintf("%s (commit: %s, built: %s, go: %s)",
		Version, GitCommit, BuildDate, runtime.Version())
}

// UserAgent returns a user agent string for API requests
func UserAgent() string {
	return fmt.Sprintf("terraform-provider-openshift/%s", Version)
}

