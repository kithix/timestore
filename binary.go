package timestore

import (
	"sort"
	"time"
)

//go:generate genny -in=$GOFILE -out=generated-$GOFILE gen "DataType=BUILTINS"

type BinaryDataTypeSamples struct {
	times []time.Time
	data  []DataType
}

func (s *BinaryDataTypeSamples) Len() int {
	return len(s.data)
}

func (s *BinaryDataTypeSamples) All() ([]DataType, []time.Time) {
	return s.data, s.times
}

func (s *BinaryDataTypeSamples) Add(t time.Time, data DataType) {
	// If this is our first data point or the new latest, add to the end
	if len(s.times) == 0 || s.times[len(s.times)-1].Before(t) {
		s.data = append(s.data, data)
		s.times = append(s.times, t)
		return
	}

	// If this is before our earliest sample, add it to the start
	if s.times[0].After(t) {
		s.data = append([]DataType{data}, s.data...)
		s.times = append([]time.Time{t}, s.times...)
		return
	}

	// Binary search
	i := sort.Search(len(s.times)-1, func(i int) bool {
		return s.times[i].After(t)
	})
	s.data = append(s.data[:i], append([]DataType{data}, s.data[i:]...)...)
	s.times = append(s.times[:i], append([]time.Time{t}, s.times[i:]...)...)
	return
}

func (s *BinaryDataTypeSamples) ClosestAfter(t time.Time) (DataType, time.Time) {
	// Nil value for unable to find
	var v DataType
	i := indexAfter(s.times, t)
	if i < 0 {
		return v, time.Time{}
	}
	return s.data[i], s.times[i]
}

func (s *BinaryDataTypeSamples) ClosestBefore(t time.Time) (DataType, time.Time) {
	// Nil value for unable to find
	var v DataType
	i := indexBefore(s.times, t)
	if i < 0 {
		return v, time.Time{}
	}
	return s.data[i], s.times[i]
}

func (s *BinaryDataTypeSamples) AllAfter(t time.Time) ([]DataType, []time.Time) {
	i := indexAfter(s.times, t)
	if i < 0 {
		return []DataType{}, []time.Time{}
	}
	return s.data[i:], s.times[i:]
}

func (s *BinaryDataTypeSamples) AllBefore(t time.Time) ([]DataType, []time.Time) {
	i := indexBefore(s.times, t)
	if i < 0 {
		return []DataType{}, []time.Time{}
	}
	return s.data[:i+1], s.times[:i+1]
}

func NewDataTypeBinary() *BinaryDataTypeSamples {
	return &BinaryDataTypeSamples{
		times: []time.Time{},
		data:  []DataType{},
	}
}
