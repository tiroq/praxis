

# RFC-056 Finance Space

**Status:** Draft  
**Authors:** Tiroq + ChatGPT  
**Last Updated:** 2026-06-28

---

## 1. Summary

This RFC defines the **Finance Space** as the RFC-050 Space specialization for personal finance, freelance income, business cash flow, assets, liabilities, budgets, financial goals, risk, taxes, investments, and long-term financial independence planning.

Finance Space is not only an expense tracker or portfolio dashboard. It is an operating model for financial decisions where income, expenses, assets, liabilities, goals, risks, taxes, invoices, payments, decisions, reviews, and projections are connected through explicit evidence and policies.

The Finance Space answers:

> How does Praxis represent money as a decision-driven, risk-aware, auditable, and goal-oriented system?

Finance Space uses the generic Space architecture from RFC-050 and specializes it for financial planning, financial memory, cash-flow visibility, and explicit financial decisions.

---

## 2. Relationship to Previous RFCs

This RFC depends on:

- RFC-000 Vision
- RFC-001 Principles
- RFC-002 Terminology
- RFC-003 Concept Model
- RFC-010 Capability Map
- RFC-011 Domain Model
- RFC-012 Artifact Model
- RFC-013 Event Model
- RFC-014 Identity & Representation Model
- RFC-015 Object Lifecycle Model
- RFC-020 Review System
- RFC-021 Decision Model
- RFC-022 State Machine
- RFC-030 System Architecture
- RFC-031 Service Contracts
- RFC-032 Data Flow
- RFC-033 Storage Model
- RFC-040 Agent Architecture
- RFC-041 LLM Routing
- RFC-042 Prompt Versioning
- RFC-043 Memory & Knowledge Graph
- RFC-050 Space Model
- RFC-051 Personal Space
- RFC-052 Work Space
- RFC-053 Product Space
- RFC-054 Freelance Space
- RFC-055 Education Space

This RFC is required before:

- RFC-060 Testing Strategy
- RFC-061 Verification Scripts
- RFC-062 Benchmarking

RFC-050 defines Space as the bounded context. This RFC defines the Finance Space specialization.

RFC-054 is especially important because Freelance Space produces invoices, payments, and income events that may flow into Finance Space.

RFC-020 and RFC-021 are especially important because high-impact financial actions require explicit review and decision records.

---

## 3. Goals

The goals of this RFC are to:

- Define Finance Space as a specialization of the Space model.
- Model accounts, income, expenses, assets, liabilities, budgets, taxes, investments, and goals.
- Support financial planning and decision history.
- Support cash-flow visibility and forecasting.
- Support financial independence and long-term planning.
- Support risk analysis and high-impact financial reviews.
- Support explicit cross-space communication with Personal, Work, Product, Freelance, and Education Spaces.
- Preserve financial memory while protecting sensitive financial data.

---

## 4. Non-Goals

This RFC does not:

- Provide financial, tax, legal, or investment advice.
- Replace professional accountants, advisors, or legal counsel.
- Define country-specific tax rules.
- Define brokerage execution.
- Define banking infrastructure.
- Define payment processing implementation.
- Define regulated financial services behavior.
- Guarantee investment outcomes.

Finance Space supports organization, analysis, evidence tracking, and decision governance. It does not make uncontrolled financial commitments.

---

## 5. Finance Space Philosophy

Money is not only transactions.

Money is a constraint system and a decision system.

Financial life connects:

```text
Income
↓
Cash Flow
↓
Budget
↓
Savings
↓
Assets
↓
Risk
↓
Goals
↓
Decisions
↓
Learning
```

Finance Space exists to make that system visible, auditable, and intentional.

A good Finance Space should answer:

- Where does money come from?
- Where does it go?
- What obligations exist?
- Which goals are funded?
- Which risks are exposed?
- Which decisions changed the trajectory?
- What should be reviewed before acting?
- What can be learned from past financial behavior?

Finance Space treats financial actions as reviewable decisions, not just transactions.

---

## 6. Scope

Finance Space covers:

