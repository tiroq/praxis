# Praxis API Service

The API exposes core Praxis resources and becomes the main boundary between interfaces, workers, and stored state.

Current scaffold:
- health endpoint
- in-memory work item endpoints
- in-memory event endpoints
- placeholder routes for decisions and agents

TODO:
- Replace in-memory storage with repository-backed persistence.
- Add authentication, validation, and review-cycle orchestration.
- Expose domain-specific read models for external projections.
