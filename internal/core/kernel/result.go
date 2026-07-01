package kernel

// PipelineResult is the immutable output of a single Kernel.Run execution.
// It contains the complete, auditable trace: Review → Decision → Actions.
// Callers (transport adapters, services) receive this and decide how to
// persist, publish, or respond — the kernel itself does neither.
type PipelineResult struct {
	EventID  string
	Review   Review
	Decision Decision
	Actions  []Action
}
