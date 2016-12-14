package main

import "time"

// Item represents an expirable item associated with a given value.
type Item struct {
	value      interface{}
	expiration time.Time
}

// NewItem creates an item with the specified value and expiring on the
// specified time.
func NewItem(value interface{}, expiration time.Time) *Item {
	return &Item{
		value:      value,
		expiration: expiration,
	}
}

// NewItemWithTTL creates an item with the specified value and expiring after
// the specified duration.
func NewItemWithTTL(value interface{}, duration time.Duration) *Item {
	return NewItem(value, time.Now().Add(duration))
}

// Value returns the value stored in the item.
func (item *Item) Value() interface{} {
	return item.value
}

// Expiration returns the item's expiration time.
func (item *Item) Expiration() time.Time {
	return item.expiration
}

// Expired checks whether the item is already expired.
func (item *Item) Expired() bool {
	return item.expiration.Before(time.Now())
}

// TTL returns the remaining duration until expiration (negative if expired).
func (item *Item) TTL() time.Duration {
	return item.expiration.Sub(time.Now())
}
