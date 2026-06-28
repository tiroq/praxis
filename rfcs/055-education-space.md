

# RFC-055 Education Space

**Status:** Draft  
**Authors:** Tiroq + ChatGPT  
**Last Updated:** 2026-06-28

---

## 1. Summary

This RFC defines the **Education Space** as the RFC-050 Space specialization for learning, skill acquisition, curriculum planning, study workflows, educational progress, mentoring, children education, and lifelong development.

Education Space is not only a collection of courses or notes. It is an operating model for learning systems where goals, skills, lessons, practice, assessments, feedback, reviews, decisions, memory, and knowledge evolve over time.

The Education Space answers:

> How does Praxis represent learning as a structured, measurable, adaptive, and knowledge-preserving system?

Education Space uses the generic Space architecture from RFC-050 and specializes it for individual learning, family education, professional upskilling, and guided study programs.

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

This RFC is required before:

- RFC-056 Finance Space
- RFC-060 Testing Strategy
- RFC-061 Verification Scripts
- RFC-062 Benchmarking

RFC-050 defines Space as the bounded context. This RFC defines the Education Space specialization.

RFC-043 is especially important because Education Space relies on memory, knowledge consolidation, retrieval, and evidence-backed learning progress.

RFC-040 through RFC-042 are especially important because education agents, LLM routing, and prompt versioning must be governed and reproducible.

---

## 3. Goals

The goals of this RFC are to:

- Define Education Space as a specialization of the Space model.
- Model learning goals, skills, courses, lessons, exercises, assessments, feedback, and progress.
- Support adaptive learning plans and structured practice.
- Preserve learning history and educational decisions.
- Support self-directed learning, family education, mentoring, and professional upskilling.
- Provide agents for tutoring, curriculum planning, assessment, and review.
- Support evidence-based progress tracking rather than vague completion tracking.
- Enable explicit cross-space communication with Personal, Work, Product, Freelance, and Finance Spaces.

---

## 4. Non-Goals

This RFC does not:

- Define a full school management system.
- Replace official educational institutions.
- Define grading standards for every country.
- Define certification authority.
- Define child safety policy in implementation detail.
- Define every possible subject domain.
- Define UI screens.
- Define concrete learning content.

---

## 5. Education Space Philosophy

Learning is not content consumption.

Learning is a loop:

```text
Goal
↓
Plan
↓
Study
↓
Practice
↓
Assessment
↓
Feedback
↓
Reflection
↓
Adjustment
```

Education Space exists to preserve and improve this loop.

A good Education Space should answer:

- What is being learned?
- Why is it being learned?
- What evidence shows progress?
- Which concepts are weak?
- What practice is needed next?
- Which learning strategy works for this learner?
- Which goals should be adjusted?

Education Space treats learning as an adaptive, evidence-backed process.

---

## 6. Scope

Education Space covers:

- learning goals;
- skills and competencies;
- courses and curricula;
- lessons and study sessions;
- exercises and practice;
- assessments and quizzes;
- feedback and corrections;
- learning plans;
- study schedules;
- certifications and credentials;
- mentoring and tutoring;
- children education;
- professional upskilling;
- educational documents;
- learning memory and knowledge graph.

---

## 7. Education Identity Model

Education identity is multi-dimensional.

| Identity Dimension | Meaning |
|-------------------|---------|
| Learner | Person learning the subject. |
| Guardian | Parent or responsible adult where applicable. |
| Mentor | Human or agent guiding learning. |
| Teacher | Formal instructor or educator. |
| Program | Structured path of learning. |
| Course | Bounded instructional unit. |
| Subject | Area of knowledge or practice. |
| Skill | Specific capability to develop. |
| Level | Beginner, intermediate, advanced, etc. |
| Credential | Certificate, exam, or proof of competence. |

Education identity must support both adult self-learning and child/family education.

---

## 8. Learning Hierarchy

Education Space may represent hierarchy.

```text
Learning Goal
└── Program
    └── Course
        └── Module
            └── Lesson
                └── Exercise
                    └── Assessment
```

Hierarchy provides structure, but learning progress must be evidence-based rather than purely completion-based.

---

## 9. Core Canonical Objects

