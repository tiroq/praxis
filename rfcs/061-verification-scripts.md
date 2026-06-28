

# RFC-061 Verification Scripts

---

## 1. Status, Authors, Last Updated

| Field         | Value                                 |
|---------------|---------------------------------------|
| Status        | DRAFT                                 |
| Authors       | Praxis Architecture Team              |
| Last Updated  | 2026-06-28                            |

---

## 2. Summary

Verification Scripts are **executable architecture assertions** which implement the Verification Invariants defined in [RFC-060 Verification Invariants](060-verification-invariants.md). Each script asserts, checks, or measures a specific property of the architecture, codebase, or system, and collectively they form the backbone of automated architecture verification. Unlike traditional unit tests, Verification Scripts encode design invariants, architectural contracts, and system-wide properties, ensuring architectural integrity as the system evolves.

---

## 3. Relationship to Previous RFCs

- **Depends on:** RFC-000 through RFC-060 (especially RFC-060 Verification Invariants)
- **Required before:** RFC-062 (Benchmarking Framework)

---

## 4. Goals

- Provide a systematic, automated, and reproducible way to verify all architecture invariants.
- Enable continuous verification of architectural properties in CI.
- Ensure every RFC invariant is mapped to one or more executable scripts.
- Standardize the organization, naming, and reporting of verification scripts.
- Support multiple verification layers (static, contract, workflow, etc.).

---

## 5. Non-Goals

- Not intended to replace unit/integration tests for implementation correctness.
- Not a framework for benchmarking (see RFC-062).
- Not a substitute for human code review or design review.
- Not a general-purpose test runner.

---

## 6. Verification Philosophy

Verification Scripts are **assertions about architecture**: they state what must always be true, what boundaries must not be crossed, and what properties must be maintained. Scripts must be deterministic, side-effect-free (unless evaluating workflows or replay), and executable in CI. Each script encodes a single invariant or a tightly related set of invariants.

---

## 7. Verification Layers

Verification is organized into the following layers (each with a dedicated script subdirectory):

| Layer         | Description                                                        |
|-------------- |--------------------------------------------------------------------|
| Static        | Checks on code structure, dependencies, and static properties      |
| Schema        | Validates data schemas, types, and interfaces                      |
| Contract      | Asserts API, module, and service contracts                         |
| Invariant     | General architectural or cross-cutting invariants                  |
| Workflow      | Validates system workflows and process sequences                   |
| Replay        | Replays historical events or traces to check system response       |
| Agent         | Verifies agent behaviors and interactions                          |
| Prompt        | Checks prompt templates, structure, and constraints                |
| Memory        | Asserts memory boundaries, access, and persistence                 |
| Space         | Validates spatial/namespace/module boundaries                      |
| Security      | Checks for security invariants and policy adherence                |
| Benchmark     | Hooks for performance/benchmark verification (see RFC-062)         |

---

## 8. Script Organization

All verification scripts reside under `/verify/`, with subdirectories per layer:

```
/verify/
  /static/
  /schema/
  /contract/
  /invariant/
  /workflow/
  /replay/
  /agent/
  /prompt/
  /memory/
  /space/
  /security/
  /benchmark/
```

Each script is a standalone executable (Python, Bash, etc.) with a manifest entry.

---

## 9. Naming Convention

- Scripts: `verify_<layer>_<invariant>[__<detail>].py`
- Example: `verify_static_no_circular_imports.py`
- All scripts must be uniquely named and easily mappable to their invariants.

---

## 10. Verification Manifest Format

Each `/verify/` directory contains a `manifest.yaml`:

```yaml
scripts:
  - name: verify_static_no_circular_imports.py
    description: "No circular imports in core modules"
    invariant: RFC-060-INV-001
    layer: static
  - name: verify_contract_api_consistency.py
    description: "API contracts match specification"
    invariant: RFC-060-INV-010
    layer: contract
...
```

---

## 11. RFC Mapping

Every invariant in RFC-060 must be mapped to at least one verification script. The mapping is maintained in `/verify/manifest.yaml` and cross-referenced in RFC-060.

---

## 12. Static Verification

Scripts in `/verify/static/` check code structure, dependencies, naming, and other static properties. Example: `verify_static_no_circular_imports.py`.

---

## 13. Schema Verification

Scripts in `/verify/schema/` validate data models, type schemas, and serialization formats. Example: `verify_schema_types_match.py`.

---

## 14. Contract Verification

