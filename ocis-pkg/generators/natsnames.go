package generators

import (
	"os"
	"strconv"
)

// NType is an enum type for the different types of NATS connections
type NType int

// Enum values for NType
const (
	NTypeBus NType = iota
	NTypeKeyValue
	NTypeRegistry
)

// String returns the string representation of a NType
func (n NType) String() string {
	return []string{"bus", "kv", "reg"}[n]
}

// GenerateConnectionName generates a connection name for a NATS connection
// The connection name will be formatted as follows: "hostname:pid:service:type"
func GenerateConnectionName(service string, ntype NType) string {
	host, err := os.Hostname()
	if err != nil {
		host = ""
	}

	return firstNRunes(host, 5) + ":" + strconv.Itoa(os.Getpid()) + ":" + service + ":" + ntype.String()
}

// firstNRunes returns the first n runes of a string
func firstNRunes(s string, n int) string {
	i := 0
	for j := range s {
		if i == n {
			return s[:j]
		}
		i++
	}
	return s
}
