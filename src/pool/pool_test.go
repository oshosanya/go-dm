package pool

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRoutinePool(t *testing.T) {
	pool := NewRoutinePool(2)
	pool.Submit(jobOne)
	assert.Equal(t, 1, pool.LastJobID)
	pool.Submit(jobTwo)
	assert.Equal(t, 2, pool.LastJobID)
	pool.Submit(jobThree)
	pool.Submit(jobThree)
}

func jobOne() {
	for i := 0; i < 10; i++ {
		println("This is job one")
	}
}

func jobTwo() {
	for i := 0; i < 10; i++ {
		println("This is job two")
	}
}

func jobThree() {
	for i := 0; i < 10; i++ {
		println("This is job three")
	}
}
