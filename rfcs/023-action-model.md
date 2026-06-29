

# RFC-023 Action Model

---
**Status:** Draft  
**Authors:** Tiroq + ChatGPT  
**Last Updated:** 2026-06-28
---

## Summary

This RFC defines the Action Model for Praxis. The Action Model formalizes how Actions are requested, planned, authorized, executed, and recorded within the system. It establishes the canonical artifacts, lifecycle, ownership, and security boundaries for all Actions, and describes how Actions interact with Decisions, Reviews, Agents, and external systems. The Action Model is foundational for auditability, automation, and reliable system orchestration.

---

## Relationships to Other RFCs

| RFC        | Title                   | Relationship                                      |
|------------|------------------------|---------------------------------------------------|
| RFC-020    | Review System          | Source of Review Packages that may lead to Actions |
| RFC-021    | Decision Model         | Decisions authorize Action Requests                |
| RFC-022    | State Machine          | Actions may transition system state                |
| RFC-030    | System Architecture    | Action Model is a core architectural component     |
| RFC-031    | Service Contracts      | Actions expose commands, events, queries           |
| RFC-033    | Storage Model          | Actions and Execution Results are persisted        |

---

## Goals

- Provide a canonical, auditable model for all Actions in Praxis.
- Ensure Actions are authorized, idempotent, and traceable.
- Decouple Action Requests, Plans, Execution, and Results.
- Support both human-initiated and automated Actions.
- Enable compensation and rollback via events, not mutation.
- Integrate with review, decision, and state machine systems.

## Non-Goals

- Define specific Action types or business logic.
- Specify UI/UX for Action management.
- Replace or duplicate the Decision or Review models.

---

## Action Philosophy

Actions are the sole means by which Praxis changes the state of the world, whether internally or via external side effects. Every Action must be:

- **Authorized**: No Action is performed without a Decision.
- **Auditable**: Every Action is recorded with immutable evidence.
- **Idempotent**: Retried Actions must not cause duplicate effects.
- **Decoupled**: Planning, requesting, and execution are distinct.
- **Compensatable**: Rollbacks are explicit and append-only.

---

## Key Artifacts and Distinctions

| Artifact           | Description                                                                                   |
|--------------------|----------------------------------------------------------------------------------------------|
| **Review**         | Assessment of a proposal or change, may recommend Actions                                    |
| **Review Package** | Set of Reviews and supporting materials for a Decision                                       |
| **Decision**       | Formal authorization for an Action Request                                                   |
| **Action Request** | Immutable, authorized intent to perform an Action                                            |
| **Action Plan**    | Versioned, structured plan for how to execute an Action Request                             |
| **Action**         | Concrete, executable unit derived from a Plan, ready for execution by an Agent or Tool       |
| **Execution**      | The actual attempt to perform an Action                                                      |
| **Execution Result** | Immutable record of the outcome of an Execution (success, failure, evidence, etc.)         |

---

## Artifact Flow Diagram

```mermaid
flowchart LR
    Review -->|recommends| ReviewPackage
    ReviewPackage -->|input to| Decision
    Decision -->|authorizes| ActionRequest
    ActionRequest -->|planned as| ActionPlan
    ActionPlan -->|yields| Action
    Action -->|executed as| Execution
    Execution -->|produces| ExecutionResult
```

---

## Canonical Definitions

| Artifact           | Definition                                                                                   |
|--------------------|---------------------------------------------------------------------------------------------|
| **Review**         | An evaluation, typically by a human or automated agent, of a proposal or change.            |
| **Review Package** | The collection of Reviews and evidence supporting a Decision.                               |
| **Decision**       | A formal record authorizing one or more Action Requests.                                    |
| **Action Request** | An immutable, uniquely identified request to perform an Action, referencing the Decision.   |
| **Action Plan**    | A versioned, structured plan detailing how the Action Request will be executed.             |
| **Action**         | The atomic, executable unit derived from an Action Plan, ready for execution.               |
| **Execution**      | The process of carrying out an Action by an Agent or Tool.                                  |
| **Execution Result** | An immutable record of the outcome, evidence, and side effects of an Execution.           |

---

## Lifecycle: Action Request

1. **Created**: Action Request is created, referencing a Decision.
2. **Planned**: An Action Plan is produced, possibly in multiple versions.
3. **Ready**: Action is derived from the latest accepted Plan.
4. **Queued**: Action is scheduled for execution.
5. **Executed**: Action is executed by an Agent or Tool.
6. **Completed**: Execution Result is persisted and linked to the Action Request.

