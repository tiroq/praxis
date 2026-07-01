package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	natstransport "github.com/tiroq/praxis/internal/transport/nats"
)

type outputResult struct {
	Decision struct {
		ID      string `json:"id"`
		Outcome string `json:"outcome"`
	} `json:"decision"`
	Actions []json.RawMessage `json:"actions"`
}

const smokeTimeout = 10 * time.Second

type smokeReport struct {
	NATSURL       string                      `json:"nats_url"`
	Stream        string                      `json:"stream"`
	InputSubject  string                      `json:"input_subject"`
	OutputSubject string                      `json:"output_subject"`
	Published     natstransport.InputMessage  `json:"published"`
	Received      natstransport.OutputMessage `json:"received"`
	VerifiedAt    time.Time                   `json:"verified_at"`
	WorkerFlowOK  bool                        `json:"worker_flow_ok"`
	Message       string                      `json:"message"`
}

func main() {
	var outPath string
	flag.StringVar(&outPath, "out", "", "optional path to also write the smoke result JSON")
	flag.Parse()

	cfg := natstransport.ConfigFromEnv()
	report, err := runSmoke(cfg, smokeTimeout)
	if err != nil {
		emitFailure(cfg, outPath, err)
		os.Exit(1)
	}

	if err := emitReport(report, outPath); err != nil {
		fmt.Fprintf(os.Stderr, "failed to emit smoke report: %v\n", err)
		os.Exit(1)
	}
}

func runSmoke(cfg natstransport.Config, timeout time.Duration) (smokeReport, error) {
	client, err := natstransport.NewClient(cfg)
	if err != nil {
		return smokeReport{}, err
	}
	defer client.Close()

	js := client.JetStream()

	sub, err := js.SubscribeSync(
		cfg.OutputSubject,
		nats.DeliverNew(),
		nats.AckExplicit(),
		nats.BindStream(cfg.StreamName),
	)
	if err != nil {
		return smokeReport{}, fmt.Errorf("subscribe to output subject %q: %w", cfg.OutputSubject, err)
	}
	defer sub.Unsubscribe() //nolint:errcheck

	input := natstransport.InputMessage{
		ID:        fmt.Sprintf("evt_smoke_%d", time.Now().UTC().UnixNano()),
		Source:    "nats-smoke",
		Text:      "urgent review: buy tickets to Shanghai",
		Timestamp: time.Now().UTC(),
		Metadata: map[string]string{
			"smoke": "true",
		},
	}

	payload, err := json.Marshal(input)
	if err != nil {
		return smokeReport{}, fmt.Errorf("marshal input message: %w", err)
	}

	if _, err := js.Publish(cfg.InputSubject, payload, nats.MsgId(input.ID)); err != nil {
		return smokeReport{}, fmt.Errorf("publish input message: %w", err)
	}

	msg, err := sub.NextMsg(timeout)
	if err != nil {
		if errors.Is(err, nats.ErrTimeout) {
			return smokeReport{}, fmt.Errorf("timeout waiting for output on %q after %s", cfg.OutputSubject, timeout)
		}
		return smokeReport{}, fmt.Errorf("receive output message: %w", err)
	}
	defer msg.Ack() //nolint:errcheck

	var output natstransport.OutputMessage
	if err := json.Unmarshal(msg.Data, &output); err != nil {
		return smokeReport{}, fmt.Errorf("decode output message: %w", err)
	}

	if err := validateOutput(input.ID, output); err != nil {
		return smokeReport{}, err
	}

	return smokeReport{
		NATSURL:       cfg.URL,
		Stream:        cfg.StreamName,
		InputSubject:  cfg.InputSubject,
		OutputSubject: cfg.OutputSubject,
		Published:     input,
		Received:      output,
		VerifiedAt:    time.Now().UTC(),
		WorkerFlowOK:  true,
		Message:       "validated worker -> kernel -> output flow over real NATS JetStream",
	}, nil
}

func validateOutput(expectedEventID string, output natstransport.OutputMessage) error {
	if output.Status != "ok" {
		return fmt.Errorf("unexpected status %q", output.Status)
	}
	if output.InputEventID != expectedEventID {
		return fmt.Errorf("unexpected input_event_id %q", output.InputEventID)
	}
	if output.Result == nil {
		return errors.New("missing result")
	}
	var result outputResult
	if err := json.Unmarshal(output.Result, &result); err != nil {
		return fmt.Errorf("decode result: %w", err)
	}
	if result.Decision.ID == "" || result.Decision.Outcome == "" {
		return errors.New("result missing decision")
	}
	if len(result.Actions) == 0 {
		return errors.New("result missing actions")
	}
	return nil
}

func emitFailure(cfg natstransport.Config, outPath string, err error) {
	report := map[string]any{
		"nats_url":       cfg.URL,
		"stream":         cfg.StreamName,
		"input_subject":  cfg.InputSubject,
		"output_subject": cfg.OutputSubject,
		"worker_flow_ok": false,
		"error":          err.Error(),
		"verified_at":    time.Now().UTC(),
	}
	if emitErr := emitJSON(report, outPath); emitErr != nil {
		fmt.Fprintf(os.Stderr, "failed to emit smoke failure JSON: %v\n", emitErr)
		fmt.Fprintf(os.Stderr, "smoke error: %v\n", err)
		return
	}
	fmt.Fprintln(os.Stderr, err.Error())
}

func emitReport(report smokeReport, outPath string) error {
	return emitJSON(report, outPath)
}

func emitJSON(v any, outPath string) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if _, err := os.Stdout.Write(data); err != nil {
		return err
	}
	if outPath == "" {
		return nil
	}
	return os.WriteFile(outPath, data, 0o644)
}
