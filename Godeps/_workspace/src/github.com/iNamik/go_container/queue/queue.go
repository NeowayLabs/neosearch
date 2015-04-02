// Package queue implements an unbounded First-In-First-Out queue.
package queue

// queue.Interface
type Interface interface {

	// Add adds a new element
	Add(interface{})

	// Remove removes the next element and returns it
	Remove() interface{}

	// Peek allows you to look at elements without removing them
	Peek(int) interface{}

	// Clear empties all elements, returning the Len() to 0
	Clear()

	// Len returns the current length of the queue
	Len() int

	// Empty returns true if there are no items in the queue
	Empty() bool

	// Cap returns the current capacity of the internal store.
	// NOTE This is not the same as the Len(), but when Len() and Cap() are the
	// same, that means that the next call to Add() will have to allocate new
	// memory to increase the size of the internal storage.
	Cap() int

	// AtCapacity returns true if Len() == Cap()
	AtCapactiy() bool
}

// New
func New(capacity int) Interface {
	return &queue{data: make([]interface{}, capacity), head: 0, tail: 0, length: 0}
}
