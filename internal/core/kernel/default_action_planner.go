package kernel

import (
	"context"
	"fmt"
)

// ActionTemplate describes a rule that maps a DecisionOutcome to one or more Actions.
type ActionTemplate struct {
	Type        ActionType
	Priority    ActionPriority
	Description string
	Parameters  map[string]any
}

// SimpleActionPlanner is a deterministic ActionPlanner.
// It maps each DecisionOutcome to a fixed set of ActionTemplates.
// If no mapping exists for the outcome, it returns a single no-op Action.
//
// No LLM calls are made.
type SimpleActionPlanner struct {
	// Rules maps a DecisionOutcome to one or more ActionTemplates.
	Rules map[DecisionOutcome][]ActionTemplate
}

// NewSimpleActionPlanner creates a planner with the default outcome-to-action rules.
// Callers may supply additional or replacement rules.
func NewSimpleActionPlanner(rules map[DecisionOutcome][]ActionTemplate) *SimpleActionPlanner {
	if rules == nil {
		rules = defaultActionRules()
	}
	return &SimpleActionPlanner{Rules: rules}
}

// defaultActionRules returns the built-in, RFC-aligned outcome → action mapping.
func defaultActionRules() map[DecisionOutcome][]ActionTemplate {
	return map[DecisionOutcome][]ActionTemplate{
		DecisionOutcomeApprove: {
			{
				Type:        ActionTypeNotify,
				Priority:    ActionPriorityMedium,
				Description: "notify actor of approval",
				Parameters:  map[string]any{"channel": "internal"},
			},
		},
		DecisionOutcomeReject: {
			{
				Type:        ActionTypeNotify,
				Priority:    ActionPriorityHigh,
				Description: "notify actor of rejection",
				Parameters:  map[string]any{"channel": "internal"},
			},
		},
		DecisionOutcomeEscalate: {
			{
				Type:        ActionTypeEscalate,
				Priority:    ActionPriorityHigh,
				Description: "escalate to human reviewer",
				Parameters:  map[string]any{"channel": "internal"},
			},
		},
		DecisionOutcomeDefer: {
			{
				Type:        ActionTypeSchedule,
				Priority:    ActionPriorityLow,
				Description: "schedule for later review",
				Parameters:  map[string]any{"delay": "24h"},
			},
		},
		DecisionOutcomeNeedsRevision: {
			{
				Type:        ActionTypeNotify,
				Priority:    ActionPriorityMedium,
				Description: "request revision from actor",
				Parameters:  map[string]any{"channel": "internal"},
			},
		},
		DecisionOutcomeNoAction: {
			{
				Type:        ActionTypeNoOp,
				Priority:    ActionPriorityLow,
				Description: "no action required",
				Parameters:  nil,
			},
		},
	}
}

// Plan maps the Decision outcome to a set of planned Actions.
// All returned Actions are plans only; the caller is responsible for execution.
func (p *SimpleActionPlanner) Plan(ctx context.Context, event Event, decision Decision) ([]Action, error) {
	templates, ok := p.Rules[decision.Outcome]
	if !ok {
		// Unknown outcome → emit a single no-op action so the pipeline always
		// returns a non-empty, auditable result.
		return []Action{
			{
				ID:             fmt.Sprintf("act-%s-noop", event.ID),
				DecisionID:     decision.ID,
				EventID:        event.ID,
				Type:           ActionTypeNoOp,
				Priority:       ActionPriorityLow,
				Description:    fmt.Sprintf("no rule defined for outcome %q", decision.Outcome),
				IdempotencyKey: fmt.Sprintf("%s-noop", decision.ID),
			},
		}, nil
	}

	actions := make([]Action, 0, len(templates))
	for i, tmpl := range templates {
		params := make(map[string]any, len(tmpl.Parameters))
		for k, v := range tmpl.Parameters {
			params[k] = v
		}
		actions = append(actions, Action{
			ID:             fmt.Sprintf("act-%s-%d", event.ID, i),
			DecisionID:     decision.ID,
			EventID:        event.ID,
			Type:           tmpl.Type,
			Priority:       tmpl.Priority,
			Description:    tmpl.Description,
			Parameters:     params,
			IdempotencyKey: fmt.Sprintf("%s-%d", decision.ID, i),
		})
	}
	return actions, nil
}
