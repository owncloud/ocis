# Mermaid Diagram Gallery

One example of each standard Mermaid.js diagram type, themed around the Aurora
review workflow.

## Flowchart

```mermaid
flowchart TD
    A[Document Submitted] --> B{Routing Queue}
    B --> C[Reviewer Assigned]
    C --> D{Decision}
    D -->|Approve| E[Published]
    D -->|Reject| F[Returned to Submitter]
    D -->|Request Changes| C
```

## Sequence Diagram

```mermaid
sequenceDiagram
    participant S as Submitter
    participant A as Review API
    participant R as Reviewer
    S->>A: Submit document
    A->>R: Notify: new review assigned
    R->>A: POST /reviews/{id}/decision
    A->>S: Notify: decision recorded
```

## Class Diagram

```mermaid
classDiagram
    class Review {
      +string id
      +string title
      +string status
      +submit()
      +decide(decision)
    }
    class Reviewer {
      +string name
      +approve(Review)
      +reject(Review)
    }
    Reviewer "1" --> "many" Review : decides
```

## State Diagram

```mermaid
stateDiagram-v2
    [*] --> Pending
    Pending --> InReview: assigned
    InReview --> Approved: approve
    InReview --> Rejected: reject
    InReview --> Pending: request changes
    Approved --> [*]
    Rejected --> [*]
```

## Entity Relationship Diagram

```mermaid
erDiagram
    REVIEWER ||--o{ REVIEW : decides
    SUBMITTER ||--o{ REVIEW : submits
    REVIEW ||--o{ COMMENT : has
    REVIEW {
      string id
      string status
      date dueDate
    }
```

## User Journey

```mermaid
journey
    title Reviewer's Day
    section Morning
      Check review queue: 5: Reviewer
      Read new document: 3: Reviewer
    section Afternoon
      Approve document: 5: Reviewer
      Leave feedback comment: 4: Reviewer
```

## Gantt Chart

```mermaid
gantt
    title Aurora Delivery Plan
    dateFormat  YYYY-MM-DD
    section Design
    Wireframes           :done, des1, 2026-06-15, 2026-06-25
    Review & sign-off     :done, des2, after des1, 5d
    section Engineering
    API implementation    :active, dev1, 2026-07-13, 25d
    Frontend integration   :dev2, after dev1, 15d
    section Rollout
    Pilot                :2026-08-24, 12d
```

## Pie Chart

```mermaid
pie title Reviews by Outcome (last 30 days)
    "Approved" : 62
    "Rejected" : 15
    "Changes Requested" : 23
```

## Quadrant Chart

```mermaid
quadrantChart
    title Reach vs Effort
    x-axis Low Effort --> High Effort
    y-axis Low Reach --> High Reach
    quadrant-1 Do First
    quadrant-2 Plan
    quadrant-3 Deprioritize
    quadrant-4 Delegate
    Bulk Approve: [0.3, 0.8]
    Mobile Layout: [0.7, 0.4]
    Audit Export: [0.5, 0.6]
```

## Requirement Diagram

```mermaid
requirementDiagram
    requirement auditTrail {
      id: 1
      text: every decision must be recorded
      risk: high
      verifymethod: test
    }
    element reviewAPI {
      type: service
    }
    reviewAPI - satisfies -> auditTrail
```

## Git Graph

```mermaid
gitGraph
    commit id: "init"
    branch feature/routing-queue
    checkout feature/routing-queue
    commit id: "add queue skeleton"
    commit id: "add tests"
    checkout main
    merge feature/routing-queue
    commit id: "release 0.1.0" tag: "v0.1.0"
```

## Mindmap

```mermaid
mindmap
  root((Aurora))
    Planning
      Charter
      Timeline
    Design
      Wireframes
      Architecture
    Engineering
      Review API
      Routing Queue
    Rollout
      Pilot
      GA
```

## Timeline

```mermaid
timeline
    title Aurora Roadmap
    2026-06 : Discovery
    2026-06 to 2026-07 : Design
    2026-07 to 2026-08 : Development
    2026-08 : Pilot Rollout
    2026-09 : General Availability
```
