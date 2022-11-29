package walker

import "time"

// ResultFunc is the callback used for gathering the
// result during graph execution.
type ResultFunc func(*ResultEntry)

type ResultEntry struct {
	vertexName string
	duration   time.Duration
}
