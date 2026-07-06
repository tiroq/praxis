# Golden Mapper Reference

**Normative architecture document. Use as the reference implementation for all transport adapters in Praxis.**

This document codifies the design of `apps/telegram/main.py:telegram_update_to_payload()` as the canonical example for transport mappers.

Future adapters (Slack, Email, Discord, HTTP, GitHub, CLI, etc.) should be compared directly against this reference to ensure compliance with the Pure Mapping Verification rules.

---

## Purpose

A **transport mapper** is a pure function that translates an external system's object representation into the internal Praxis wire contract.

This document is **normative, not descriptive**.

It defines what all transport mappers must be.

Deviation from this reference requires explicit architecture review and documented justification.

---

## Structural Properties

### The Reference Mapper

**File:** `apps/telegram/main.py`

**Function:** `telegram_update_to_payload(update: Update) -> dict`

**Properties:**

- ✅ Performs exactly one translation (Telegram Update → InputMessage dict)
- ✅ Has exactly one input (external `Update` object)
- ✅ Has exactly one output (internal wire contract `dict`)
- ✅ Contains no side effects (no publish, no log, no I/O)
- ✅ Contains no business logic (no classification, no enrichment, no decisions)
- ✅ Contains no infrastructure calls (no HTTP, no NATS, no storage)
- ✅ Contains no hidden state (stateless transformation)
- ✅ Contains no speculative abstractions (explicit, concrete function)

**Code Volume:** 13 lines (excluding docstring)

**Complexity:** Fits on one screen. Comprehensible in under one minute.

---

## Ownership

### What This Mapper Owns

The mapper owns **representation translation only**:

- Extract fields from external object
- Rename fields to match wire contract
- Convert primitive types
- Normalize formats (timezone, encoding, string conversion)
- Construct deterministic identifiers
- Provide transport-level default values

### What This Mapper Does NOT Own

The mapper explicitly does not own:

- ❌ **Transport** — connection, polling, publishing (owned by `run()` and `handle_message()`)
- ❌ **Business logic** — intent classification, enrichment, validation, business rules
- ❌ **Retries** — error recovery, retry logic, exponential backoff
- ❌ **Persistence** — caching, storage, databases, event sourcing
- ❌ **Orchestration** — coordination, workflows, state management
- ❌ **Contracts** — wire format definition (owned by RFCs and `internal/transport/nats.InputMessage`)

Clear ownership prevents responsibility from drifting into mappers as they evolve.

---

## Allowed Operations

These operations are **permitted** in mappers. Examples from the reference:

### 1. Copy Field

```python
"text": msg.text,
```

**Why:** Direct field pass-through. No transformation. Transparent.

---

### 2. Rename Field

```python
"message_id": str(msg.message_id),
```

**Why:** Maps external name to internal contract name. Explicit and localized.

---

### 3. Primitive Type Conversion

```python
chat_id = str(msg.chat_id)
message_id = str(msg.message_id)
```

**Why:** Type conversions are deterministic. No logic, just representation change.

---

### 4. Deterministic Formatting

```python
"timestamp": msg.date.astimezone(timezone.utc).isoformat(),
```

**Why:** Normalization is deterministic. Same input always produces same output. No external dependencies.

---

### 5. Deterministic Identifier Construction

```python
"id": f"tg-{chat_id}-{message_id}",
```

**Why:** Identifier is derived **solely from input fields**. Deterministic. Reconstructible from input.

---

### 6. Transport-Level Default Values

```python
"username": (user.username or "") if user else "",
"first_name": (user.first_name or "") if user else "",
```

**Why:** Provides safe fallback for optional external fields. Prevents null-pointer errors. Transport-specific precondition handling, not business logic.

---

## Forbidden Operations

These operations **MUST NOT** appear in mappers. Examples of why:

### ❌ Intent Classification

**Forbidden:**

```python
def telegram_update_to_payload(update: Update) -> dict:
    msg = update.message
    
    # FORBIDDEN: Intent classification
    if "buy" in msg.text.lower():
        intent = "purchase"
    elif "help" in msg.text.lower():
        intent = "support"
    else:
        intent = "unknown"
    
    return {
        "intent": intent,  # ← NOT from transport object
        ...
    }
```

