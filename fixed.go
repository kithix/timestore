package timestore

import "time"

//go:generate genny -in=$GOFILE -out=generated-$GOFILE gen "DataType=BUILTINS"

type FixedDataTypeSamples struct {
	since       time.Time
	until       time.Time
	granularity time.Duration
	earliestPos *int
	latestPos   *int
	data        []DataType
	times       []time.Time
}

func (s *FixedDataTypeSamples) Len() int {
	return len(s.data)
}

func (s *FixedDataTypeSamples) All() ([]DataType, []time.Time) {
	return s.data, s.times
}

func (s *FixedDataTypeSamples) Add(t time.Time, data DataType) {
	pos := s.position(t)
	if s.earliestPos == nil || pos < *s.earliestPos {
		s.earliestPos = &pos
	}
	if s.latestPos == nil || pos > *s.latestPos {
		s.latestPos = &pos
	}
	s.data[pos] = data
}

func (s *FixedDataTypeSamples) Earliest() (DataType, time.Time) {
	if s.earliestPos == nil {
		return s.data[0], time.Time{}
	}
	return s.data[*s.earliestPos], s.times[*s.earliestPos]
}

func (s *FixedDataTypeSamples) Latest() (DataType, time.Time) {
	if s.latestPos == nil {
		return s.data[len(s.data)-1], time.Time{}
	}
	return s.data[*s.latestPos], s.times[*s.latestPos]
}

func (s *FixedDataTypeSamples) AllBefore(t time.Time) ([]DataType, []time.Time) {
	pos := s.position(t)
	if pos < 0 {
		return []DataType{}, []time.Time{}
	}
	return s.data[:pos], s.times[:pos]
}

func (s *FixedDataTypeSamples) AllAfter(t time.Time) ([]DataType, []time.Time) {
	pos := s.position(t)
	if pos < 0 {
		return []DataType{}, []time.Time{}
	}
	return s.data[pos:], s.times[pos:]
}

func (s *FixedDataTypeSamples) ClosestBefore(t time.Time) (DataType, time.Time) {
	var d DataType
	pos := s.position(t)
	if pos < 0 {
		return d, time.Time{}
	}
	return s.data[pos], s.times[pos]
}

func (s *FixedDataTypeSamples) ClosestAfter(t time.Time) (DataType, time.Time) {
	var d DataType
	pos := s.position(t) + 1
	if pos < 0 {
		return d, time.Time{}
	}
	return s.data[pos], s.times[pos]
}

func (s *FixedDataTypeSamples) position(t time.Time) int {
	pos := workoutSamplePosition(s.granularity, s.since, t)
	if pos > len(s.data) {
		return -1
	}
	return pos
}

func NewDataTypeFixed(since time.Time, granularity, keepFor time.Duration) *FixedDataTypeSamples {
	until := since.Add(keepFor)
	amountOf := workoutSamplePosition(granularity, since, until)
	return &FixedDataTypeSamples{
		since,
		until,
		granularity,
		nil,
		nil,
		make([]DataType, amountOf),
		makeTimes(since, granularity, amountOf),
	}
}