- accounts;
- income;
- expenses;
- transfers;
- invoices;
- payments;
- budgets;
- cash flow;
- assets;
- liabilities;
- net worth;
- savings goals;
- emergency fund;
- investment tracking;
- tax-related records;
- subscriptions;
- insurance;
- financial documents;
- financial reviews;
- financial decisions;
- financial risk;
- financial independence planning.

---

## 7. Financial Identity Model

Finance identity is multi-dimensional.

| Identity Dimension | Meaning |
|-------------------|---------|
| Owner | Person or entity owning the financial Space. |
| Household | Family or group context where applicable. |
| Business | Freelance or company entity where applicable. |
| Account | Bank, wallet, brokerage, cash, or ledger account. |
| Currency | Currency of account, transaction, or goal. |
| Jurisdiction | Tax or legal context where applicable. |
| Counterparty | Person, company, platform, bank, or client. |
| Category | Income, expense, asset, liability, tax, etc. |
| Period | Month, quarter, year, project, or custom window. |

Financial identity must be explicit because many financial objects depend on ownership, jurisdiction, currency, and counterparty.

---

## 8. Finance Hierarchy

Finance Space may represent hierarchy.

```text
Financial Life
└── Goal
    └── Plan
        └── Budget
            └── Account
                └── Transaction
```

Alternative portfolio-oriented hierarchy:

```text
Net Worth
├── Assets
│   ├── Cash
│   ├── Investments
│   ├── Property
│   └── Receivables
└── Liabilities
    ├── Loans
    ├── Credit
    ├── Taxes
    └── Payables
```

Hierarchy provides structure, but financial truth must be evidence-backed by transactions, statements, contracts, invoices, or explicit manual records.

---

## 9. Core Canonical Objects

| Object | Description |
|--------|-------------|
| Account | Financial account, wallet, cash box, brokerage, or ledger. |
| Transaction | Money movement or financial event. |
| Income | Incoming money or receivable. |
| Expense | Outgoing money or payable. |
| Transfer | Movement between owned accounts. |
| Invoice | Request for payment. |
| Payment | Payment received or sent. |
| Budget | Planned allocation of money for a period or category. |
| Asset | Resource with economic value. |
| Liability | Financial obligation or debt. |
| Subscription | Recurring financial obligation. |
| Tax Record | Tax-relevant financial record. |
| Insurance | Risk-transfer policy or coverage record. |
| Financial Goal | Desired financial outcome. |
| Savings Plan | Plan to accumulate funds. |
| Investment | Position, fund, equity, crypto, or other tracked holding. |
| Review | Periodic financial review. |
| Decision | Financial decision or commitment. |
| Document | Statement, receipt, contract, tax document, or report. |

---

## 10. Finance Agents

| Agent | Responsibility |
|------|----------------|
| Budget Analyst | Reviews spending, budgets, and category drift. |
| Cash Flow Planner | Forecasts income, expenses, and liquidity. |
| Invoice Tracker | Tracks invoices, payments, and overdue items. |
| Subscription Auditor | Finds recurring charges and cancellation opportunities. |
| Tax Organizer | Groups tax-relevant records and reminders. |
| Risk Reviewer | Evaluates financial risk and exposure. |
| Goal Planner | Maps goals to savings and funding plans. |
| Investment Tracker | Summarizes portfolio state without executing trades. |
| Finance Critic | Challenges assumptions, weak plans, and risky decisions. |
| Document Curator | Organizes receipts, statements, contracts, and reports. |

Agents operate under Finance Space policies and must not execute high-impact financial actions without explicit human approval.

---

## 11. Financial Workflow Model

Finance Space supports a general financial workflow.

```text
Capture
↓
Classify
↓
Reconcile
↓
Analyze
↓
Review
↓
Decide
↓
Act
↓
Observe
↓
Learn
```

This workflow applies to budgeting, cash flow, tax organization, subscriptions, investments, and financial goals.

---

## 12. Account Model

Accounts are first-class objects.

An Account may have:

- owner;
- institution;
- account type;
- currency;
- opening balance;
- current balance;
- available balance;
- statement history;
- sync status;
- reconciliation status;
- privacy classification;
- linked documents.

Accounts may be manual, imported, or integration-backed.

---

## 13. Transaction Model

Transactions represent financial events.

A Transaction may include:

