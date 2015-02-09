package future

import (
	"sync"
)

// Returns a future that captures any Ready future in a set. Note that
// Do() DOES NOT capture a future that has failed or been discarded. If
// no such ready future will ever exist then a discarded future is returned.
func Select(args ...Interface) (fof Interface) {
	return New(nilFuture, func() (interface{}, error) {
		var wg sync.WaitGroup
		wg.Add(len(args))
		ch := make(chan Interface, len(args))
		for _, f := range args {
			OnAny(f, func(f Interface) {
				defer wg.Done()
				if IsReady(f) {
					ch <- f
				}
			})
		}
		wg.Wait()
		select {
		case f := <-ch:
			return f, nil
		default:
			// there will never be a *ready* future
			return nilFuture, nil
		}
	})
}
