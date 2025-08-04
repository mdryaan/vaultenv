package version

import (
	"fmt"
	"runtime"
)

var (
	Version   = "0.1.0"
	GitCommit = "none"
	BuildDate = "unknown"
	GoVersion = runtime.Version()
)

func String() string {
	return fmt.Sprintf("vaultenv %s (commit: %s, built: %s, %s/%s, %s)",
		Version, GitCommit, BuildDate, runtime.GOOS, runtime.GOARCH, GoVersion)
}

func Short() string {
	return Version
}
