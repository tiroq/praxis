package kernel

import "time"

// ReviewRecommendation is the proposed outcome of a Review.
// Per RFC-020 §10.
type ReviewRecommendation string

const (
	ReviewRecommendationApprove       ReviewRecommendation = "approve"
	ReviewRecommendationReject        ReviewRecommendation = "reject"
	ReviewRecommendationNeedsRevision ReviewRecommendation = "needs_revision"
	ReviewRecommendationEscalate      ReviewRecommendation = "escalate"
	ReviewRecommendationInformational ReviewRecommendation = "informational"
)

// ReviewCategory classifies the purpose of a Review.
// Per RFC-020 §9.
type ReviewCategory string

const (
	ReviewCategoryTechnical    ReviewCategory = "technical"
	ReviewCategoryBusiness     ReviewCategory = "business"
	ReviewCategoryRisk         ReviewCategory = "risk"
	ReviewCategoryQuality      ReviewCategory = "quality"
	ReviewCategoryPolicy       ReviewCategory = "policy"
	ReviewCategorySecurity     ReviewCategory = "security"
	ReviewCategoryCost         ReviewCategory = "cost"
	ReviewCategoryPriority     ReviewCategory = "priority"
	ReviewCategoryCompleteness ReviewCategory = "completeness"
	ReviewCategoryConsistency  ReviewCategory = "consistency"
)

// ReviewAssessment is the factual evaluation produced by a reviewer.
// Per RFC-020 §8.
type ReviewAssessment struct {
	Score         float64  // 0.0–1.0 normalised score
	Findings      []string // concrete findings from the review
	Risks         []string // identified risks
	Opportunities []string // identified opportunities
	Confidence    float64  // 0.0–1.0 confidence in the assessment
}

// Review is an immutable, structured evaluation of an Event.
// Reviews are opinions, not facts. They inform Decisions but never enact them.
// Per RFC-020: reviews are append-only and must never be modified after creation.
type Review struct {
	ID         string
	EventID    string // reference to the reviewed Event
	Reviewer   string // identity of the reviewer (agent, engine, or person)
	ReviewType ReviewCategory
	Timestamp  time.Time

	Assessment     ReviewAssessment
	Recommendation ReviewRecommendation
	Explanation    string   // human-readable rationale
	Evidence       []string // cited sources, calculations, or rules
	Version        string   // version of the reviewed object/artifact
}
