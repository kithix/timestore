package timestore

import (
	"reflect"
	"testing"
	"time"
)

func timeFromYear(year int) time.Time {
	return time.Date(year, 0, 0, 0, 0, 0, 0, time.UTC)
}

func timeFromSeconds(seconds time.Duration) time.Time {
	return time.Unix(0, 0).Add(seconds)
}

func TestStoreAdd(t *testing.T) {
	s := BinaryDataTypeSamples{}
	s.Add(timeFromYear(1), "1")
	if !reflect.DeepEqual(s.data, []DataType{"1"}) {
		t.Error("Unable to add a single value", s)
	}

	s.Add(timeFromYear(3), "3")
	if !reflect.DeepEqual(s.data, []DataType{"1", "3"}) {
		t.Error("Did not append 3", s)
	}

	s.Add(timeFromYear(2), "2")
	if !reflect.DeepEqual(s.data, []DataType{"1", "2", "3"}) {
		t.Error("Did not insert 2 between 1 and 3", s)
	}
}

func TestStoreAllAfter(t *testing.T) {
	store := NewDataTypeBinary()
	store.Add(timeFromYear(2001), "2001")
	store.Add(timeFromYear(2010), "2010")
	store.Add(timeFromYear(2020), "2020")
	store.Add(timeFromYear(2025), "2025")
	store.Add(timeFromYear(2030), "2030")
	store.Add(timeFromYear(2031), "2031")
	store.Add(timeFromYear(2040), "2040")
	store.Add(timeFromYear(2050), "2050")

	testCases := []struct {
		in  int
		out interface{}
	}{
		{0, store.data},
		{2001, store.data},
		{2050, []DataType{"2050"}},
		{2051, []DataType{}},
	}

	for _, c := range testCases {
		value, time := store.AllAfter(timeFromYear(c.in))
		_ = time
		if !reflect.DeepEqual(value, c.out) {
			t.Errorf(
				"Expected year %v to retrieve all of %v, instead retrieved %v",
				c.in,
				c.out,
				value,
			)
		}
	}
}

func TestStoreAllBefore(t *testing.T) {
	store := NewDataTypeBinary()
	store.Add(timeFromYear(2001), "2001")
	store.Add(timeFromYear(2010), "2010")
	store.Add(timeFromYear(2020), "2020")
	store.Add(timeFromYear(2025), "2025")
	store.Add(timeFromYear(2030), "2030")
	store.Add(timeFromYear(2031), "2031")
	store.Add(timeFromYear(2040), "2040")
	store.Add(timeFromYear(2050), "2050")

	testCases := []struct {
		in  int
		out interface{}
	}{
		{0, []DataType{}},
		{2001, []DataType{"2001"}},
		{2002, []DataType{"2001"}},
		{2050, store.data},
		{2051, store.data},
	}
	for _, c := range testCases {
		value, time := store.AllBefore(timeFromYear(c.in))
		_ = time
		if !reflect.DeepEqual(value, c.out) {
			t.Errorf(
				"Expected year %v to retrieve all of %v, instead retrieved %v",
				c.in,
				c.out,
				value,
			)
		}
	}
}

func TestStoreClosestAfter(t *testing.T) {
	store := NewDataTypeBinary()
	store.Add(timeFromYear(2001), "2001")
	store.Add(timeFromYear(2010), "2010")
	store.Add(timeFromYear(2020), "2020")
	store.Add(timeFromYear(2025), "2025")
	store.Add(timeFromYear(2030), "2030")
	store.Add(timeFromYear(2031), "2031")
	store.Add(timeFromYear(2040), "2040")
	store.Add(timeFromYear(2050), "2050")

	testCases := []struct {
		in  int
		out interface{}
	}{
		{0, "2001"},
		{2001, "2001"},
		{2002, "2010"},
		{2010, "2010"},
		{2011, "2020"},
		{2040, "2040"},
		{2050, "2050"},
		{2051, nil},
		{3051, nil},
	}

	for _, c := range testCases {
		value, time := store.ClosestAfter(timeFromYear(c.in))
		_ = time
		if value != c.out {
			t.Errorf(
				"Expected year %v to retrieve year %v, instead retrieved %v",
				c.in,
				c.out,
				value,
			)
		}
	}
}

