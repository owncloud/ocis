# Review API Reference (draft)

## `GET /reviews`

Returns the current user's pending reviews.

**Response**
```json
[
  { "id": "string", "title": "string", "submitter": "string", "dueDate": "string", "status": "pending|in_review|approved|rejected" }
]
```

## `POST /reviews/{id}/decision`

Submit a decision for a review.

**Request**
```json
{ "decision": "approve|reject|request_changes", "comment": "string (optional)" }
```

**Response:** `204 No Content`

## `GET /reviews/{id}/audit`

Returns the full decision history for a document.

[Add authentication details here]