**Why:** Classification is business logic. Belongs in a `ClassificationEngine`, not a mapper. Mappers are translation, not reasoning.

---

### ❌ Enrichment with External Knowledge

**Forbidden:**

```python
def telegram_update_to_payload(update: Update) -> dict:
    msg = update.message
    
    # FORBIDDEN: HTTP enrichment
    user_profile = requests.get(f"https://api.telegram.org/user/{msg.from_user.id}").json()
    
    return {
        "user_email": user_profile["email"],  # ← Where does it come from? Transport? No.
        ...
    }
```

**Why:** Enrichment is a separate concern. Belongs in an `EnrichmentEngine` or `Coordinator`. Mapper should not call HTTP.

---

### ❌ LLM Calls

**Forbidden:**

```python
def telegram_update_to_payload(update: Update) -> dict:
    msg = update.message
    
    # FORBIDDEN: LLM inference
    summary = openai.ChatCompletion.create(
        messages=[{"role": "user", "content": msg.text}]
    )["choices"][0]["message"]["content"]
    
    return {
        "summary": summary,  # ← Inferred, not from transport
        ...
    }
```

**Why:** Inference is business logic. Belongs in a `ReasoningEngine`, not a mapper. Mapper should not call LLMs.

---

### ❌ Repository Access

**Forbidden:**

```python
def telegram_update_to_payload(update: Update) -> dict:
    msg = update.message
    
    # FORBIDDEN: Storage access
    chat_history = chat_repository.get(msg.chat_id)
    
    return {
        "context": chat_history[-5:],  # ← From storage, not transport
        ...
    }
```

**Why:** Storage access is orchestration. Belongs in a `Coordinator`, not a mapper. Mapper should not access repositories.

---

### ❌ Retries

**Forbidden:**

```python
def telegram_update_to_payload(update: Update) -> dict:
    msg = update.message
    
    # FORBIDDEN: Retry logic
    for attempt in range(3):
        try:
            result = expensive_operation(msg.text)
            break
        except:
            if attempt < 2:
                time.sleep(2 ** attempt)
    
    return {
        "result": result,
        ...
    }
```

**Why:** Retries belong in a `ResilienceEngine` or middleware, not a mapper. Mapper is deterministic, not resilient.

---

### ❌ Metrics Collection

**Forbidden:**

```python
def telegram_update_to_payload(update: Update) -> dict:
    msg = update.message
    
    # FORBIDDEN: Metrics
    metrics.increment("telegram.messages_received")
    metrics.histogram("telegram.text_length", len(msg.text))
    
    return { ... }
```

**Why:** Metrics are observability infrastructure. Belongs in middleware or a `MetricsCollector`. Mapper should not publish telemetry.

---

### ❌ Business Validation

**Forbidden:**

```python
def telegram_update_to_payload(update: Update) -> dict:
    msg = update.message
    
    # FORBIDDEN: Business validation
    if len(msg.text) > 500:
        raise ValueError("Message too long for processing")
    
    if not msg.text.isprintable():
        raise ValueError("Non-printable characters not allowed")
    
    return { ... }
```

**Why:** Business rules belong in a `ValidationEngine` or `BusinessRuleEngine`. Mapper performs only transport-level validation (preconditions), not business logic.

---

## Field Traceability Table

**Template for reviewing future mappers.**

For `telegram_update_to_payload()`, every output field is traceable to input:

| Output Field | Source | Type | Justification |
|---|---|---|---|
| `"id"` | `f"tg-{msg.chat_id}-{msg.message_id}"` | Deterministic ID | Constructed from input fields only. Reconstructible. |
| `"source"` | `"telegram"` | Literal | Wire contract constant. Identifies transport source. |
| `"text"` | `msg.text` | Field copy | Direct pass-through from input. |
| `"timestamp"` | `msg.date.astimezone(timezone.utc).isoformat()` | Format conversion | Deterministic timezone normalization. Same input → same output. |
| `"metadata.chat_id"` | `str(msg.chat_id)` | Type conversion | Integer to string. Deterministic. |
| `"metadata.message_id"` | `str(msg.message_id)` | Type conversion | Integer to string. Deterministic. |
| `"metadata.username"` | `user.username or ""` if `user` else `""` | Field copy + default | Transport-level precondition: optional field. Safe fallback. |
| `"metadata.first_name"` | `user.first_name or ""` if `user` else `""` | Field copy + default | Transport-level precondition: optional field. Safe fallback. |