| Object | Description |
|--------|-------------|
| Learner | Person whose learning is being modeled. |
| Learning Goal | Desired educational outcome. |
| Skill | Capability to acquire or improve. |
| Competency | Measurable bundle of related skills. |
| Subject | Domain of knowledge or practice. |
| Program | Structured learning path. |
| Course | Bounded educational unit. |
| Module | Subdivision of a course. |
| Lesson | Learning unit. |
| Exercise | Practice task. |
| Assessment | Measurement of understanding or ability. |
| Feedback | Correction, evaluation, or guidance. |
| Study Session | Time-bounded learning activity. |
| Learning Plan | Planned sequence of learning activities. |
| Resource | Book, video, article, course, teacher, or tool. |
| Credential | Certificate, exam result, diploma, or proof. |
| Review | Periodic learning review. |
| Decision | Learning-related choice or direction. |

---

## 10. Education Agents

| Agent | Responsibility |
|------|----------------|
| Learning Planner | Builds learning paths and study schedules. |
| Tutor Agent | Explains concepts and adapts explanations. |
| Practice Generator | Creates exercises and drills. |
| Assessment Agent | Evaluates answers and detects weak concepts. |
| Feedback Agent | Provides corrections and improvement guidance. |
| Curriculum Reviewer | Reviews programs and course structure. |
| Parent Assistant | Helps manage child learning, homework, and routines. |
| Skill Gap Analyst | Detects missing prerequisites or weak skills. |
| Study Coach | Supports consistency, motivation, and reflection. |
| Credential Planner | Tracks exams, certificates, and requirements. |

Agents operate under Education Space policies and must preserve learner safety, privacy, and evidence traceability.

---

## 11. Learning Workflow Model

Education Space supports a general learning workflow.

```text
Goal
↓
Diagnose
↓
Plan
↓
Study
↓
Practice
↓
Assess
↓
Review
↓
Adjust
↓
Consolidate
```

This workflow is adaptive. A learner may return to diagnosis or practice when evidence shows gaps.

---

## 12. Skill Model

Skills are first-class objects.

A Skill may have:

- prerequisites;
- level;
- evidence requirements;
- practice history;
- assessment history;
- confidence score;
- retention score;
- related concepts;
- associated resources;
- linked credentials.

Skill mastery should be inferred from evidence, not from content completion alone.

---

## 13. Competency Model

A Competency groups skills into a meaningful capability.

Examples:

| Competency | Skills |
|-----------|--------|
| Python Automation | Python syntax, file IO, HTTP APIs, testing, packaging. |
| English Speaking | Vocabulary, pronunciation, listening, fluency. |
| Algebra Basics | Equations, functions, graphing, factoring. |
| Product Management | Discovery, prioritization, metrics, roadmap. |
| QA Engineering | test design, automation, CI, performance, resiliency. |

Competency progress depends on underlying skill evidence.

---

## 14. Curriculum Model

A Curriculum is a structured path through subjects, skills, resources, exercises, and assessments.

Curriculum may be:

- official;
- custom;
- AI-generated;
- teacher-created;
- parent-created;
- self-directed;
- imported from external platforms.

Curriculum must remain editable and reviewable.

---

## 15. Study Session Model

A Study Session is a time-bounded learning event.

It may include:

- learner;
- goal;
- subject;
- resource;
- exercises;
- notes;
- assessment;
- feedback;
- duration;
- focus level;
- outcome;
- next action.

Study Sessions produce Events and may update learning memory.

---

## 16. Practice Model

Practice is how skills become durable.

Practice may be:

- drill-based;
- project-based;
- problem-based;
- spaced repetition;
- simulation;
- guided exercise;
- freeform application.

Practice should be tied to skills and evidence.

---

## 17. Assessment Model

Assessment measures learning evidence.

Assessment types:

| Type | Purpose |
|------|---------|
| Diagnostic | Identify current level and gaps. |
| Formative | Guide ongoing learning. |
| Summative | Measure achievement at end of unit. |
| Self-assessment | Capture learner confidence and reflection. |
| Peer review | Human review by peer or mentor. |
| Agent review | AI-assisted evaluation with traceable criteria. |
| Credential exam | External formal evaluation. |

Assessments must preserve evidence, criteria, and result interpretation.

---

## 18. Feedback Model

