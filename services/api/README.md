# Praxis API Service

The API exposes core Praxis resources and becomes the main boundary between interfaces, workers, and stored state.

Current scaffold:
- health endpoint
- in-memory work item endpoints
- in-memory event endpoints
- placeholder routes for decisions and agents

Important scaffold note:
- The current in-memory route storage is for local single-process development only.
- Production-ready deployments must replace it with shared persistence and authentication.

TODO:
- Replace in-memory storage with repository-backed persistence.
- Add authentication, validation, and review-cycle orchestration.
- Expose domain-specific read models for external projections.
