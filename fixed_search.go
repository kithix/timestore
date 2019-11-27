package timestore

import "time"

func workoutSamplePosition(granularity time.Duration, since, t time.Time) int {
	return int(t.Sub(since) / granularity)
}

func makeTimes(since time.Time, granularity time.Duration, amountOf int) []time.Time {
	times := make([]time.Time, amountOf)
	i := 0
	for i < amountOf {
		times[i] = since.Add(granularity * time.Duration(i))
		i++
	}
	return times
}
