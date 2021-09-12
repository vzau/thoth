package version

import "fmt"

var (
	Version   = "dev"
	GitCommit = "HEAD"
	BuildTime = "unknown"
)

func FriendlyVersion() string {
	return fmt.Sprintf("%s-%s (built: %s)", Version, GitCommit, BuildTime)
}