Feedback converts assessment into improvement.

Feedback should include:

- what was correct;
- what was incorrect;
- why it matters;
- which concept is weak;
- what to practice next;
- confidence level;
- evidence reference;
- reviewer identity.

Feedback must not be only judgment. It must guide the next learning action.

---

## 19. Learning Plan Model

A Learning Plan organizes goals, skills, resources, sessions, and assessments over time.

A Learning Plan may include:

- target outcome;
- target date;
- weekly cadence;
- resources;
- milestones;
- practice schedule;
- review schedule;
- assessment checkpoints;
- adjustment rules.

Plans should change when evidence changes.

---

## 20. Progress Model

Progress is evidence-backed.

Progress signals include:

- completed sessions;
- completed exercises;
- assessment results;
- retention score;
- confidence score;
- feedback quality;
- repeated mistakes;
- time spent;
- project outputs;
- credential results.

Progress should distinguish activity from competence.

---

## 21. Education Memory Model

Education Memory includes:

- learning history;
- skill progress;
- repeated mistakes;
- preferred explanations;
- successful practice methods;
- learner motivation patterns;
- assessment history;
- feedback history;
- curriculum changes;
- teacher or parent observations.

Education Memory is scoped to Education Space and governed by Education Space privacy policies.

---

## 22. Education Knowledge Graph

The Education Knowledge Graph connects:

- learners;
- goals;
- skills;
- competencies;
- subjects;
- resources;
- lessons;
- exercises;
- assessments;
- feedback;
- credentials;
- mentors;
- study sessions.

Example relationships:

| Relationship | Meaning |
|-------------|---------|
| requires | Skill requires prerequisite Skill. |
| teaches | Resource teaches Skill. |
| practices | Exercise practices Skill. |
| assesses | Assessment assesses Skill. |
| improves | Feedback improves Skill. |
| belongs_to | Lesson belongs to Module or Course. |
| supports | Resource supports Learning Goal. |
| achieved_by | Credential achieved by Assessment. |

---

## 23. Education Policies

Education Space policies may include:

- learner privacy policy;
- child safety policy;
- guardian approval policy;
- AI tutoring policy;
- assessment integrity policy;
- content source policy;
- retention policy;
- notification policy;
- cross-space sharing policy;
- external platform integration policy.

Policies must be explicit and auditable.

---

## 24. Child and Family Learning Model

Education Space may support children and family learning.

This requires stricter rules:

- guardian visibility;
- age-appropriate content;
- limited automation;
- human approval for sensitive recommendations;
- separation from unrelated adult memory;
- careful sharing with school or tutors;
- auditability of AI-generated guidance.

Child learning data must not leak into other Spaces by default.

---

## 25. Education AI Governance

AI use in Education Space must be governed.

Governance rules:

- AI explanations must be age-appropriate where applicable.
- AI-generated assessments must be reviewable.
- AI feedback must be evidence-backed.
- AI should not fabricate learning progress.
- AI should distinguish uncertainty from mastery.
- AI must not replace guardian or teacher authority where required.
- Prompt versions must be traceable.
- Model routing must respect privacy and child safety policy.

---

## 26. Integrations

Education Space may integrate with:

- calendars;
- learning management systems;
- online courses;
- school portals;
- note-taking tools;
- document storage;
- flashcard systems;
- coding platforms;
- language learning platforms;
- exam systems;
- video platforms;
- reading trackers.

Integrations must normalize external records into Education Space contracts.

---

## 27. Projections

Education Space projections include:

| Projection | Purpose |
|-----------|---------|
| Today | Study sessions, homework, reviews, and practice. |
| This Week | Learning plan and scheduled assessments. |
| Skills | Skill graph, mastery, prerequisites, gaps. |
| Courses | Course progress and next lessons. |
| Practice | Exercises, drills, spaced repetition. |
| Assessments | Results, weak areas, feedback. |
| Child Learning | Guardian-facing learning view. |
| Credentials | Exams, certificates, and requirements. |

Projections are derived and rebuildable.

---

## 28. Cross-Space Communication

Education Space may communicate with:

| Target Space | Example |
|-------------|---------|
| Personal Space | Study schedule, family routines, child reminders. |
| Work Space | Professional upskilling linked to work goals. |
| Product Space | Product learning or market research skills. |
| Freelance Space | Skills needed for freelance positioning. |
| Finance Space | Education budget, courses, tuition, certification costs. |

