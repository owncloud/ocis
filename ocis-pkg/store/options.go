package store

import "time"

// Option provides an option to configure the store
type Option func(*Options)

// Type defines the type of the store
func Type(typ string) Option {
	return func(o *Options) {
		o.Type = typ
	}
}

// Addresses defines the addresses where the store can be reached
func Addresses(addrs ...string) Option {
	return func(o *Options) {
		o.Addresses = addrs
	}
}

// Database defines the Database the store should use
func Database(db string) Option {
	return func(o *Options) {
		o.Database = db
	}
}

// Table defines the table the store should use
func Table(t string) Option {
	return func(o *Options) {
		o.Table = t
	}
}

// Size defines the maximum capacity of the store.
// Only applicable when using "ocmem" store
func Size(s int) Option {
	return func(o *Options) {
		o.Size = s
	}
}

// TTL defines the time to life for elements in the store.
// Only applicable when using "natsjs" store
func TTL(t time.Duration) Option {
	return func(o *Options) {
		o.TTL = t
	}
}
