package svc

import (
	"strings"

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
	return strings.ToLower(*s.spacesSlice[i].Name) > strings.ToLower(*s.spacesSlice[j].Name)
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
	return strings.ToLower(*s.spacesSlice[i].Name) > strings.ToLower(*s.spacesSlice[j].Name)
}

type userSlice []*libregraph.User

// Len is the number of elements in the collection.
func (d userSlice) Len() int { return len(d) }

// Swap swaps the elements with indexes i and j.
func (d userSlice) Swap(i, j int) { d[i], d[j] = d[j], d[i] }

type usersByDisplayName struct {
	userSlice
}

type usersByMail struct {
	userSlice
}

type usersByOnPremisesSamAccountName struct {
	userSlice
}

// Less reports whether the element with index i
// must sort before the element with index j.
func (u usersByDisplayName) Less(i, j int) bool {
	return strings.ToLower(u.userSlice[i].GetDisplayName()) > strings.ToLower(u.userSlice[j].GetDisplayName())
}

// Less reports whether the element with index i
// must sort before the element with index j.
func (u usersByMail) Less(i, j int) bool {
	return strings.ToLower(u.userSlice[i].GetMail()) > strings.ToLower(u.userSlice[j].GetMail())
}

// Less reports whether the element with index i
// must sort before the element with index j.
func (u usersByOnPremisesSamAccountName) Less(i, j int) bool {
	return strings.ToLower(u.userSlice[i].GetOnPremisesSamAccountName()) > strings.ToLower(u.userSlice[j].GetOnPremisesSamAccountName())
}

type groupSlice []*libregraph.Group

// Len is the number of elements in the collection.
func (d groupSlice) Len() int { return len(d) }

// Swap swaps the elements with indexes i and j.
func (d groupSlice) Swap(i, j int) { d[i], d[j] = d[j], d[i] }

type groupsByDisplayName struct {
	groupSlice
}

// Less reports whether the element with index i
// must sort before the element with index j.
func (g groupsByDisplayName) Less(i, j int) bool {
	return strings.ToLower(g.groupSlice[i].GetDisplayName()) > strings.ToLower(g.groupSlice[j].GetDisplayName())
}