Cross-space communication must use explicit references or events.

Learning memory must not leak into Work, Product, or Freelance Spaces by default.

---

## 29. Security Model

Education Space security must account for sensitive learning data.

Security dimensions:

- learner identity;
- guardian access;
- age sensitivity;
- content safety;
- assessment integrity;
- external platform access;
- AI agent permissions;
- cross-space sharing;
- audit logging.

Access control must be enforceable at learner, course, assessment, memory, and projection levels.

---

## 30. Lifecycle

Education Space lifecycle follows RFC-050.

```text
Draft
↓
Active
↓
Suspended
↓
Archived
```

An Education Space may be archived when a course, program, school year, or learning path is complete but should remain available for history and credentials.

---

## 31. Storage Mapping

| Store | Education Space Use |
|------|---------------------|
| Canonical Store | Learners, Skills, Goals, Courses, Assessments, Credentials. |
| Event Store | Study sessions, assessment events, feedback events. |
| Review Store | Learning reviews, curriculum reviews, assessment reviews. |
| Decision Store | Learning path decisions, course selection, credential decisions. |
| Action Store | Scheduled study actions, reminders, practice assignments. |
| Projection Store | Today, Skills, Courses, Assessments, Credentials views. |
| Search Index | Notes, lessons, resources, feedback, documents. |
| Vector Store | Semantic retrieval over learning content and notes. |
| Knowledge Graph | Learners, skills, prerequisites, resources, assessments. |
| Blob Store | PDFs, worksheets, certificates, recordings, assignments. |

---

## 32. Failure Modes

| Failure | Description |
|--------|-------------|
| Completion Illusion | Learner completes content without skill mastery. |
| Weak Evidence | Progress recorded without assessment evidence. |
| Forgotten Gaps | Repeated mistakes are not consolidated into memory. |
| Over-Automation | Agent makes sensitive educational decisions without review. |
| Child Privacy Leak | Child data leaks into unrelated Spaces. |
| Bad Recommendation | AI recommends unsuitable or unsafe content. |
| Assessment Drift | Assessments stop matching target skill. |
| Motivation Blindness | Plan ignores learner energy, interest, or frustration. |

---

## 33. Invariants

The following invariants must hold:

- Education objects are scoped to Education Space.
- Learning progress is evidence-backed.
- Skill mastery is not equal to content completion.
- Assessments preserve criteria and evidence.
- Feedback references assessed work or behavior.
- Education memory is private by default.
- Child learning data receives stricter protection.
- AI-generated learning guidance is traceable.
- Cross-space communication is explicit and auditable.
- Credentials reference evidence or external authority.

---

## 34. Architectural Consequences

The Education Space model enables:

- adaptive learning plans;
- evidence-backed progress;
- reusable learning memory;
- skill gap detection;
- AI-assisted tutoring with governance;
- family and child learning support;
- professional upskilling integration;
- long-term knowledge accumulation.

The cost is discipline: learning activities must produce evidence, feedback, or reflection to become useful memory.

---

## 35. Dependencies

Depends on:

- RFC-000 through RFC-054

Required before:

- RFC-056 Finance Space
- RFC-060 Testing Strategy
- RFC-061 Verification Scripts
- RFC-062 Benchmarking

---

## 36. Acceptance Criteria

This RFC can be accepted when:

- Education Space is defined as a specialization of RFC-050.
- Learning hierarchy is defined.
- Core education objects are defined.
- Education agents are defined.
- Skill and competency models are defined.
- Assessment and feedback models are defined.
- Education memory and knowledge graph are defined.
- Child and family learning constraints are defined.
- AI governance is defined.
- Cross-space communication is explicit.
- Invariants are agreed upon.

---

## 37. Decision Log

| Date | Decision | Author |
|------|----------|--------|
| 2026-06-28 | Initial draft of Education Space specialization. | Tiroq + ChatGPT |
| 2026-06-28 | Defined Education Space as evidence-backed adaptive learning system rather than content collection. | Tiroq + ChatGPT |

---

> **Learning is not content completion. Learning is evidence-backed change in capability.**