- account;
- amount;
- currency;
- date;
- counterparty;
- category;
- description;
- source;
- evidence document;
- tax relevance;
- recurrence flag;
- reconciliation status;
- linked object;
- notes.

Transactions must preserve source and evidence where possible.

---

## 14. Income Model

Income represents incoming money or expected receivables.

Income may originate from:

- salary;
- freelance invoices;
- product revenue;
- dividends;
- interest;
- sale of assets;
- refunds;
- reimbursements;
- gifts;
- other sources.

Income should be linked to source Space where applicable.

Example:

- Freelance Space Engagement → Invoice → Payment → Finance Space Income.

---

## 15. Expense Model

Expenses represent outgoing money or expected payables.

Expense categories may include:

- housing;
- food;
- transport;
- education;
- healthcare;
- tools;
- subscriptions;
- taxes;
- travel;
- business expenses;
- family expenses;
- investments;
- insurance.

Expense classification should be explainable and correctable.

---

## 16. Budget Model

A Budget defines intended allocation of money.

A Budget may include:

- period;
- category;
- planned amount;
- actual amount;
- variance;
- owner;
- goal link;
- policy;
- review cadence.

Budgets are planning tools, not truth. Actual transactions remain the evidence.

---

## 17. Cash Flow Model

Cash Flow tracks money timing.

Cash Flow includes:

- expected income;
- expected expenses;
- due dates;
- recurring obligations;
- invoice payment timing;
- account liquidity;
- forecast horizon;
- cash buffer.

Cash Flow is especially important for freelance, product, and household planning.

---

## 18. Net Worth Model

Net Worth is derived from Assets minus Liabilities.

```text
Net Worth = Assets - Liabilities
```

Net Worth is a projection, not a manually owned fact.

It must be traceable to asset and liability records.

---

## 19. Asset Model

Assets represent economic value.

Asset types may include:

- cash;
- bank balances;
- investments;
- property;
- vehicles;
- receivables;
- business assets;
- crypto assets;
- valuable personal assets.

Asset valuation may be exact, estimated, imported, or manually entered.

Valuation confidence must be visible where relevant.

---

## 20. Liability Model

Liabilities represent financial obligations.

Liability types may include:

- credit card debt;
- personal loans;
- mortgages;
- taxes payable;
- invoices payable;
- subscriptions;
- unpaid bills;
- contractual obligations.

Liabilities must preserve due dates and evidence where possible.

---

## 21. Investment Model

Investments are tracked for visibility, planning, and review.

Investment records may include:

- instrument;
- quantity;
- cost basis;
- current value;
- currency;
- account;
- allocation;
- risk category;
- thesis;
- review cadence;
- decision history.

Finance Space may summarize investments but must not execute trades unless a future RFC defines regulated execution boundaries.

---

## 22. Tax Model

Tax-related records must be organized but not interpreted as professional advice.

Tax model may include:

- jurisdiction;
- tax year;
- income records;
- deductible expenses;
- invoices;
- receipts;
- documents;
- filing deadlines;
- advisor notes;
- payment records.

Tax decisions should be reviewable and may require professional confirmation.

---

## 23. Financial Goal Model

Financial Goals represent intended outcomes.

Examples:

- emergency fund;
- family relocation;
- education fund;
- home purchase;
- runway target;
- debt repayment;
- financial independence;
- travel fund;
- product launch budget.

A Financial Goal may include:

- target amount;
- target date;
- currency;
- priority;
- funding source;
- progress;
- linked plan;
- risk assumptions;
- review cadence.

---

## 24. Financial Independence Model

Finance Space may support long-term financial independence planning.

This model may include:

- monthly expense baseline;
- safe runway estimate;
- income diversification;
- savings rate;
- asset allocation;
- risk scenarios;
- family obligations;
- relocation assumptions;
- tax assumptions;
- inflation assumptions.

All assumptions must be explicit and reviewable.

---

## 25. Risk Model

Financial risk must be represented explicitly.

Risk categories:

| Risk | Description |
|------|-------------|
| Liquidity Risk | Not enough cash when needed. |
| Income Risk | Income source may stop or decline. |
| Currency Risk | Currency movement affects obligations or assets. |
| Market Risk | Investment values may change. |
| Tax Risk | Incorrect or delayed tax handling. |
| Counterparty Risk | Client, bank, platform, or payer may fail. |
| Concentration Risk | Too much exposure to one source or asset. |
| Lifestyle Risk | Spending baseline exceeds sustainable level. |

