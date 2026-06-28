You are acting as Principal Software Architect.

Project: Praxis

Your task is NOT to rewrite documents.

Your task is to perform a complete architecture review of every RFC in this repository.

You should behave like an architecture review board preparing a system for long-term development.

The architecture should be evaluated against the following qualities:

• consistency
• completeness
• separation of concerns
• terminology consistency
• bounded contexts
• object ownership
• identity model
• lifecycle model
• event model
• decision model
• review model
• scalability
• extensibility
• future maintainability
• CQRS compatibility
• Event-driven compatibility
• multi-agent compatibility
• provider independence
• observability
• replayability
• auditability

Review ALL RFCs.

For each RFC:

1. Summarize its responsibility in one sentence.

2. Verify that it does not overlap responsibilities with other RFCs.

3. Check terminology consistency.

4. Check whether concepts introduced here already exist elsewhere.

5. Detect circular dependencies.

6. Detect missing abstractions.

7. Detect duplicated abstractions.

8. Detect missing invariants.

9. Detect missing diagrams.

10. Detect contradictions.

11. Detect responsibilities that belong in another RFC.

12. Check naming quality.

13. Check if future implementation would become harder because of this RFC.

14. Check if concepts are over-engineered.

15. Check if concepts are under-specified.

16. Verify dependency ordering.

17. Verify roadmap ordering.

18. Verify architectural layering.

19. Verify that every RFC has:

- Goals
- Non-Goals
- Invariants
- Dependencies
- Acceptance Criteria
- Decision Log

20. Suggest improvements.

Important rules:

Do NOT rewrite documents.

Do NOT make stylistic suggestions.

Only report architectural issues.

For every issue provide:

Severity:
Critical
Major
Minor

Reason

Recommendation

Affected RFC(s)

Finally produce:

1. Foundation Architecture Score (0-100)

2. Intelligence Layer Score (0-100)

3. Runtime Architecture Score (0-100)

4. Overall Architecture Score (0-100)

5. Top 20 improvements sorted by impact.