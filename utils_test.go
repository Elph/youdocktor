package main

import (
	"testing"
	"time"

	assert "github.com/stretchr/testify/assert"
)

func TestRangeDate(t *testing.T) {
	tests := []struct {
		start time.Time
		end   time.Time

		expectedCount int
	}{
		{time.Now(), time.Now().AddDate(0, 0, 3), 4},
		{time.Now(), time.Now(), 1},
		{time.Now(), time.Now().AddDate(0, 0, 60), 61},
	}

	for _, theory := range tests {

		// act
		count := 0
		for rd := rangeDate(theory.start, theory.end); ; {
			date := rd()
			if date.IsZero() {
				break
			}
			count = count + 1
		}

		assert.Equal(t, theory.expectedCount, count)
	}
}
