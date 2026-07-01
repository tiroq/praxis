# ADR-003 — Internal Event Bus

**Status:** PROPOSED (not accepted — ratifies existing implementation, pending RFC record)
**Date:** 2026-07-01
**Authors:** Architecture
**Supersedes:** —
**Superseded by:** —
**Related:** ADR-002, RFC-013, RFC-030, RFC-031, RFC-032, RFC-033, RFC-022, ADQ-003

---

## Context

Praxis is an event-driven system. The RFCs make this a first-class architectural commitment,
not an implementation detail:

- **RFC-030 §8** ("Communication Model"): "Praxis uses **synchronous APIs only at system
  edges** … **Internally, services communicate asynchronously over an event bus.**"
- **RFC-030 §8.1** ("Communication Rules"): "Events are broadcast to all interested
  subscribers", "Events never expect responses or acknowledgments", and "Commands are
  point-to-point, targeting a single service owner."
- **RFC-030 §6** ("High-Level Architecture"): "Every stage communicates exclusively through
  immutable events, ensuring clear boundaries and traceability."
- **RFC-030 §14** ("Architectural Invariants"): "The runtime can be rebuilt from event
  streams" and "Runtime remains reconstructable from immutable Events."
- **RFC-032 §9** ("Event Flow"): "Every stage communicates via immutable Events. Events are
  append-only, timestamped, and carry correlation metadata."
- **RFC-013 §11** ("Event Immutability"): "Events are **append-only**. After being persisted
  as an Event Record, they can never be modified or deleted."
- **RFC-013 §16** ("Replay"): replay rebuilds "Projections, Read Models, Knowledge Graph
  relationships, Search Indexes, Analytics Views … by reprocessing the event log from a known
  state, without modifying original Event Records."

The RFCs therefore impose a precise set of demands on whatever technology carries events
between services: durable append-only streams, broadcast fan-out to many independent
subscribers, point-to-point command delivery, ordered replay from a known position, and
carriage of correlation/causation/trace metadata (RFC-013 §8, §13; RFC-031 §9; RFC-032 §12).

**What the RFCs do *not* do** is name a concrete broker. RFC-030, RFC-031, and RFC-032 all
describe "an event bus" abstractly. RFC-033 §4 ("Non-Goals") is explicit that concrete
engine selection is deferred: "Concrete implementation choices are deferred to implementation
RFCs and engineering design documents." No accepted RFC selects a message-bus technology.

**Evidence of an implicit, undocumented decision.** The implementation has already partially
committed to a broker without an ADR or RFC to record it:

| Location | Evidence | Implication |
|----------|----------|-------------|
| `go.mod` | Direct dependency `github.com/nats-io/nats.go v1.52.0` | NATS client is already a first-class dependency |
| `docker-compose.yml` | `nats:2.10-alpine` service on ports `4222`, `8222` | NATS is already the deployed broker |
| `internal/transport/nats/` | NATS pub/sub adapter package exists | Transport adapter already targets NATS |
| `services/worker` | "requires NATS" per Taskfile `run:worker` | Worker consumption already assumes NATS |
| ADR-002 (Architecture §) | Kestra publishes to `praxis.events.*` subjects, workers "consume without knowledge of their origin" | NATS subjects are already the integration seam |

This creates the exact failure mode the RFC process exists to prevent: an architecturally
load-bearing choice (the internal event backbone) is embedded in code and infrastructure but
is not recorded in any RFC or ADR. The current implementation additionally uses **core NATS**
(fire-and-forget pub/sub), which does **not** by itself satisfy the RFC-013 §11/§16 durability
and replay invariants — core NATS does not persist messages or support replay from a known
position. There is therefore both an **undocumented decision** and a **latent invariant gap**.

This ADR exists to (a) make the broker choice explicit and reviewable, and (b) select the
delivery mode that actually satisfies the durability, ordering, and replay invariants the RFCs
require.

---

## Decision

> **This ADR proposes adopting NATS with the JetStream persistence layer as the internal
> event backbone of Praxis. Core NATS request/reply is retained only for edge-synchronous
> and ephemeral control-plane messaging. The decision is PROPOSED, not accepted; no RFC is
> changed until this ADR is promoted to ACCEPTED. Until then, the existing NATS dependency
> stands as an implementation fact this ADR seeks to ratify and correct (core NATS →
> JetStream for the durable event path).**

