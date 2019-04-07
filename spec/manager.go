package spec

import (
	"fmt"
)

// BuildUserAgent creates User-Agent, also used in agent-name of host's meta
func BuildUserAgent(version, revision string) string {
	return fmt.Sprintf("mackerel-remora/%s (Revision %s)", version, revision)
}
