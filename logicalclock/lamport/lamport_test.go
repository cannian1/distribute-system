package lamport

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogicalClockBasicOps(t *testing.T) {
	t.Run("初始化", func(t *testing.T) {
		var time int64 = 20
		c := NewLogicalClock(time)
		assert.Equal(t, time, c.Now())
	})

	t.Run("Increment", func(t *testing.T) {
		c := NewLogicalClock(0)
		var expected int64 = 1
		assert.Equal(t, expected, c.Increment())
		assert.Equal(t, expected, c.Now())
	})

	t.Run("Update", func(t *testing.T) {
		testCases := []struct {
			local    int64
			received int64
			expected int64
		}{
			{5, 3, 6}, // 本地时间较新
			{3, 5, 6}, // 接收时间较新
			{0, 0, 1}, // 相等情况
		}

		for _, tc := range testCases {
			c := NewLogicalClock(tc.local)
			result := c.Update(tc.received)
			if result != tc.expected || c.Now() != tc.expected {
				t.Errorf("Update(%d from %d): expected %d, got %d",
					tc.received, tc.local, tc.expected, result)
			}
		}
	})

	t.Run("Advance", func(t *testing.T) {
		t.Run("Valid", func(t *testing.T) {
			c := NewLogicalClock(10)
			var expected int64 = 15
			c.Advance(expected)
			assert.Equal(t, expected, c.Now())
		})

		t.Run("Invalid ", func(t *testing.T) {
			c := NewLogicalClock(15)
			var expected int64 = 15
			c.Advance(5)
			assert.Equal(t, expected, c.Now())
		})
	})
}

func TestConcurrentAccess(t *testing.T) {
	c := NewLogicalClock(0)
	const goroutines = 100
	var wg sync.WaitGroup

	for i := 0; i < goroutines; i++ {
		wg.Add(2)
		go func() {
			c.Increment()
			wg.Done()
		}()
		go func() {
			c.Now()
			wg.Done()
		}()
	}

	wg.Wait()

	assert.Equal(t, int64(goroutines), c.Now())
}
