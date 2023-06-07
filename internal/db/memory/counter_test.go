package memDB

import (
	"fmt"
	"sync"
	"testing"
)

func TestCounter(t *testing.T) {
	c := newCounter(0)
	c.increment()

	if c.read() != 1 {
		t.Fatalf("Unexpected counter value after initial increment, expected %d, got %d", 1, c.read())
	}

	routines := 15
	expected := routines + c.read()

	var wg sync.WaitGroup
	wg.Add(routines)

	for i := 1; i <= routines; i++ {
		go func(i int) {
			defer wg.Done()
			val := c.increment()
			fmt.Printf("Routine %d set counter to %d\n", i, val)
		}(i)

	}

	wg.Wait()

	if c.read() != expected {
		t.Fatalf("Unexpected counter value after concurrent increments, expected %d, got %d", expected, c.read())
	}

}