**Review Rule:** For every output field in a new mapper, the reviewer **MUST** trace it back to the source. If a field appears without a clear source, it is **not** a translation—it is business logic or enrichment.

---

## Architectural Review Checklist

**Use this checklist for every transport mapper.**

### Purity Verification

- [ ] Does the mapper accept only the external transport object?
- [ ] Does the mapper return only a dict (or the wire contract object)?
- [ ] Does the mapper call any I/O function (publish, HTTP, storage)?
- [ ] Does the mapper access any global state?
- [ ] Does the mapper perform any retry logic?
- [ ] Can this mapper be called 1,000,000 times with the same input and produce identical output?
- [ ] Is the mapper completely free of business logic?

### Structural Triviality

- [ ] Can every output field be traced directly to one or more input fields?
- [ ] Is every transformation deterministic?
- [ ] Would another engineer understand the mapper in under one minute?
- [ ] Could the mapper be rewritten as a simple table of field mappings?
- [ ] If removed, would business behavior remain unchanged?

### Responsibility

- [ ] Does the mapper perform exactly one translation?
- [ ] Has the mapper not accumulated unrelated behavior?
- [ ] Does the mapper remain the obvious place only for representation translation?

### Ownership

- [ ] Is the mapper completely free of infrastructure calls (HTTP, NATS, storage)?
- [ ] Is the mapper completely free of business logic?
- [ ] Is the mapper completely free of classification, enrichment, intent detection?

### Abstraction Quality

- [ ] No mapper framework invented?
- [ ] No base mapper class created?
- [ ] No mapper inheritance hierarchy?
- [ ] No reusable mapping utilities extracted prematurely?
- [ ] No mapper registry or factory?
- [ ] Explicit, concrete function?
- [ ] No over-engineering?

**Stop-and-Review Rule:** If any checkbox is unchecked, perform an architecture review before proceeding.

---

## Common Failure Modes

Mappers degrade over time through recognizable patterns. **Catch these early.**

### Pattern 1: "Just One More If"

**Example:**

```python
# Month 1: Pure mapper
def telegram_update_to_payload(update: Update) -> dict:
    return { "text": update.message.text, ... }

# Month 2: One business case
def telegram_update_to_payload(update: Update) -> dict:
    if update.message.text.startswith("/admin"):
        return None  # Silently drop
    return { "text": update.message.text, ... }

# Month 3: More cases
def telegram_update_to_payload(update: Update) -> dict:
    if update.message.text.startswith("/admin"):
        return None
    if len(update.message.text) > 1000:
        return None
    if not is_user_verified(update.from_user.id):
        return None
    return { "text": update.message.text, ... }

# Month 4: Mapper is now a filter, not a translator
```

**Why This Fails:** Business filtering belongs in `handle_message()`, not in the mapper.

**Fix:** Move conditions to caller. Mapper returns dict. Caller decides whether to publish.

---

### Pattern 2: Helper Extraction

**Example:**

```python
# Month 1: Inline
def telegram_update_to_payload(update: Update) -> dict:
    timestamp = update.message.date.astimezone(timezone.utc).isoformat()
    return { "timestamp": timestamp, ... }

# Month 2: Extract helper
def _format_timestamp(dt):
    return dt.astimezone(timezone.utc).isoformat()

def telegram_update_to_payload(update: Update) -> dict:
    return { "timestamp": _format_timestamp(update.message.date), ... }

# Month 3: Generalize helper
def _format_timestamp(dt, tz=timezone.utc):
    return dt.astimezone(tz).isoformat()

# Month 4: Move to utility module
# util/datetime.py
def format_timestamp(dt, tz=None, include_ms=False, locale=None):
    ...

# Month 5: Utility module becomes a framework
```

**Why This Fails:** Premature abstraction. Extract only after duplication.

**Fix:** Keep inline. Duplication in a second mapper is evidence for extraction, not hypothesis.

---

### Pattern 3: Shared Mapper Framework

**Example:**