func TestStoreEmpty(t *testing.T) {
	store := NewDataTypeBinary()
	dd, tt := store.ClosestBefore(timeFromYear(1999))
	if dd != nil {
		t.Error("retrieved incorrect value: ", dd)
	}
	notime := time.Time{}
	if tt != notime {
		t.Error("Incorrect time retrieved: ", tt)
	}
	dd, tt = store.ClosestAfter(timeFromYear(1999))
	if dd != nil {
		t.Error("retrieved incorrect value: ", dd)
	}
	if tt != notime {
		t.Error("Incorrect time retrieved: ", tt)
	}
}
func TestStoreSingleValue(t *testing.T) {
	store := NewDataTypeBinary()
	store.Add(timeFromYear(2000), "2000")

	d, tt := store.ClosestBefore(timeFromYear(2001))
	if d != "2000" {
		t.Error("retrieved incorrect value: ", d)
	}
	if tt != timeFromYear(2000) {
		t.Error("Incorrect time retrieved: ", tt)
	}
	d, tt = store.ClosestBefore(timeFromYear(1999))
	if d != nil {
		t.Error("retrieved incorrect value: ", d)
	}
	notime := time.Time{}
	if tt != notime {
		t.Error("Incorrect time retrieved: ", tt)
	}
}

func TestStoreClosestBefore(t *testing.T) {
	store := NewDataTypeBinary()
	store.Add(timeFromYear(2001), "2001")
	store.Add(timeFromYear(2010), "2010")
	store.Add(timeFromYear(2020), "2020")
	store.Add(timeFromYear(2025), "2025")
	store.Add(timeFromYear(2030), "2030")
	store.Add(timeFromYear(2031), "2031")
	store.Add(timeFromYear(2040), "2040")
	store.Add(timeFromYear(2050), "2050")

	testCases := []struct {
		in  int
		out string
	}{
		{2001, "2001"},
		{2002, "2001"},
		{2010, "2010"},
		{2011, "2010"},
		{2040, "2040"},
		{2050, "2050"},
		{2051, "2050"},
		{3051, "2050"},
	}

	for _, c := range testCases {
		value, time := store.ClosestBefore(timeFromYear(c.in))
		_ = time
		if value != c.out {
			t.Errorf(
				"Expected year %v to retrieve year %v, instead retrieved %v",
				c.in,
				c.out,
				value,
			)
		}
	}
}

// Benchmarks

func makeStoreDataType(i int, b *testing.B) BinaryDataTypeSamples {
	s := BinaryDataTypeSamples{}
	for t := 0; t < i; t++ {
		var d DataType
		s.Add(timeFromYear(t), d)
	}
	b.ResetTimer()
	return s
}
func bStoreBefore(i int, b *testing.B) {
	t := timeFromYear(i / 2)
	s := makeStoreDataType(i, b)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		s.ClosestBefore(t)
	}
}
func bStoreAfter(i int, b *testing.B) {
	t := timeFromYear(i / 2)
	s := makeStoreDataType(i, b)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		s.ClosestAfter(t)
	}
}

func BenchmarkStoreClosestBefore100(b *testing.B)      { bStoreBefore(100, b) }
func BenchmarkStoreClosestBefore100000(b *testing.B)   { bStoreBefore(100000, b) }
func BenchmarkStoreClosestBefore10000000(b *testing.B) { bStoreBefore(10000000, b) }

func BenchmarkStoreClosestAfter100(b *testing.B)      { bStoreAfter(100, b) }
func BenchmarkStoreClosestAfter100000(b *testing.B)   { bStoreAfter(100000, b) }
func BenchmarkStoreClosestAfter10000000(b *testing.B) { bStoreAfter(10000000, b) }