Risk Reviews should be required for high-impact decisions.

---

## 26. Finance Memory Model

Finance Memory includes:

- spending patterns;
- income patterns;
- budget drift;
- recurring obligations;
- financial decisions;
- past mistakes;
- successful savings strategies;
- tax reminders;
- investment theses;
- risk reviews;
- family financial preferences;
- advisor notes.

Finance Memory is highly sensitive and private by default.

---

## 27. Finance Knowledge Graph

The Finance Knowledge Graph connects:

- accounts;
- transactions;
- budgets;
- goals;
- assets;
- liabilities;
- invoices;
- payments;
- clients;
- subscriptions;
- tax records;
- documents;
- decisions;
- reviews.

Example relationships:

| Relationship | Meaning |
|-------------|---------|
| paid_by | Payment paid by Account. |
| belongs_to | Transaction belongs to Category. |
| funds | Income funds Financial Goal. |
| derived_from | Net Worth derived from Assets and Liabilities. |
| linked_to | Invoice linked to Freelance Engagement. |
| supports | Document supports Tax Record. |
| reviewed_by | Financial Decision reviewed by Review. |
| exposes | Asset or income source exposes Risk. |

---

## 28. Finance Policies

Finance Space policies may include:

- privacy policy;
- approval policy;
- transaction classification policy;
- budget review policy;
- cash-flow warning policy;
- investment review policy;
- tax retention policy;
- document retention policy;
- AI agent permission policy;
- cross-space sharing policy;
- integration policy.

Policies must be explicit and auditable.

---

## 29. Finance AI Governance

AI use in Finance Space must be conservative.

Governance rules:

- AI must not provide uncontrolled financial, tax, legal, or investment advice.
- AI recommendations must be evidence-backed and clearly framed as analysis.
- High-impact actions require human approval.
- Investment-related outputs must preserve assumptions and uncertainty.
- Tax-related outputs should recommend professional review where needed.
- AI must not execute trades, payments, or transfers without explicit future policy support.
- Prompt versions and routing decisions must be traceable.
- Sensitive financial data must respect privacy and provider routing policies.

---

## 30. Integrations

Finance Space may integrate with:

- banks;
- payment processors;
- Stripe;
- Wise;
- PayPal;
- accounting systems;
- brokerage exports;
- crypto wallets;
- invoice systems;
- receipt scanners;
- tax portals;
- spreadsheets;
- freelance platforms;
- budgeting tools.

Integrations must normalize external records into Finance Space contracts.

---

## 31. Projections

Finance Space projections include:

| Projection | Purpose |
|-----------|---------|
| Cash Flow | Expected income, expenses, due dates, and liquidity. |
| Budget | Planned vs actual spending by category and period. |
| Net Worth | Assets, liabilities, and derived net worth. |
| Income | Income by source, Space, client, or period. |
| Expenses | Expenses by category, merchant, period, or goal. |
| Subscriptions | Recurring charges and renewal dates. |
| Taxes | Tax-relevant income, expenses, documents, and deadlines. |
| Goals | Financial goals, progress, funding gaps, and risks. |
| Investments | Holdings, allocation, thesis, and review status. |
| Risk | Financial risks, exposure, and review queue. |

Projections are derived and rebuildable.

---

## 32. Cross-Space Communication

Finance Space may communicate with:

| Target Space | Example |
|-------------|---------|
| Personal Space | Household budgets, travel funds, family goals. |
| Work Space | Compensation, benefits, work expenses. |
| Product Space | Product budget, revenue, costs, experiments. |
| Freelance Space | Invoices, payments, client revenue, project costs. |
| Education Space | Tuition, courses, certification budget. |

Cross-space communication must use explicit references or events.

Financial data must not leak into other Spaces by default.

Other Spaces may publish financial events into Finance Space only through approved policies.

---

## 33. Security Model

Finance Space security must be strict.

Security dimensions:

