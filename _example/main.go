package main

import (
	"fmt"
	"time"

	"github.com/kithix/timestore"
)

func main() {
	start := time.Now()
	granularity := 2 * time.Microsecond
	until := 2 * time.Second

	fixed := timestore.NewDataTypeFixed(start, granularity, until)
	flexible := timestore.NewDataTypeBinary()

	times := make([]time.Time, 0)
	for {
		now := time.Now()
		if now.After(start.Add(until)) {
			break
		}
		times = append(times, now)
		fixed.Add(now, now.String())
		flexible.Add(now, now.String())
	}

	fmt.Println("Times Added")
	fmt.Println(len(times))
	fmt.Println("Fixed")
	fmt.Println(fixed.Len())
	fmt.Println("Flex")
	fmt.Println(flexible.Len())
}
