package future

// convenience func that returns true if f.Done() is closed and f.Err() == nil
func IsReady(f Interface) bool {
	select {
	case <-f.Done():
		return f.Err() == nil && !IsDiscarded(f)
	default:
		return false
	}
}

// convenience func that returns true if f.Done() is closed and f.Err() != nil
func HasError(f Interface) bool {
	select {
	case <-f.Done():
		return f.Err() != nil && !IsDiscarded(f)
	default:
		return false
	}
}

// convenience func that returns true if f.Done() is closed
func IsDone(f Interface) bool {
	select {
	case <-f.Done():
		return true
	default:
		return false
	}
}

// convenience func that returns true if f.Discarded() is closed
func IsDiscarded(f Interface) bool {
	select {
	case <-f.Discarded():
		return true
	default:
		return false
	}
}

// installs a callback handler to be invoked once the done channel has been closed
func OnAny(r Interface, f func(Interface)) {
	go func() {
		Wait(r)
		f(r)
	}()
}

// installs a callback handler to be invoked if Err() returns non-nil, but only once the done channel has been closed
func OnError(r Interface, f func(error)) {
	go func() {
		select {
		case <-r.Discarded():
			// noop
		case <-r.Done():
			if err := r.Err(); err != nil && !IsDiscarded(r) {
				f(err)
			}
		}
	}()
}

// installs a callback handler to be invoked if Err() returns nil, but only once the done channel has been closed
func OnReady(r Interface, f func(interface{})) {
	go func() {
		select {
		case <-r.Discarded():
			// noop
		case <-r.Done():
			if err := r.Err(); err == nil && !IsDiscarded(r) {
				f(r.Result())
			}
		}
	}()
}

// installs a callback handler to be invoked once the discard channel has been closed
func OnDiscard(r Interface, f func()) {
	go func() {
		select {
		case <-r.Discarded():
			f()
		}
	}()
}

// wait for the future to complete or to be discarded
func Wait(r Interface) {
	select {
	case <-r.Done():
	case <-r.Discarded():
	}
}
