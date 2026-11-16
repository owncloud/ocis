Bugfix: Retry activity listing on empty response in acceptance tests

The `apiActivities` acceptance test suite queried the activities API
immediately after performing a file operation. Because activities are
written asynchronously via NATS, the event had sometimes not yet been
processed when the query arrived, returning an empty list and causing a
spurious `minItems` assertion failure.

The step that lists activities for a resource now retries up to ten
times with a 500 ms interval whenever the API returns HTTP 200 with an
empty `value` array, matching the retry pattern already used for the
share-listing step.
