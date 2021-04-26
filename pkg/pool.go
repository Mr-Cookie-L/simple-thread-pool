package pkg

import "sync"

type Pool struct {
	capacity int
	running  int
	// TODO
	worker []string
	lock   sync.Mutex
}

func NewPool(size int) *Pool {
	return &Pool{
		capacity: size,
	}
}

func (p *Pool) Cap() int {
	return p.capacity
}
