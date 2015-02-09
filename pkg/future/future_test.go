package future

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFuture_nilType(t *testing.T) {
	assert := assert.New(t)
	f := New(nil, func() (interface{}, error) {
		return nil, nil
	})

	assert.NotNil(f)
	Wait(f)
	err := f.Err()
	assert.Nil(err)
	assert.True(IsReady(f))
	assert.Nil(f.Result())
}

func TestNewFuture_nilTypeNonNilResult(t *testing.T) {
	assert := assert.New(t)
	f := New(nil, func() (interface{}, error) {
		i := 1
		return i, nil
	})

	assert.NotNil(f)
	Wait(f)
	err := f.Err()
	assert.Nil(err)
	assert.Equal(1, f.Result())
}

func TestNewFuture_nilTypeZeroValResult(t *testing.T) {
	assert := assert.New(t)
	f := New(nil, func() (interface{}, error) {
		return &bar{}, nil
	})

	assert.NotNil(f)
	Wait(f)
	err := f.Err()
	assert.Nil(err)
	assert.Equal(&bar{}, f.Result())
}

func TestNewFuture_intpType(t *testing.T) {
	assert := assert.New(t)
	var itype int
	f := New(&itype, func() (interface{}, error) {
		i := 1
		return &i, nil
	})

	assert.NotNil(f)
	Wait(f)
	err := f.Err()
	assert.Nil(err)
	obj := f.Result()
	assert.IsType(&itype, obj)
	assert.Equal(1, *(obj.(*int)))
}

func TestNewFuture_intType(t *testing.T) {
	assert := assert.New(t)
	var itype int
	f := New(itype, func() (interface{}, error) {
		i := 2
		return i, nil
	})

	assert.NotNil(f)
	Wait(f)
	assert.True(IsDone(f))
	assert.False(HasError(f))

	err := f.Err()
	assert.Nil(err)
	assert.True(IsReady(f))
	assert.False(IsDiscarded(f))

	obj := f.Result()
	assert.IsType(itype, obj)
	assert.Equal(2, obj)

	f.Discard()
	assert.True(IsDone(f))
	assert.True(IsDiscarded(f))
	assert.False(IsReady(f))
}

func TestNewFuture_intTypeWithError(t *testing.T) {
	assert := assert.New(t)
	e := errors.New("expected error")
	var itype int
	f := New(itype, func() (x interface{}, err error) {
		err = e
		return
	})

	assert.NotNil(f)
	Wait(f)
	assert.True(IsDone(f))
	assert.True(HasError(f))

	err := f.Err()
	assert.NotNil(err)
	assert.Equal(e, err)
	assert.False(IsReady(f))
	assert.False(IsDiscarded(f))

	f.Discard()
	assert.True(IsDone(f))
	assert.True(IsDiscarded(f))
	assert.False(HasError(f))
}

func TestNewFuture_intTypeIllegalResultType(t *testing.T) {
	assert := assert.New(t)
	var itype int
	f := New(itype, func() (interface{}, error) {
		return "", nil
	})

	assert.NotNil(f)
	Wait(f)
	assert.True(IsDone(f))
	assert.True(HasError(f))

	err := f.Err()
	assert.NotNil(err)
	assert.IsType(nilTypeError, err, err.Error()) // type mismatch error
	assert.NotEqual(nilTypeError, err, err.Error())
	assert.False(IsReady(f))
	assert.False(IsDiscarded(f))

	f.Discard()
	assert.True(IsDone(f))
	assert.True(IsDiscarded(f))
	assert.False(HasError(f))
}

func TestNewFuture_interfaceTypeNilSpecType(t *testing.T) {
	assert := assert.New(t)
	var itype a
	f := New(itype, func() (interface{}, error) {
		return itype, nil
	})

	assert.NotNil(f)
	Wait(f)
	err := f.Err()
	assert.Nil(err)
	assert.Nil(f.Result())
}

func TestNewFuture_interfaceTypeNilInGetter(t *testing.T) {
	assert := assert.New(t)
	b := &bar{}
	f := New(b, func() (interface{}, error) {
		return nil, nil
	})

	assert.NotNil(f)
	Wait(f)
	err := f.Err()
	assert.Nil(err)
	obj := f.Result()
	assert.Nil(obj)
}

func TestNewFuture_interfaceType(t *testing.T) {
	assert := assert.New(t)
	b := &bar{}
	f := New(b, func() (interface{}, error) {
		c := &bar{}
		return c, nil
	})

	assert.NotNil(f)
	Wait(f)
	err := f.Err()
	assert.Nil(err)

	obj := f.Result()
	assert.IsType(b, obj)
	assert.Equal(b, obj)
}

type a interface {
	foo()
}

type bar struct{ zoom int }

func (b *bar) foo() {}

// this is pretty terrible, using pointers to interfaces. but if you have
// to do it, here's how
func TestNewFuture_interfacepType(t *testing.T) {
	assert := assert.New(t)

	var itype *a
	b := &bar{}
	c := &bar{}
	d := &bar{1}

	f := New(itype, func() (interface{}, error) {
		cast := a(b)
		return &cast, nil
	})

	assert.NotNil(f)
	Wait(f)
	err := f.Err()
	assert.Nil(err)

	obj := f.Result()
	assert.IsType(itype, obj)
	cast := a(b)
	assert.Equal(&cast, obj)
	cast = a(c)
	assert.Equal(&cast, obj)
	cast = a(d)
	assert.IsType(&cast, obj)
	assert.NotEqual(&cast, obj)
}
