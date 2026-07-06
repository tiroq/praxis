---
name: "Pure Mapper Discipline"
description: "Mandatory rules for all transport mapping functions (adapters). Enforces purity, pure function shape, referential transparency, and structural triviality."
applyTo: ["services/**", "packages/**", "apps/**", "scripts/**", "infra/**"]
---

# Pure Mapper Discipline

**Mandatory rule for all transport mapping functions (adapters).**

Every mapping function that translates an external transport object to a Praxis wire contract **MUST** be a pure function with no side effects, no external dependencies, and no business logic.

---

## Mapper Shape (Required)

```
Input:  external transport object (e.g., Telegram Update, HTTP request, AMQP message)
Output: Praxis wire contract dict (matching internal/transport/nats.InputMessage)
```

---

## Mapper MUST NOT

- ❌ Publish messages
- ❌ Acknowledge messages
- ❌ Call HTTP endpoints
- ❌ Call NATS
- ❌ Call storage (databases, caches, files)
- ❌ Access environment variables
- ❌ Perform retries
- ❌ Mutate global state
- ❌ Generate business identifiers
- ❌ Execute business logic
- ❌ Classify intent
- ❌ Infer semantics
- ❌ Enrich domain objects
- ❌ Create workflows

---

## Forbidden Complexity in Mappers

A mapper MUST NOT:

- Call repositories
- Call databases
- Call HTTP services
- Call LLMs
- Call event stores
- Call action planners
- Call reviewers
- Call decision makers
- Perform retries
- Perform caching
- Maintain state
- Generate business identifiers
- Compute business outcomes
- Infer intent
- Classify content
- Enrich with external knowledge
- Coordinate workflows

If any of these become necessary, stop and move the logic into the correct owner.

---

## Allowed Complexity

A mapper MAY ONLY:

- ✅ Rename fields
- ✅ Copy fields
- ✅ Drop unused fields
- ✅ Normalize formatting
- ✅ Convert primitive types
- ✅ Construct transport DTOs
- ✅ Perform deterministic serialization
- ✅ Validate transport-level preconditions
- ✅ Create deterministic identifiers derived solely from the input

Nothing else.

---

## Referential Transparency Invariant

A mapper must be **referentially transparent**:

> Given the same input object, it MUST always produce the same output object.
>
> No observable side effects are permitted.

---

## Verification Checklist

For every adapter with a mapping function, answer all:

| Question | Answer |
|---|---|
| Does the mapper accept only the external transport object? | _(required: must be YES)_ |
| Does the mapper return only a dict (wire contract)? | _(required: must be YES)_ |
| Does the mapper call any I/O function (publish, HTTP, storage)? | _(required: must be NO)_ |
| Does the mapper access any global state? | _(required: must be NO)_ |
| Does the mapper perform any retry logic? | _(required: must be NO)_ |
| Can this mapper be called 1,000,000 times with the same input and produce identical output? | _(required: must be YES)_ |
| Is the mapper completely free of business logic? | _(required: must be YES)_ |

**Stop condition:** If any answer fails purity verification, the mapper **MUST** be refactored before implementation proceeds.

---

## Structural Triviality Test

Review every mapper using these questions:

1. Can every output field be traced directly to one or more input fields?
2. Is every transformation deterministic?
3. Would another engineer understand the mapper in under one minute?
4. Could the mapper be rewritten as a simple table of field mappings?
5. If removed, would business behavior remain unchanged?

If any answer is "No", the mapper likely owns behavior it should not own.

Stop implementation and perform an architecture review.

---

## Extraction Rule

A mapper is not a reusable abstraction.

Do not extract mapper frameworks, mapper interfaces, generic converters, transformation pipelines, or reusable mapping engines unless duplication already exists.

Prefer explicit mapping functions.

Extract only after repeated duplication.

---

## Wire Contract Rule

Transport mappers own only the translation into the published wire contract.

They do not own the contract itself.

Changing a mapper must never silently change a transport contract.

Transport contracts remain owned by RFCs.

---

## Reference Implementation

Review `docs/architecture/GOLDEN_MAPPER.md` before writing any mapper.

Transport mappers must match the Golden Mapper reference (`apps/telegram/main.py::telegram_update_to_payload()`).

