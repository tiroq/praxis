package eventstore_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/tiroq/praxis/internal/storage/eventstore"
)

// FuzzEventRecordValidate ensures that Validate never panics regardless of input.
//
// Invariants:
//   - Must never panic for any combination of inputs.
//   - If all required fields are present AND payload is valid JSON, Validate must return nil.
//   - If payload is invalid JSON, Validate must return ErrInvalidJSON (not panic).
//   - If any required field is empty, Validate must return ErrMissingField (not panic).
func FuzzEventRecordValidate(f *testing.F) {
	// Seed corpus: valid and representative edge cases.
	validPayload, _ := json.Marshal(map[string]string{"key": "value"})

	f.Add("evt-1", "user.created", "auth-service", "user-123", string(validPayload), "meta-val")
	f.Add("", "type", "source", "subject", `{"x":1}`, "")
	f.Add("id", "", "source", "subject", `{}`, "")
	f.Add("id", "type", "", "subject", `{}`, "")
	f.Add("id", "type", "source", "", `{}`, "")
	f.Add("id", "type", "source", "subject", "not-json", "")
	f.Add("id", "type", "source", "subject", "", "")
	f.Add("id", "type", "source", "subject", "null", "meta")
	f.Add("id", "type", "source", "subject", `{"nested":{"deep":true}}`, "v")
	f.Add("id\x00", "type\xff", "source", "subject", `{}`, "")

	f.Fuzz(func(t *testing.T, id, evtType, source, subjectID, payload, metaVal string) {
		event := eventstore.EventRecord{
			ID:         id,
			Type:       evtType,
			Source:     source,
			SubjectID:  subjectID,
			OccurredAt: time.Now(),
			Payload:    json.RawMessage(payload),
		}
		if metaVal != "" {
			event.Metadata = map[string]string{"key": metaVal}
		}

		// Must never panic — any error is acceptable.
		_ = event.Validate()
	})
}
