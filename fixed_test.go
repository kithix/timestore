package timestore

import (
	"testing"
	"time"
)

// benchmark this
func TestWorkoutSamplePosition(t *testing.T) {
	testCases := []struct {
		granularity time.Duration
		since, time time.Time
		expect      int
	}{
		{500 * time.Millisecond, time.Unix(0, 0), time.Unix(2, 0), 4},
		{500 * time.Millisecond, time.Unix(0, 0), time.Unix(2, int64(499*time.Millisecond)), 4},
		{500 * time.Millisecond, time.Unix(0, 0), time.Unix(2, int64(500*time.Millisecond)), 5},
		{500 * time.Millisecond, time.Unix(0, 0), time.Unix(2, int64(1*time.Second)), 6},
		{1 * time.Second, time.Unix(0, 0), time.Unix(2, 0), 2},
		{2 * time.Second, time.Unix(0, 0), time.Unix(2, 0), 1},
		{1 * time.Second, time.Unix(5, 0), time.Unix(5, 0), 0},
		{8 * time.Second, time.Unix(5, 0), time.Unix(4804, 0), 599},
		{8 * time.Second, time.Unix(5, 0), time.Unix(4809, 0), 600},
	}

	for _, c := range testCases {
		samplePosition := workoutSamplePosition(
			c.granularity,
			c.since,
			c.time,
		)
		if samplePosition != c.expect {
			t.Error("Expected sample position to be", c.expect, ", got", samplePosition)
		}
	}
}

func shouldPanic(fn func(), t *testing.T) (paniced bool) {
	defer func() {
		r := recover()
		if r != nil {
			paniced = true
		}
	}()
	fn()
	return
}

func TestEmptyFixedStore(t *testing.T) {
	store := NewDataTypeFixed(
		time.Unix(1, 0), 1*time.Second, 10*time.Second,
	)
	dd, tt := store.ClosestBefore(time.Unix(1, 0))
	if dd != nil {
		t.Error("retrieved incorrect value: ", dd)
	}
	if tt != time.Unix(1, 0) {
		t.Error("Incorrect time retrieved: ", tt)
	}
	dd, tt = store.ClosestAfter(time.Unix(0, 0))
	if dd != nil {
		t.Error("retrieved incorrect value: ", dd)
	}
	if tt != time.Unix(1, 0) {
		t.Error("Incorrect time retrieved: ", tt)
	}
}

func TestFixedStore(t *testing.T) {
	// With range 1-10
	store := NewDataTypeFixed(
		time.Unix(1, 0), 1*time.Second, 10*time.Second,
	)
	store.Add(time.Unix(1, 0), "1")
	dd, tt := store.Earliest()
	if dd != "1" || tt != time.Unix(1, 0) {
		t.Error("Did not set earliest to 1")
	}
	dd, tt = store.Latest()
	if dd != "1" || tt != time.Unix(1, 0) {
		t.Error("Did not set latest to 1")
	}
	dd, tt = store.ClosestBefore(time.Unix(1, 0))
	if dd != "1" || tt != time.Unix(1, 0) {
		t.Error("Closest before did not retrieve correctly")
	}
	dd, tt = store.ClosestAfter(time.Unix(0, 0))
	if dd != "1" || tt != time.Unix(1, 0) {
		t.Error("Closest after did not retrieve correctly")
	}
	// Adding before a stores time is not doable (error? currently panics)
	if !shouldPanic(func() { store.Add(time.Unix(0, 0), "fails adding before") }, t) {
		t.Error("Expected adding earlier than a store to panic")
	}
	// Adding after a stores time is not doable (error?)
	if !shouldPanic(func() { store.Add(time.Unix(11, 0), "fails adding after") }, t) {
		t.Error("Expected adding after than a store to panic")
	}
	// Adding a close enough time replaces the previous
	store.Add(timeFromSeconds(1), "not1")
	dd, tt = store.ClosestBefore(time.Unix(1, 0))
	if dd != "not1" || tt != time.Unix(1, 0) {
		t.Error("Did not replace value as expected")
	}
}

func makeFixedStoreDataType(i int, b *testing.B) *FixedDataTypeSamples {
	s := NewDataTypeFixed(
		time.Unix(0, 0),
		1*time.Second,
		1*time.Second*time.Duration(i),
	)
	for t := 0; t < i; t++ {
		var d DataType
		s.Add(time.Unix(int64(t), 0), d)
	}
	b.ResetTimer()
	return s
}

func bFixedAfter(i int, b *testing.B) {
	t := time.Unix(int64(i)/2, 0)
	s := makeFixedStoreDataType(i, b)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		s.ClosestAfter(t)
	}
}
func bFixedBefore(i int, b *testing.B) {
	t := time.Unix(int64(i)/2, 0)
	s := makeFixedStoreDataType(i, b)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		s.ClosestBefore(t)
	}
}

// TODO benchmark sample position

func BenchmarkFixedClosestAfter100(b *testing.B)       { bFixedAfter(100, b) }
func BenchmarkFixedClosestAfter100000(b *testing.B)    { bFixedAfter(100000, b) }
func BenchmarkFixedClosestAfter10000000(b *testing.B)  { bFixedAfter(10000000, b) }
func BenchmarkFixedClosestBefore100(b *testing.B)      { bFixedBefore(100, b) }
func BenchmarkFixedClosestBefore100000(b *testing.B)   { bFixedBefore(100000, b) }
func BenchmarkFixedClosestBefore10000000(b *testing.B) { bFixedBefore(10000000, b) }
