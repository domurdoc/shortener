package utils

import "sync"

// Resetter is an interface for objects that can reset their state to a clean, reusable condition.
// This is used in conjunction with Pool to allow safe reuse of objects without reallocation.
type Resetter interface {
	Reset()
}

// Pool[T] is a type-safe wrapper around sync.Pool for managing reusable objects of type T.
// It ensures that objects are properly reset before being returned to the pool,
// making it safe to reuse mutable objects in concurrent environments.
// The type parameter T must implement the Resetter interface.
type Pool[T Resetter] struct {
	sync.Pool
}

func NewPool[T Resetter](newF func() T) *Pool[T] {
	return &Pool[T]{
		Pool: sync.Pool{New: func() any { return newF() }},
	}
}

func (p *Pool[T]) Get() T {
	return p.Pool.Get().(T)
}

func (p *Pool[T]) Put(x T) {
	x.Reset()
	p.Pool.Put(x)
}