Scripts in `/verify/contract/` assert API contracts, module boundaries, and service interfaces. Example: `verify_contract_api_consistency.py`.

---

## 15. Invariant Verification

Scripts in `/verify/invariant/` check cross-cutting or global architectural invariants. Example: `verify_invariant_singleton_enforced.py`.

---

## 16. Workflow Verification

Scripts in `/verify/workflow/` simulate or validate system workflows and process sequences. Example: `verify_workflow_onboarding.py`.

---

## 17. Replay Verification

Scripts in `/verify/replay/` replay captured events or request traces to verify system behavior. Example: `verify_replay_last_migration.py`.

---

## 18. Agent Verification

Scripts in `/verify/agent/` check agent behaviors, policies, and interactions. Example: `verify_agent_role_assignment.py`.

---

## 19. Prompt Verification

Scripts in `/verify/prompt/` check prompt templates, structure, and constraints. Example: `verify_prompt_template_completeness.py`.

---

## 20. Memory Verification

Scripts in `/verify/memory/` assert memory boundary, access, and persistence invariants. Example: `verify_memory_retention_policy.py`.

---

## 21. Space Boundary Verification

Scripts in `/verify/space/` check module/namespace boundaries and isolation. Example: `verify_space_no_cross_boundary_access.py`.

---

## 22. Security Verification

Scripts in `/verify/security/` check security invariants, access policies, and potential vulnerabilities. Example: `verify_security_no_hardcoded_secrets.py`.

---

## 23. Integration Verification

Scripts in `/verify/integration/` (if present) check integration points and system-wide interactions. Example: `verify_integration_external_service.py`.

---

## 24. Migration Verification

Scripts in `/verify/migration/` (if present) validate migration paths, data transformations, and upgrade invariants. Example: `verify_migration_schema_evolution.py`.

---

## 25. Benchmark Hooks

Scripts in `/verify/benchmark/` provide hooks for performance verification, delegating actual benchmarking to RFC-062.

---

## 26. Reporting Model

Verification scripts output:

- **JSON**: machine-readable detail per script.
- **Markdown summary**: human-readable summary for CI.

Example JSON output:
```json
{
  "script": "verify_static_no_circular_imports.py",
  "status": "PASS",
  "messages": [],
  "exit_code": 0
}
```

Example Markdown summary:
```
## Verification Report
| Script                                | Status | Messages                         |
|----------------------------------------|--------|----------------------------------|
| verify_static_no_circular_imports.py   | PASS   |                                  |
| verify_security_no_hardcoded_secrets.py| FAIL   | Found secret in config.py        |
```

---

## 27. CI Integration

- All scripts are invoked in CI as a dedicated verification stage.
- CI fails if any script returns a nonzero exit code.
- Reports are archived as artifacts and posted to PRs.

---

## 28. Exit Codes Table

| Code | Meaning       |
|------|--------------|
| 0    | PASS         |
| 1    | ERROR        |
| 2    | WARNING      |
| 3    | INFO         |
| 4    | SKIPPED      |

---

## 29. Failure Classification Table

| Level   | Description                                      | CI Action   |
|---------|--------------------------------------------------|-------------|
| ERROR   | Invariant violated; must be fixed                | Fail build  |
| WARNING | Non-fatal issue; should be addressed             | Warn only   |
| INFO    | Informational; no action required                | Log only    |

---

## 30. Verification Invariants

All scripts must encode one or more invariants from RFC-060, and must be traceable to their source invariant via manifest.

---

## 31. Architectural Consequences

- Enables automated enforcement of architectural decisions.
- Makes architectural drift visible and actionable.
- Provides a living, executable documentation of system invariants.
- Increases confidence in system evolution and refactoring.

---

## 32. Dependencies

- Scripting language interpreters (Python 3.10+, Bash)
- RFC-060 Verification Invariants
- CI system (GitHub Actions, etc.)

---

## 33. Acceptance Criteria

- All RFC-060 invariants are mapped to at least one verification script.
- Scripts are organized as described, with manifest and reporting.
- CI integration as specified.
- Scripts are deterministic and side-effect-free (except workflow/replay).
- Reporting model produces both JSON and Markdown.

---

## 34. Decision Log

| Date       | Entry                                                                 |
|------------|-----------------------------------------------------------------------|
| 2026-06-28 | Initial draft of Verification Scripts RFC created.                    |
| 2026-06-28 | Adopted manifest-based mapping of invariants to scripts.              |

---