**Immutability:** Once created, an Action Request cannot be modified.

---

## Lifecycle: Action

1. **Planned**: Action is derived from an Action Plan.
2. **Assigned**: Action is assigned to an Agent or Tool for execution.
3. **In Progress**: Execution is underway.
4. **Succeeded/Failed**: Execution Result is recorded.
5. **Compensated** (optional): If needed, a compensating Action is created and executed.

---

## Action Ownership

| Artifact       | Owner                        |
|----------------|-----------------------------|
| Action Request | The system, post-Decision    |
| Action Plan    | Planner (human or agent)     |
| Action         | Executor (Agent or Tool)     |
| Execution      | Agent or Tool                |
| Execution Result | System (immutable record)  |

---

## Idempotency Requirements

- All Actions must be idempotent: repeated execution must yield the same effect as a single execution.
- Execution Results are immutable and must record all attempts.
- Compensation is performed via new Actions, not by mutating history.

---

## Human Approval Model

- Human approval is expressed via Reviews and Decisions.
- No Action is executed without a Decision.
- Human-in-the-loop is enforced at the Decision stage, not during execution.

---

## Agent Execution Model

- Agents are responsible for picking up Actions from the Action Store.
- Agents may be human or automated.
- Agents must only execute Actions via Action Requests.
- Agent responsibility includes recording Execution Results.

---

## Tool Invocation Model

- Tools are invoked by Agents as part of Action execution.
- Tool invocations are recorded as part of the Execution Result.
- Tools must not bypass Action Requests or Plans.

---

## External Side Effects

- All external side effects (e.g., API calls, resource provisioning) must be auditable.
- Execution Results must record evidence of side effects.
- External changes must be idempotent or safely compensatable.

---

## Compensation and Rollback Strategy

- Rollbacks are performed via compensating Actions, not by mutating prior records.
- Compensation creates new Events and Execution Results.
- All compensation attempts are auditable and append-only.

---

## Action Events

| Event                  | Description                                              |
|------------------------|---------------------------------------------------------|
| ActionRequested        | Action Request created                                  |
| ActionPlanned          | Action Plan produced                                    |
| ActionQueued           | Action scheduled for execution                          |
| ActionStarted          | Execution begins                                        |
| ActionSucceeded        | Execution completed successfully                        |
| ActionFailed           | Execution failed                                        |
| ActionCompensated      | Compensating Action executed                            |
| ActionRetried          | Action retried due to failure or timeout                |

---

## Action Service Contracts

| Contract Type | Example Commands / Queries / Events                         |
|---------------|------------------------------------------------------------|
| Commands      | `RequestAction`, `PlanAction`, `ExecuteAction`, `CompensateAction` |
| Queries       | `GetActionRequest`, `ListActions`, `GetExecutionResult`    |
| Events        | See Action Events above                                    |

---

## Storage Mapping

| Artifact           | Store                |
|--------------------|---------------------|
| Action Request     | Action Store        |
| Action Plan        | Action Store        |
| Action             | Action Store        |
| Execution Result   | Execution Store     |

All stores must support immutable, append-only semantics.

---

## Security Model

- Only authorized Agents may execute Actions.
- All Actions must be linked to a valid Decision.
- All Execution Results are immutable and audit-logged.
- Action Requests and Plans are immutable after creation.

---

## Architectural Invariants

1. **Reviews never execute Actions.**
2. **Decisions authorize Actions.**
3. **Action Requests are immutable after creation.**
4. **Action Plans are versioned.**
5. **Actions are executable units.**
6. **Execution Results are immutable evidence.**
7. **Agents execute only through Action Requests.**
8. **External side effects are auditable.**
9. **Retries must be idempotent.**
10. **Compensation creates new Events rather than mutating history.**

---

## Dependencies

- RFC-020 Review System
- RFC-021 Decision Model
- RFC-022 State Machine
- RFC-030 System Architecture
- RFC-031 Service Contracts
- RFC-033 Storage Model

---

## Acceptance Criteria

- All Actions in Praxis must adhere to this model.
- No Actions are executed without a Decision.
- All Action Requests and Execution Results are immutable and auditable.
- Idempotency and compensation rules are enforced in implementation.
- Action Store and Execution Store are append-only.
- Security and ownership boundaries are respected.

---

## Decision Log

| Date       | Decision / Change                  | Author(s)           |
|------------|------------------------------------|---------------------|
| 2024-06-11 | Initial draft                      | Praxis Core Team    |
| TBD        | Review and feedback                |                     |

---