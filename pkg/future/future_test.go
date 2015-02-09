package future

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func nilGetter() (interface{}, error) {
	return nil, nil
}

func intpGetter() (interface{}, error) {
	i := 1
	return &i, nil
}

func intGetter() (interface{}, error) {
	i := 2
	return i, nil
}

func TestNewFuture_nilType(t *testing.T) {
	assert := assert.New(t)
	f := New(nil, nilGetter)

	assert.NotNil(f)
	Wait(f)
	err := f.Err()
	assert.NotNil(err)
	assert.Equal(nilTypeError, err, err.Error())
}

func TestNewFuture_intpType(t *testing.T) {
	assert := assert.New(t)
	var itype int
	f := New(&itype, intpGetter)

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
	f := New(itype, intGetter)

	assert.NotNil(f)
	Wait(f)
	err := f.Err()
	assert.Nil(err)

	obj := f.Result()
	assert.IsType(itype, obj)
	assert.Equal(2, obj)
}

func TestNewFuture_interfaceTypeNilInCtor(t *testing.T) {
	assert := assert.New(t)
	var itype a
	f := New(itype, func() (interface{}, error) {
		return itype, nil
	})

	assert.NotNil(f)
	Wait(f)
	err := f.Err()
	assert.NotNil(err)
	assert.Equal(nilTypeError, err, err.Error())
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
