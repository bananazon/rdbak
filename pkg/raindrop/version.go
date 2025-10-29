package raindrop

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/djherbis/times"
)

var (
	Major  = "0"
	Minor  = "5"
	Patch  = "6"
	Suffix = "dev"
)

// Version returns a version string based on the SemVer parts defined at compile time. Dev builds will result in
// 0.0.0-dev. Prefix (v) and suffix can be optionally included, while suffix will only be included if one is defined.
func Version(prefix, suffix, versionFull bool) string {
	version := fmt.Sprintf("%s.%s.%s", Major, Minor, Patch)

	if prefix {
		version = fmt.Sprintf("v%s", version)
	}

	if suffix && Suffix != "" {
		version = fmt.Sprintf("%s-%s", version, Suffix)
	}

	if versionFull {
		version = fmt.Sprintf("%s-%s-%s", version, runtime.GOOS, runtime.GOARCH)

		var creationTime time.Time
		path, err := exec.LookPath("raindrop")
		if err == nil && path != "" {
			t, err := times.Stat(path)
			if err == nil {
				if t.HasBirthTime() {
					creationTime = t.BirthTime()
				}
			}
		}

		if !creationTime.IsZero() {
			version = fmt.Sprintf("%s - built %v", version, creationTime.Format("2006-01-02 15:04:05"))
		}
	}

	return version
}