The event bus is treated as **Infrastructure-layer** technology (RFC-030 §5.1). Per RFC-030
§5.2, "Business rules must never migrate into the Infrastructure layer", and per RFC-030
§13.1, "Infrastructure never owns business rules." The bus therefore transports events but
never interprets them.

---

## Architectural Role

The event bus is the transport that realizes three RFC-defined communication primitives. It
does not define them; RFC-030 §8.1 and RFC-031 §8–§9 do.

```mermaid
flowchart LR
    subgraph Domain/Application
      Ing[Ingestion]
      Und[Understanding]
      Can[Canonical Domain]
      Rev[Review]
      Dec[Decision]
      Act[Action]
      Lrn[Learning]
      Prj[Projection]
      Srch[Search]
    end
    subgraph Infrastructure
      Bus[(Event Bus\nNATS JetStream)]
    end

    Ing -- publish event --> Bus
    Und -- publish event --> Bus
    Can -- publish event --> Bus
    Rev -- publish event --> Bus
    Dec -- publish event --> Bus
    Act -- publish event --> Bus

    Bus -- fan-out --> Und
    Bus -- fan-out --> Can
    Bus -- fan-out --> Rev
    Bus -- fan-out --> Dec
    Bus -- fan-out --> Act
    Bus -- fan-out --> Lrn
    Bus -- fan-out --> Prj
    Bus -- fan-out --> Srch
```

The pipeline stages are defined by RFC-032 §6–§7 (Gateway → Ingestion → Understanding →
Canonical Domain → Review → Decision → Action → Learning → Projection). The bus is the wire
between them; the stages and their ownership are RFC-owned, not bus-owned.

---

## Invariants

These invariants derive from existing RFCs; this ADR selects a technology that satisfies them.

