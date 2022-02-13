package syncx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReentrantMutex(t *testing.T) {
	var rmu ReentrantMutex

	result := 0

	fn1 := func() {
		rmu.Lock()
		rmu.Unlock()

		result = 100
	}

	fn2 := func() {
		rmu.Lock()
		rmu.Unlock()

		fn1()
	}

	fn3 := func() {
		rmu.Lock()
		rmu.Unlock()

		fn2()
	}

	fn3()

	assert.Equal(t, 100, result)
}
