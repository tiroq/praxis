package kernel

// ActionType classifies the kind of work an Action represents.
type ActionType string

const (
	ActionTypeNotify   ActionType = "notify"
	ActionTypeCreate   ActionType = "create"
	ActionTypeUpdate   ActionType = "update"
	ActionTypeDelete   ActionType = "delete"
	ActionTypeSchedule ActionType = "schedule"
	ActionTypeEscalate ActionType = "escalate"
	ActionTypeDelegate ActionType = "delegate"
	ActionTypeArchive  ActionType = "archive"
	ActionTypeNoOp     ActionType = "noop"
)

// ActionPriority indicates the urgency of an Action.
type ActionPriority string

const (
	ActionPriorityHigh   ActionPriority = "high"
	ActionPriorityMedium ActionPriority = "medium"
	ActionPriorityLow    ActionPriority = "low"
)

// Action is the concrete, executable unit derived from a Decision.
// Per RFC-023: Actions are the sole means by which Praxis changes the state
// of the world. Every Action must be authorized, auditable, and idempotent.
// Actions returned by the kernel are planned but not yet executed.
type Action struct {
	ID         string
	DecisionID string // the decision that authorized this action
	EventID    string // the original triggering event

	Type        ActionType
	Priority    ActionPriority
	Description string
	Parameters  map[string]any // action-specific parameters; immutable once set

	// Idempotency key: callers must use this to deduplicate executions.
	IdempotencyKey string
}
