package promise

import (
	"errors"
	"sync"
	"sync/atomic"

	"github.com/jdef/sync/pkg/future"
)

var (
	DiscardedError = errors.New("promise discarded")
)

// interface common to all promises
type Interface interface {
	Discard()
	Future() future.Interface
	Get() (interface{}, error)
	Set(interface{}) bool
}

type Promise struct {
	ch          chan interface{}
	getter      sync.Once
	setter      sync.Once
	cached      interface{}
	discarded   int32
	closeMutex  sync.Mutex
	discardOnce sync.Once
	future      future.Interface
}

func New(ptrToType interface{}) Interface {
	b := &Promise{
		ch: make(chan interface{}, 1),
	}
	b.future = future.New(ptrToType, b.Get)
	future.OnDiscard(b.future, func() {
		// b.Discard will also trigger a discard of the future,
		// but that turns into a noop since the future is already
		// discarded.
		b.Discard()
	})
	return b
}

func (b *Promise) Future() future.Interface {
	return b.future
}

func (b *Promise) Get() (val interface{}, err error) {
	b.getter.Do(func() {
		select {
		// if we're discarded, the chan will already be closed, otherwise
		// we'll block until we either get a value from Set(), or else we're
		// Discard()ed.
		case v, ok := <-b.ch:
			if atomic.LoadInt32(&b.discarded) != 0 {
				err = DiscardedError
			} else if ok {
				b.cached = v
			} else {
				//programming error..
				panic("chan closed but promise not yet discarded")
			}
		}
	})
	val = b.cached
	return
}

func (b *Promise) Set(val interface{}) (success bool) {
	b.setter.Do(func() {
		// don't let someone change the channel while we're working
		b.closeMutex.Lock()
		defer b.closeMutex.Unlock()

		// write to the chan, if it's still open
		select {
		case <-b.ch:
			// this Set func is the only one that ever writes to the chan.
			// if it's already closed (discarded) this becomes a noop.
		default:
			b.ch <- val
			success = true
		}
	})
	return
}

func (b *Promise) Discard() {
	b.discardOnce.Do(func() {
		b.closeMutex.Lock()
		defer b.closeMutex.Unlock()
		atomic.StoreInt32(&b.discarded, 1)
		close(b.ch)

		go func() {
			select {
			case <-b.future.Discarded():
				//noop
			default:
				b.future.Discard()
			}
		}()
	})
}