```python
# Month 1: Telegram adapter
def telegram_update_to_payload(update: Update) -> dict:
    ...

# Month 2: Slack adapter
def slack_event_to_payload(event: SlackEvent) -> dict:
    ...

# Month 3: "These are similar. Let's extract."
class TransportMapper:
    def extract_field(self, obj, path):
        ...
    def rename(self, field, new_name):
        ...
    def map(self, transform_map):
        ...

class TelegramMapper(TransportMapper):
    def get_source(self):
        return "telegram"

class SlackMapper(TransportMapper):
    def get_source(self):
        return "slack"
```

**Why This Fails:** Framework turns two concrete implementations into an abstraction that owns logic. The framework becomes the place where responsibility accumulates.

**Fix:** Keep mappers explicit and concrete. Duplication is cheaper than abstraction.

---

### Pattern 4: Hidden Business Rules

**Example:**

```python
# Month 1: Pure mapping
def telegram_update_to_payload(update: Update) -> dict:
    return {
        "user_id": update.from_user.id,
        ...
    }

# Month 2: "We need to normalize user IDs"
def telegram_update_to_payload(update: Update) -> dict:
    # Prepend "tg-" to all user IDs for disambiguation
    return {
        "user_id": f"tg-{update.from_user.id}",
        ...
    }

# Month 3: "Some users are internal"
def telegram_update_to_payload(update: Update) -> dict:
    if update.from_user.id in INTERNAL_USER_IDS:
        return {
            "user_id": f"internal-{update.from_user.id}",
            "is_internal": True,
        }
    return {
        "user_id": f"tg-{update.from_user.id}",
        "is_internal": False,
    }

# Month 4: Mapper now contains business user classification logic
```

**Why This Fails:** User categorization is a business decision. Belongs in a `UserClassifier` or `UserEnricher`.

**Fix:** Mapper returns deterministic ID only. Classification happens elsewhere.

---

### Pattern 5: Infrastructure Leakage

**Example:**

```python
# Month 1: Pure mapping
def telegram_update_to_payload(update: Update) -> dict:
    return { "text": update.message.text, ... }

# Month 2: "We need to validate against our rules"
def telegram_update_to_payload(update: Update) -> dict:
    rules = validation_repository.get_rules(update.from_user.id)
    if not rules.validate(update.message.text):
        return None
    return { "text": update.message.text, ... }

# Month 3: "Also log violations"
def telegram_update_to_payload(update: Update) -> dict:
    rules = validation_repository.get_rules(update.from_user.id)
    if not rules.validate(update.message.text):
        logger.warning(f"validation failed: {update.from_user.id}")
        metrics.increment("validation.failures")
        return None
    return { "text": update.message.text, ... }

# Month 4: Mapper is now calling storage, logging, metrics
```

**Why This Fails:** Validation and logging belong in middleware or a `ValidationEngine`, not a mapper.

**Fix:** Keep mapper pure. Move all I/O to the caller (`handle_message()`).

---

## Golden Rule

> **A transport mapper should be boring.**
>
> If a mapper becomes interesting, it is almost certainly owning responsibility that belongs somewhere else.

Mappers are not the place for innovation, logic, or architectural sophistication. They are the place for mundane, transparent, deterministic field translation.

If you find yourself writing interesting code in a mapper, stop. Move the logic to its proper owner.

---

## Summary

This reference mapper exemplifies the Pure Mapping Verification rules:

- **Pure:** No side effects, no I/O, no state.
- **Structurally trivial:** One screen. One minute to understand.
- **Single responsibility:** Translate Telegram → InputMessage. Nothing else.
- **Clear ownership:** Owns representation only. Owns nothing else.
- **No accumulated behavior:** Single, explicit translation.

Use this mapper as the template for all future adapters.

When new adapters are proposed, compare directly against this reference. If they deviate, document why. If they cannot be justified, refactor to match this pattern.

One concrete, audited "golden mapper" is more valuable than dozens of abstract architecture rules.

---

**Document Status:** Normative. Update only after explicit architecture review and RFC amendment.

**Last Updated:** 2026-07-04

**Related:** [Praxis Architecture Guardian](../../.github/instructions/praxis-architecture-guardian.instructions.md), [Pure Mapping Verification](../../.github/instructions/praxis-architecture-guardian.instructions.md#pure-mapping-verification)
