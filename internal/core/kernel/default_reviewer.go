package kernel

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// KeywordReviewer is a deterministic, rule-based Reviewer.
// It scans Event.Text for configurable keywords and produces a structured
// Review with findings, risks, and a recommendation — no LLM calls.
//
// If no keywords match, it produces an informational review with a neutral
// score. If keywords match, the score, recommendation, and findings reflect
// the number and nature of matched keywords.
type KeywordReviewer struct {
	// Keywords maps a keyword (lower-cased) to its risk label.
	// Example: {"urgent": "time-sensitive", "blocked": "blocker"}
	Keywords map[string]string

	// ReviewerName is the identity string recorded in the Review.
	ReviewerName string
}

// NewKeywordReviewer creates a KeywordReviewer with the supplied keyword map.
// ReviewerName defaults to "keyword-reviewer" if empty.
func NewKeywordReviewer(keywords map[string]string, reviewerName string) *KeywordReviewer {
	if reviewerName == "" {
		reviewerName = "keyword-reviewer"
	}
	kw := make(map[string]string, len(keywords))
	for k, v := range keywords {
		kw[strings.ToLower(k)] = v
	}
	return &KeywordReviewer{Keywords: kw, ReviewerName: reviewerName}
}

// Review scans Event.Text for known keywords and builds a Review accordingly.
func (r *KeywordReviewer) Review(ctx context.Context, event Event) (Review, error) {
	lower := strings.ToLower(event.Text)

	var findings []string
	var risks []string
	matched := 0

	for keyword, label := range r.Keywords {
		if strings.Contains(lower, keyword) {
			matched++
			findings = append(findings, fmt.Sprintf("keyword %q matched (label: %s)", keyword, label))
			risks = append(risks, label)
		}
	}

	recommendation := ReviewRecommendationInformational
	score := 0.5
	explanation := "no significant keywords detected; informational review"

	if matched > 0 {
		// Scale score by fraction of keywords matched, capped at 0.9.
		score = min(0.9, float64(matched)/float64(len(r.Keywords)))
		recommendation = ReviewRecommendationApprove
		explanation = fmt.Sprintf("%d keyword(s) matched; review recommends approval with noted risks", matched)
	}

	return Review{
		ID:         fmt.Sprintf("rev-%s", event.ID),
		EventID:    event.ID,
		Reviewer:   r.ReviewerName,
		ReviewType: ReviewCategoryQuality,
		Timestamp:  time.Now().UTC(),
		Assessment: ReviewAssessment{
			Score:      score,
			Findings:   findings,
			Risks:      risks,
			Confidence: event.Confidence,
		},
		Recommendation: recommendation,
		Explanation:    explanation,
		Version:        "1",
	}, nil
}

// min returns the smaller of two float64 values.
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