1. **Durability of the event path.** Events that carry the canonical pipeline forward MUST be
   persisted before acknowledgement (RFC-013 §11 "append-only … never be modified or
   deleted"; RFC-033 §6 Event Store "Rebuildable: No").
2. **Replay from a known position.** The bus (or the Event Store fed by it) MUST support
   reprocessing from a known offset without mutating history (RFC-013 §16; RFC-032 §16
   "Replay reconstructs all derived state from the event log").
3. **Broadcast fan-out.** A single event MUST be observable by many independent subscribers,
   each tracking its own position (RFC-030 §8.1 "Events are broadcast to all interested
   subscribers"; RFC-031 §9 "Events may be observed by many services").
4. **Point-to-point commands.** Commands MUST target exactly one owning service (RFC-030 §8.1
   "Commands are point-to-point, targeting a single service owner").
5. **Metadata carriage.** Every message MUST carry Correlation ID, Causation ID, and Trace ID
   (RFC-013 §8, §13; RFC-031 §9; RFC-032 §12) unchanged end-to-end.
6. **No business logic in the bus.** The bus performs no routing decisions based on business
   meaning (RFC-030 §5.2, §13.1). Subject hierarchies are transport addresses, not policy.
7. **Ownership isolation.** Services own their persistence; the bus is not a shared database
   (RFC-033 §31 "Services never read another service's database directly").

---

## Options Considered

### Option A — NATS JetStream *(PROPOSED)*

**Description.** NATS core provides subject-based pub/sub and request/reply. JetStream adds a
persistence layer: durable append-only streams, consumer cursors (durable/ephemeral,
push/pull), replay from sequence/time, at-least-once delivery with acknowledgements,
deduplication windows, and per-subject stream retention. Praxis already depends on NATS
(`nats.go v1.52.0`, `nats:2.10-alpine`), so this option upgrades the *delivery mode* rather
than introducing a new dependency.

**Fits.** RFC-030 §8/§8.1 (subject broadcast + point-to-point via distinct subjects/queue
groups). RFC-013 §11/§16 (JetStream streams are append-only and replayable). RFC-032 §9/§12
(headers carry correlation/causation/trace metadata). RFC-033 §6 (JetStream streams back the
Event Store; derived stores rebuild via replay). RFC-030 §5.1 (Infrastructure layer).

**Conflicts.** Core NATS (currently used) does not satisfy invariants 1–2 on its own; adoption
requires migrating the durable event path to JetStream. Exactly-once is not native (see
"exactly-once vs at-least-once" below).

**Advantages.** Single technology for both request/reply (edges) and durable streaming
(internal), reducing operational surface. Lightweight single binary; trivial local dev via
Docker Compose (already present). Native subject hierarchy maps cleanly to `praxis.events.*`
(already used by ADR-002). Built-in replay by sequence/time. Low latency. No ZooKeeper/KRaft
equivalent to operate.

**Disadvantages.** Smaller ecosystem than Kafka. JetStream clustering and stream sizing
require deliberate capacity planning. Exactly-once must be achieved via idempotent consumers +
dedup windows, not a broker guarantee.

**Operational impact.** One broker to run for both edge and internal messaging. JetStream adds
a persistence tier (file/memory store) that must be sized, backed up, and monitored. Clustering
uses RAFT; three-node quorum for HA.

**Implementation complexity.** Low incremental: the client dependency and adapter package
already exist. Work is confined to enabling JetStream streams/consumers and adding
idempotency, not integrating a new system.

**Long-term maintainability.** High. One dependency, one wire protocol, one mental model.
Subjects are stable addresses; consumers evolve independently. Aligns with the "lightest
full-featured" philosophy already stated for infrastructure in ADR-002.

**Verdict.** Best fit. Ratifies and corrects the existing implicit choice, satisfies every
RFC-013/030/032 invariant, and minimizes operational and cognitive surface.

---

### Option B — Apache Kafka (or Redpanda)

**Description.** A partitioned, log-based streaming platform with strong ordering per partition,
long retention, consumer groups, and a large ecosystem (Connect, Streams, Schema Registry).

**Fits.** RFC-013 §11/§16 (the log *is* an append-only, replayable store — arguably the
strongest fit for event sourcing). RFC-030 §8.1 broadcast (consumer groups) and RFC-032 §12
(headers). RFC-033 §6 Event Store maps naturally to Kafka topics.

**Conflicts.** RFC-030 §8.1 point-to-point commands and request/reply are awkward on Kafka
(no native request/reply; requires correlation topics). Introduces a **new** heavyweight
dependency not currently present, contradicting the minimal-infrastructure posture of ADR-002.

**Advantages.** Best-in-class durability and replay. Very high throughput. Mature tooling,
schema registry, exactly-once semantics (EOS) within Kafka boundaries. Strong per-partition
ordering.

**Disadvantages.** Operational weight (brokers + KRaft/ZooKeeper; Redpanda reduces but does
not eliminate this). Ordering is per-partition, so global ordering requires careful key design.
Request/reply and command point-to-point patterns need extra machinery. Heavier local-dev
footprint than NATS.

**Operational impact.** Highest of the candidates. Partition planning, retention tuning,
rebalancing, and broker operations become a standing concern. Two messaging paradigms would
likely coexist (Kafka for streams, something else for request/reply).

**Implementation complexity.** High. New client, new adapter, new local-dev stack; migration
from the existing NATS dependency; probable second technology for edge request/reply.

**Long-term maintainability.** Medium. Powerful but heavy; the operational tax is permanent and
disproportionate to a single-operator personal-scale system.

**Verdict.** Strongest raw event-sourcing substrate, but its operational weight and poor
request/reply fit make it disproportionate. Rejected for current scale; revisitable if
throughput or multi-tenant scale demands it.

---

### Option C — RabbitMQ

**Description.** A mature AMQP broker with rich routing (exchanges, bindings), acknowledgements,
dead-letter queues, and (via streams plugin) a log abstraction.

**Fits.** RFC-030 §8.1 broadcast (fanout/topic exchanges) and point-to-point commands (direct
exchanges/queues) — the routing model maps cleanly to both primitives. RFC-031 §8 command
delivery.

**Conflicts.** Classic queues are consume-and-delete, not append-only replayable logs; RFC-013
§16 replay requires the Streams plugin, which is less mature than JetStream/Kafka logs.
Introduces a new dependency.

**Advantages.** Excellent, flexible routing. Battle-tested reliability. Clear command/queue
semantics. Good management UI.

**Disadvantages.** Replay/event-sourcing is not its native strength. Two abstractions (queues
vs streams) to reason about. Erlang/OTP operational model unfamiliar to a Go/Python stack. New
dependency versus the existing NATS.

**Operational impact.** Medium. Broker + management plane; streams plugin adds another
subsystem to operate for the replay path.

**Implementation complexity.** Medium-high. New client and adapter; replay path needs the
streams plugin; migration off NATS.

**Long-term maintainability.** Medium. Solid but the replay story is bolted on, splitting the
mental model.

**Verdict.** Great message router, weaker event-sourcing substrate. Rejected because RFC-013
§16 replay is a core invariant, not an add-on.

---

### Option D — Redis Streams

**Description.** Redis' append-only stream type with consumer groups, IDs, and range reads,
optionally persisted via AOF/RDB.

**Fits.** RFC-013 §11 append-only IDs; RFC-030 §8.1 fan-out via consumer groups; low latency.

**Conflicts.** Durability depends on AOF/RDB configuration and is weaker than a purpose-built
log; RFC-013 §11 ("never be modified or deleted") and RFC-033 §6 (Event Store "Rebuildable:
No") demand strong durability guarantees. Memory-first design constrains long retention needed
for replay (RFC-013 §16).

**Advantages.** Extremely simple, very low latency, familiar. Doubles as cache (RFC-033 §6
Cache category). Trivial local dev.

**Disadvantages.** Durability and retention are not first-class for a system-of-record event
log. Long-horizon replay is memory-bound and operationally awkward. No native request/reply
beyond ad-hoc patterns.

**Operational impact.** Low to run, but risky as the *authoritative* event backbone; likely
needs a separate durable store, defeating the single-bus goal.

**Implementation complexity.** Low to start, higher to make durable/replayable to RFC standard.

**Long-term maintainability.** Medium-low as an event backbone; excellent as a cache.

**Verdict.** Excellent cache and ephemeral-stream tool; unsuitable as the durable
system-of-record event backbone. Rejected for the primary role; may serve RFC-033 §6 Cache.

---

### Option E — Google Cloud Pub/Sub (managed)

**Description.** Fully managed, horizontally scalable pub/sub with at-least-once delivery,
optional ordering keys, and replay via retained acknowledged messages / snapshots.

**Fits.** RFC-030 §8.1 broadcast; RFC-013 §16 replay (message retention + seek); effectively
infinite scale with zero broker operations.

**Conflicts.** Managed cloud dependency contradicts the local-first, Docker-Compose,
self-hostable posture evidenced by `docker-compose.yml`, `infra/`, and ADR-002's "run locally
in Docker Compose" rationale. Introduces vendor lock-in and network egress cost. No native
request/reply.

**Advantages.** No broker to operate. Elastic scale. Strong durability SLAs.

**Disadvantages.** Cloud lock-in; not runnable offline/local; per-message and egress cost;
latency higher than an in-cluster broker; ordering keys constrain throughput. Contradicts the
self-hosted infrastructure model.

**Operational impact.** Near-zero infra ops, but a hard external dependency and billing
surface; unusable in local/offline development.

**Implementation complexity.** Medium; new client/adapter and a divergent local-dev emulator.

**Long-term maintainability.** Medium; low ops but strategic lock-in and cost exposure.

**Verdict.** Attractive at large managed scale but conflicts with the self-hosted, local-first
architecture. Rejected for a single-operator, self-hosted system.

---

### Option F — No broker (in-process / direct calls)

**Description.** Services call each other directly (HTTP/gRPC) or run in-process; no
asynchronous bus.

**Fits.** Nothing in the async model. It would satisfy only a degenerate monolith.

**Conflicts.** Directly violates RFC-030 §8 ("Internally, services communicate asynchronously
over an event bus"), §8.1 (broadcast, no acknowledgement), §14 ("rebuilt from event streams"),
RFC-013 §11/§16 (append-only, replay), and RFC-032 §9/§16 (event flow, replay). There is no
event log to replay.

**Advantages.** Zero infrastructure; simplest to start.

**Disadvantages.** Eliminates durability, replay, fan-out, temporal decoupling, and audit — the
system's defining properties. Reintroduces tight coupling the architecture forbids.

**Operational impact.** None to run; catastrophic to the architecture's guarantees.

**Implementation complexity.** Low now, unbounded rework later to reintroduce the bus.

**Long-term maintainability.** Lowest; structurally incompatible with the RFCs.

**Verdict.** Non-viable. Included only to demonstrate the async bus is load-bearing, not
optional.

---

### Comparison Matrix

| Criterion | NATS JetStream | Kafka/Redpanda | RabbitMQ | Redis Streams | Google Pub/Sub | No broker |
|---|---|---|---|---|---|---|
| Durable append-only log (RFC-013 §11) | ✅ | ✅✅ | ⚠️ streams plugin | ⚠️ AOF-dependent | ✅ | ❌ |
| Replay from position (RFC-013 §16) | ✅ | ✅✅ | ⚠️ | ⚠️ | ✅ | ❌ |
| Broadcast fan-out (RFC-030 §8.1) | ✅ | ✅ | ✅ | ✅ | ✅ | ❌ |
| Point-to-point commands (RFC-030 §8.1) | ✅ | ⚠️ | ✅✅ | ⚠️ | ⚠️ | ✅ (sync) |
| Native request/reply (edges) | ✅ | ❌ | ⚠️ | ❌ | ❌ | ✅ |
| Metadata headers (RFC-032 §12) | ✅ | ✅ | ✅ | ⚠️ | ✅ | n/a |
| Ordering | ✅ per-stream | ✅ per-partition | ⚠️ | ✅ per-stream | ⚠️ keys | n/a |
| Exactly-once | ⚠️ idempotent | ✅ EOS (in-boundary) | ⚠️ | ⚠️ | ⚠️ | n/a |
| Operational weight | 🟢 low | 🔴 high | 🟡 medium | 🟢 low | 🟢 managed | 🟢 none |
| Local Docker-Compose dev | ✅ present | ⚠️ heavy | ✅ | ✅ | ❌ | ✅ |
| Already a dependency | ✅ | ❌ | ❌ | ❌ | ❌ | n/a |
| Self-hostable / offline | ✅ | ✅ | ✅ | ✅ | ❌ | ✅ |

---

## Proposed Decision

**Adopt NATS JetStream (Option A)** as the durable internal event backbone; retain **core NATS
request/reply** for edge-synchronous and ephemeral control messaging.

**Why A over the alternatives.**

- **Over Kafka (B):** Kafka is the strongest raw log, but its operational weight and weak
  native request/reply are disproportionate for a single-operator, self-hosted, personal-scale
  system. JetStream satisfies the same RFC-013 §11/§16 invariants at a fraction of the ops cost
  while also covering RFC-030 §8.1 point-to-point/request-reply in one technology.
- **Over RabbitMQ (C) and Redis Streams (D):** both bolt replay onto a non-log core; RFC-013
  §16 replay is a first-class invariant, so a native replayable log (JetStream) is preferred.
  Redis is retained instead for the RFC-033 §6 Cache category.
- **Over Google Pub/Sub (E):** managed cloud conflicts with the local-first, self-hosted,
  Docker-Compose posture evidenced across `infra/`, `docker-compose.yml`, and ADR-002.
- **Over No broker (F):** structurally violates RFC-030 §8/§14 and RFC-013 — non-viable.
- **Decisive factor:** NATS is *already* the de-facto dependency (`nats.go`, `nats:2.10`,
  `internal/transport/nats/`). Choosing A ratifies an existing implicit decision and closes the
  invariant gap by moving the durable path from core NATS to JetStream — the smallest correct
  change.

**Delivery-mode policy.**

| Message class | Transport | Delivery | Rationale |
|---|---|---|---|
| Canonical pipeline events (RFC-032 §9) | JetStream stream | at-least-once + idempotent consumers | RFC-013 §11/§16 durability + replay |
| Commands to a single owner (RFC-030 §8.1) | JetStream work-queue / core queue group | at-least-once, single delivery | RFC-031 §8 point-to-point |
| Edge synchronous queries | core NATS request/reply | best-effort, no persistence | RFC-030 §8 "synchronous APIs only at system edges" |
| Ephemeral control/health | core NATS pub/sub | fire-and-forget | Not part of the event log |

**Exactly-once vs at-least-once.** The bus provides **at-least-once**. Effective
exactly-once processing is achieved at the consumer via idempotency keyed on Event ID
(RFC-013 §8) and JetStream deduplication windows. This is consistent with RFC-013 §11's
correction model ("a new event is appended that references the original") — the system never
relies on in-place mutation or a broker-level exactly-once guarantee.

---

## Consequences

### Positive
- Satisfies RFC-013 §11/§16 durability and replay, RFC-030 §8.1 broadcast + point-to-point,
  and RFC-032 §12 metadata carriage with a single technology.
- Ratifies the existing dependency; no new broker to learn or operate.
- One wire protocol for both edge request/reply and internal streaming.
- Local-first: unchanged Docker-Compose developer experience.

### Negative
- Requires migrating the durable path from core NATS to JetStream (streams, consumers,
  idempotency), which the current code does not yet do.
- Exactly-once must be engineered at consumers rather than assumed from the broker.
- JetStream persistence tier adds sizing, retention, and backup responsibilities.

### Trade-offs
- Accepts a smaller ecosystem than Kafka in exchange for far lower operational weight and
  native request/reply.
- Accepts consumer-side idempotency work in exchange for a single-technology footprint.

### Future flexibility
- Subjects are stable addresses; consumers can be added/replaced without producer changes
  (RFC-030 §8.1). If scale ever demands Kafka-class throughput, the Event Store abstraction
  (RFC-033 §6) allows swapping the stream substrate behind service-owned repositories without
  touching domain logic (RFC-030 §13.1).

### Migration cost
- Enable JetStream streams for `praxis.events.*`; convert existing core-NATS publishes on the
  canonical path to JetStream publishes; add durable pull consumers per service; add
  idempotency keyed on Event ID. Edge request/reply and ephemeral control paths remain on core
  NATS unchanged.

### Operational impact
- Run NATS with JetStream enabled (file store); plan a 3-node cluster for HA. Monitor stream
  depth, consumer lag, redelivery counts, and storage utilization. Back up JetStream file
  store alongside Postgres (per ADR-002 backup pipeline candidate).

### Development impact
- Producers/consumers gain durable semantics and must be idempotent. The bus remains behind
  the `internal/transport/nats` adapter so domain/application code depends on a port, not on
  NATS APIs (RFC-030 §13.1).

### Testing impact
- New Replay tests (RFC-060 §8 "Replay Tests"; RFC-061 §17 Replay layer) assert derived state
  rebuilds from the stream. New Invariant tests assert append-only and idempotent consumption.

### Performance impact
- JetStream adds persistence latency versus core NATS fire-and-forget, but stays well within
  interactive budgets for a single-operator system; benchmarking is covered by RFC-062
  (latency/reliability dimensions).

### Failure modes
- Broker unavailability: producers block/queue; consumers resume from last ack on recovery.
- Duplicate delivery: neutralized by consumer idempotency (Event ID) + dedup window.
- Stream storage exhaustion: mitigated by retention policy + monitoring; canonical Event Store
  retention follows RFC-033 §6 ("Rebuildable: No" → retained).
- Poison messages: routed to a dead-letter subject after max-deliver; never silently dropped
  (preserves RFC-013 §11 auditability).

---

## Required RFC Edits (if this ADR is accepted)

| RFC | Required change | Scope |
|-----|-----------------|-------|
| RFC-030 | Add a note that the event bus is realized by NATS JetStream (Infrastructure layer), without embedding vendor detail in domain sections. | System architecture |
| RFC-031 | Specify delivery semantics per message class (events at-least-once + idempotent; commands single-owner) referencing the bus. | Service contracts |
| RFC-032 | Note that canonical events are carried on durable streams and that replay reads from stream positions. | Data flow |
| RFC-033 | Clarify that the Event Store (§6) is backed by durable streams and how replay maps to stream positions. | Storage model |
| RFC-013 | Add an implementation note that at-least-once + Event-ID idempotency provides effective exactly-once; no in-place mutation. | Event model |
| RFC-060/061 | Add explicit Replay and idempotency invariant tests to the verification manifest. | Testing / Verification |

---

## Implementation Impact

### Safe to implement immediately
- Keeping the existing NATS dependency and `internal/transport/nats` adapter.
- Defining `praxis.events.*` subject conventions (already used by ADR-002).
- Adding correlation/causation/trace headers to all messages (RFC-013 §8, §13).

### Blocked until this ADR is accepted
- Enabling JetStream streams/consumers as the *authoritative* Event Store path (RFC-033 §6).
- Declaring at-least-once + idempotency as the canonical delivery contract in RFC-031.
- Any code that assumes replay-from-position semantics (needs JetStream, not core NATS).

---

## Verification Impact

### Existing verification affected
- Any test that assumes fire-and-forget core NATS must be updated for durable delivery and
  redelivery on the canonical path.
- Contract tests (RFC-061 §14) must assert that domain/application code imports the bus adapter
  port, not `nats.go` directly.

### New verification required
- **Replay test** (RFC-060 §8; RFC-061 §17): rebuild a projection purely from the stream and
  assert equivalence to live state (RFC-013 §16; RFC-032 §16).
- **Append-only invariant** (RFC-061 §15): assert no delete/overwrite on the event stream
  (RFC-013 §11).
- **Idempotency test**: duplicate delivery of an Event ID produces a single effect.
- **Boundary test** (RFC-061 §14): only the transport adapter imports NATS; no domain package
  references `nats.go`.

### Testing changes
- Add a JetStream-backed integration harness; use ephemeral streams per test.

### Coverage impact
- New Replay/Invariant/Boundary categories increase verification breadth on the transport path;
  aligns with RFC-060 §8 categories.

### Acceptance criteria
- Derived state is reconstructable from the stream (RFC-013 §16 / RFC-032 §16).
- No domain package imports a broker SDK.
- Duplicate deliveries are neutralized; poison messages are dead-lettered, never dropped.

---

## Rejected Alternatives

- **Kafka/Redpanda (B):** operationally disproportionate for single-operator scale and weak
  native request/reply; JetStream meets the same invariants at lower cost.
- **RabbitMQ (C):** replay is a bolt-on (streams plugin), but RFC-013 §16 replay is a
  first-class invariant.
- **Redis Streams (D):** durability/retention insufficient for a system-of-record log;
  retained instead for the RFC-033 §6 Cache role.
- **Google Pub/Sub (E):** managed-cloud lock-in conflicts with the self-hosted, local-first
  posture (ADR-002, `infra/`, `docker-compose.yml`).
- **No broker (F):** violates RFC-030 §8/§14 and RFC-013; eliminates the system's defining
  properties.

---

## Open Questions

1. **Event Store boundary:** Is the JetStream stream itself the authoritative Event Store
   (RFC-033 §6), or is it a transport feeding a separate durable Event Store? This affects
   backup and replay tooling and must be settled with RFC-033.
2. **Retention horizon:** How long must canonical event streams be retained for replay
   (RFC-013 §16) versus offloaded to cold storage? Requires a retention policy per stream.
3. **Stream topology:** One stream per pipeline stage, per subject-tree, or a single global
   stream partitioned by subject? Impacts ordering guarantees and consumer design.
4. **Clustering target:** Is HA (3-node JetStream cluster) in scope for the personal-scale
   deployment, or is single-node with backups acceptable for Phase 1–2?
5. **Command transport:** Should commands use JetStream work-queue streams or core NATS queue
   groups? Trade-off between durability and latency for point-to-point (RFC-030 §8.1).

---

## References

- [docs/adr/ADR-002-external-workflow-orchestrator.md](ADR-002-external-workflow-orchestrator.md) — orchestrator publishes to `praxis.events.*`
- [rfcs/013-event-model.md](../../rfcs/013-event-model.md) — event immutability (§11), envelope (§8), correlation/causation (§13), replay (§16)
- [rfcs/030-system-architecture.md](../../rfcs/030-system-architecture.md) — communication model (§8, §8.1), layers (§5.1–§5.2), invariants (§14), boundaries (§13.1)
- [rfcs/031-service-contracts.md](../../rfcs/031-service-contracts.md) — command rules (§8), event structure (§9)
- [rfcs/032-data-flow.md](../../rfcs/032-data-flow.md) — pipeline (§6–§7), event flow (§9), correlation model (§12), invariants (§16)
- [rfcs/033-storage-model.md](../../rfcs/033-storage-model.md) — storage categories (§6), non-goals (§4), ownership (§30–§31)
- [rfcs/022-state-machine.md](../../rfcs/022-state-machine.md) — runtime state machines driven by events
- [docs/ARCHITECTURE_DECISION_QUEUE.md](../ARCHITECTURE_DECISION_QUEUE.md) — ADQ-003 (workflow model, orchestration boundary)
