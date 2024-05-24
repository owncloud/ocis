---
title: "28. Activity Service"
date: 2024-05-16T15:00:00+01:00
weight: 28
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/adr
geekdocFilePath: 0028-activities-service.md
---

* Status: draft
* Deciders: [@kobergj](https://github.com/kobergj), [@fschade](https://github.com/fschade)
* Date: 2024-05-16

## Context and Problem Statement

The user should be able to see all activities for a resource.
Besides the current resource, the user should also be able to decide if he wants to include child resource activities or not.

## Decision Drivers <!-- optional -->

* The user should be able to see all activities for a resource.
* The user should be able to decide if he wants to include child resource activities.
* Activities should be stored space efficiently.
* Activities should be stored in a way that they can be queried efficiently.
* Activities should stay in place even if the resource is gone.
* Activities reflect the state at a given point in time and not the current state.
* The Service should only store a configurable number of activities per resource.

## Considered Options

### Activity store

* Use a go-micro store to store the individual activities.
* Use a time series database to store the activities.
* Use a graph database to store the activities.
* Use a relational database to store the activities.
* Use the file system to store the activities.

### Activity format

* Normalize the activities before storing them.
* Only store relevant data to get the related event from the event-history service when needed, e.g.,
  ```go
    package pseudo_code

    import (
        "time"
    )

    type Activity struct {
        ResourceId string
        EventID string
        Depth int64
        Timestamp time.Time
    }
  ```
* Store the activity in a human-readable way e.g. "resource A has been shared with user B."
* Store each activity only on the resource itself.
* Store each activity only on the resource itself and all its parents.

## Decision Outcome

* Activity store:
  * Use a go-micro store to store the individual activities.
* Activity format:
  * Store each activity only on the resource itself and all its parents.
  * Only store event ids and get the related event from the event-history service when needed.

### Positive Consequences:

* Activity store (go-micro store):
  * Reuse existing technology.
  * We can use nats-js-kv store which already proved reliable in production.
  * No need to introduce any kind of new technology, e.g., a time series database, a relational database.
* Activity Format:
  * Having each activity stored on each resource (the resource itself and its parents)
    makes it easy to retrieve the timeline of activities for a resource and its children.
  * Only storing the event id and getting the related event from the event-history we benefit
    from the event-history services capabilities to store and query events.
  * Walking the resource tree from the resource to the root is a linear operation and can be done efficiently.

### Negative Consequences:

* Activity store:
  * Other database types might be more efficient for storing activities.
  * Using the go-micro-store only allows storing the activity in a key-value format.
* Activity Format:
  * Storing only the event ids and getting the related data from the event-history service when needed
    might introduce additional latency when querying activities.
  * Adding each event-id to each resource parent leads to a lot of duplicated data.

## Pros and Cons of the Options <!-- optional -->

* Activity store:
  * (PRO) Introducing a new database type might be more efficient for storing activities.
  * (CON) Introducing a new database type brings extra complexity and maintenance overhead.
  * (CON) Using the file system to store the activities might be inefficient and could be problematic especially in a distributed environment.
* Activity format:
  * (PRO) Normalizing the activities before storing them might make it easier and more efficient to query them.
  * (PRO) Storing each activity only on the resource itself is more space-efficient.
  * (CON) Storing each activity only on the resource itself increases the complexity of querying activities.
  * (CON) Storing each activity in a human-readable format is not space-efficient.

## Links <!-- optional -->

* [Story](https://github.com/owncloud/ocis/issues/8881)