- owner identity;
- account access;
- data sensitivity;
- integration tokens;
- document privacy;
- tax records;
- payment records;
- investment records;
- AI provider routing;
- cross-space sharing;
- audit logging.

Access control must be enforceable at account, transaction, document, projection, memory, and integration levels.

---

## 34. Lifecycle

Finance Space lifecycle follows RFC-050.

```text
Draft
↓
Active
↓
Suspended
↓
Archived
```

Finance Space archival must preserve tax records, decisions, invoices, payments, and financial documents according to retention policies.

Deletion or export must be explicit and auditable.

---

## 35. Storage Mapping

| Store | Finance Space Use |
|------|-------------------|
| Canonical Store | Accounts, Transactions, Budgets, Assets, Liabilities, Goals. |
| Event Store | Transaction events, payment events, invoice events, review events. |
| Review Store | Budget reviews, risk reviews, investment reviews, tax reviews. |
| Decision Store | Financial decisions, budget decisions, investment decisions. |
| Action Store | Payment reminders, invoice follow-ups, review actions. |
| Projection Store | Cash flow, net worth, budget, tax, risk views. |
| Search Index | Receipts, statements, invoices, documents, notes. |
| Vector Store | Semantic retrieval over financial documents and notes. |
| Knowledge Graph | Accounts, goals, invoices, payments, risks, decisions. |
| Blob Store | Receipts, statements, contracts, tax documents, reports. |

---

## 36. Failure Modes

| Failure | Description |
|--------|-------------|
| Misclassification | Transaction assigned to wrong category or tax bucket. |
| Missing Income | Income source not captured or reconciled. |
| Hidden Liability | Debt, tax, or subscription not represented. |
| Cash Flow Surprise | Upcoming obligation not forecasted. |
| Duplicate Transaction | Imported transaction counted twice. |
| Currency Confusion | Amounts mixed across currencies incorrectly. |
| Unsafe AI Advice | Agent produces overconfident financial recommendation. |
| Privacy Leak | Sensitive financial data leaks across Spaces or providers. |
| Tax Evidence Gap | Tax-relevant record lacks supporting document. |
| Integration Drift | Bank or platform data diverges from Finance projections. |

---

## 37. Invariants

The following invariants must hold:

- Finance objects are scoped to Finance Space.
- Financial data is private by default.
- Transactions preserve source and evidence where possible.
- Net Worth is derived, not manually owned truth.
- Invoices and payments are traceable to source where applicable.
- Financial decisions are explicit and auditable.
- High-impact financial actions require human approval.
- AI-generated financial analysis is traceable and evidence-backed.
- Cross-space communication is explicit and policy-approved.
- Tax-related records preserve retention and evidence.
- Currency context is preserved for all monetary values.

---

## 38. Architectural Consequences

The Finance Space model enables:

- cash-flow visibility;
- budget discipline;
- explicit financial goals;
- financial decision history;
- income diversification analysis;
- freelance income integration;
- product finance integration;
- risk-aware planning;
- financial independence modeling;
- safer AI-assisted financial analysis.

The cost is discipline: financial records must preserve evidence, currency, ownership, and decision context to remain trustworthy.

---

## 39. Dependencies

Depends on:

- RFC-000 through RFC-055

Required before:

- RFC-060 Testing Strategy
- RFC-061 Verification Scripts
- RFC-062 Benchmarking

---

## 40. Acceptance Criteria

This RFC can be accepted when:

- Finance Space is defined as a specialization of RFC-050.
- Core financial objects are defined.
- Account, transaction, income, expense, asset, and liability models are defined.
- Budget, cash flow, net worth, and goal models are defined.
- Risk and tax models are defined.
- Finance agents are defined.
- Finance memory and knowledge graph are defined.
- AI governance is conservative and explicit.
- Cross-space communication is explicit.
- Security boundaries are defined.
- Invariants are agreed upon.

---

## 41. Decision Log

| Date | Decision | Author |
|------|----------|--------|
| 2026-06-28 | Initial draft of Finance Space specialization. | Tiroq + ChatGPT |
| 2026-06-28 | Defined Finance Space as evidence-backed financial decision system rather than expense tracker. | Tiroq + ChatGPT |

---

> **Money is not only tracked. Money is governed through evidence, decisions, risk, and goals.**