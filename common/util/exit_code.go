package util

const (
	ExitError      int = 1
	ExitClean      int = 0
	ExitBadOptions int = 3
	ExitKill       int = 4
	// Go reserves exit code 2 for its own use
)
