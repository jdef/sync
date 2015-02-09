package future

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelect_noFutures(t *testing.T) {
	assert := assert.New(t)
	f := Select()
	assert.NotNil(f)

	Wait(f)
	err := f.Err()
	assert.Nil(err)

	assert.True(IsReady(f))
	obj := f.Result()
	assert.NotNil(obj)

	f2, ok := obj.(Interface)
	assert.True(ok)
	assert.True(IsDiscarded(f2))
	assert.Nil(f2.Err())
}

func TestSelect_nilFuture(t *testing.T) {
	assert := assert.New(t)
	f := Select(nilFuture)
	assert.NotNil(f)

	Wait(f)
	err := f.Err()
	assert.Nil(err)

	assert.True(IsReady(f))
	obj := f.Result()
	assert.NotNil(obj)

	f2, ok := obj.(Interface)
	assert.True(ok)
	assert.True(IsDiscarded(f2))
	assert.Nil(f2.Err())
}
