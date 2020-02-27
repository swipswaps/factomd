package worker

// a single execution - return false to signal parent loop to break
type step func() bool

// syntactic sugar to wrap the execution of sequential steps
// returns False after first failure or True if all complete properly
func RunSteps(steps ...step) bool {
	for i := range steps {
		if !steps[i]() {
			return false
		}
	}
	return true
}
