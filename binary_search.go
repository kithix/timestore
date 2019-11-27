package timestore

import (
	"sort"
	"time"
)

func indexBefore(times []time.Time, t time.Time) int {
	if len(times) == 0 {
		return -1
	}
	// If the latest time is before or equal to this time, max length
	if times[len(times)-1].Before(t) ||
		times[len(times)-1].Equal(t) {
		return len(times) - 1
	}
	// Earliest time is after this time, we can't get a time before
	if times[0].After(t) {
		return -1
	}
	// If our earliest time is equal
	if times[0].Equal(t) {
		return 0
	}
	i := sort.Search(len(times)-1, func(i int) bool {
		return times[i].After(t)
	})

	return i - 1
}

func indexAfter(times []time.Time, t time.Time) int {
	if len(times) == 0 {
		return -1
	}
	if times[len(times)-1].Before(t) {
		return -1
	}
	if times[0].After(t) || times[0].Equal(t) {
		return 0
	}
	return sort.Search(len(times)-1, func(i int) bool {
		return times[i].After(t) || times[i].Equal(t)
	})
}
