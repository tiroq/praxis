package main

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/tiroq/praxis/internal/core/kernel"
	natstransport "github.com/tiroq/praxis/internal/transport/nats"
)

func mustResult(t *testing.T, result kernel.PipelineResult) json.RawMessage {
	t.Helper()
	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("marshal result: %v", err)
	}
	return data
}

func TestValidateOutput_SucceedsOnExpectedResult(t *testing.T) {
	out := natstransport.OutputMessage{
		InputEventID: "evt_123",
		Status:       "ok",
		Result: mustResult(t, kernel.PipelineResult{
			Decision: kernel.Decision{ID: "dec_123", Outcome: kernel.DecisionOutcomeApprove},
			Actions:  []kernel.Action{{ID: "act_123"}},
		}),
	}

	if err := validateOutput("evt_123", out); err != nil {
		t.Fatalf("validateOutput returned error: %v", err)
	}
}

func TestValidateOutput_FailsOnInvalidResult(t *testing.T) {
	tests := []struct {
		name    string
		output  natstransport.OutputMessage
		wantErr string
	}{
		{
			name:    "status",
			output:  natstransport.OutputMessage{InputEventID: "evt_123", Status: "error"},
			wantErr: "unexpected status",
		},
		{
			name:    "event id",
			output:  natstransport.OutputMessage{InputEventID: "evt_other", Status: "ok"},
			wantErr: "unexpected input_event_id",
		},
		{
			name:    "missing result",
			output:  natstransport.OutputMessage{InputEventID: "evt_123", Status: "ok"},
			wantErr: "missing result",
		},
		{
			name: "missing decision",
			output: natstransport.OutputMessage{
				InputEventID: "evt_123",
				Status:       "ok",
				Result:       mustResult(t, kernel.PipelineResult{}),
			},
			wantErr: "result missing decision",
		},
		{
			name: "missing actions",
			output: natstransport.OutputMessage{
				InputEventID: "evt_123",
				Status:       "ok",
				Result: mustResult(t, kernel.PipelineResult{
					Decision: kernel.Decision{ID: "dec_123", Outcome: kernel.DecisionOutcomeApprove},
				}),
			},
			wantErr: "result missing actions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateOutput("evt_123", tt.output)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("expected error containing %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}
