package svc

import (
	libregraph "github.com/owncloud/libre-graph-api-go"
)

type spacesSlice []*libregraph.Drive

// Len is the number of elements in the collection.
func (d spacesSlice) Len() int { return len(d) }

// Swap swaps the elements with indexes i and j.
func (d spacesSlice) Swap(i, j int) { d[i], d[j] = d[j], d[i] }

type spacesByName struct {
	spacesSlice
}
type spacesByLastModifiedDateTime struct {
	spacesSlice
}

// Less reports whether the element with index i
// must sort before the element with index j.
func (s spacesByName) Less(i, j int) bool {
	return *s.spacesSlice[i].Name > *s.spacesSlice[j].Name
}

// Less reports whether the element with index i
// must sort before the element with index j.
func (s spacesByLastModifiedDateTime) Less(i, j int) bool {
	// compare the items when both dates are set
	if s.spacesSlice[i].LastModifiedDateTime != nil && s.spacesSlice[j].LastModifiedDateTime != nil {
		return s.spacesSlice[i].LastModifiedDateTime.After(*s.spacesSlice[j].LastModifiedDateTime)
	}
	// move left item down if it has no value
	if s.spacesSlice[i].LastModifiedDateTime == nil && s.spacesSlice[j].LastModifiedDateTime != nil {
		return false
	}
	// move right item down if it has no value
	if s.spacesSlice[i].LastModifiedDateTime != nil && s.spacesSlice[j].LastModifiedDateTime == nil {
		return true
	}
	// fallback to name if no dateTime is set on both items
	return *s.spacesSlice[i].Name > *s.spacesSlice[j].Name
}
