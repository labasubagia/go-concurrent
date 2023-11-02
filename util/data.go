package util

import "time"

// GenNestedDuration
// generate dummy data for test
func GenNestedDuration(outerSize, innerSize int, duration time.Duration) (data [][]time.Duration) {
	for i := 0; i < outerSize; i++ {
		inner := []time.Duration{}
		for j := 0; j < innerSize; j++ {
			inner = append(inner, duration)
		}
		data = append(data, inner)
	}
	return data
}
