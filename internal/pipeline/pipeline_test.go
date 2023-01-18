package pipeline

import (
	"context"
	"github.com/tj/assert"
	"io"
	"sync"
	"testing"
	"time"
)

func TestOpen(t *testing.T) {
	p := Open(Settings{Id: "p"})
	p1 := Open(Settings{Id: "p"})
	assert.Equal(t, p, p1, "Must be same pipe")
	np := Open(Settings{Id: "np"})
	assert.NotEqual(t, p, np, "There are different pipes")
	Stop()
}

func TestPipeline_Do(t *testing.T) {
	p := Open(Settings{Id: "pipe"})
	result, err := p.Do(context.Background(), func() Result {
		return Result{Done: true, Result: 133}
	})
	assert.Equal(t, 133, result, "Result must be equal")
	assert.Nil(t, err, "No errors")

	result, err = p.Do(context.Background(), func() Result {
		return Result{Done: true, Err: io.EOF}
	})
	assert.Nil(t, result)
	assert.EqualError(t, err, io.EOF.Error())

	cnt := 0
	result, err = p.Do(context.Background(), func() Result {
		cnt++
		return Result{Done: false, Err: io.EOF}
	})
	assert.Equal(t, MaxAttempts, cnt, "Must repeated while max attempts counter will be reached")
	assert.Nil(t, result)
	assert.EqualError(t, err, io.EOF.Error())

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		_, err := p.Do(ctx, func() Result {
			select {
			case <-time.After(100 * time.Second):
				return Result{Done: true}
			case <-ctx.Done():
				return Result{Done: false}
			}
		})
		assert.EqualError(t, err, context.Canceled.Error())
	}()
	<-time.After(20 * time.Millisecond)
	cancel()

	s := map[int]struct{}{}
	wg := sync.WaitGroup{}
	const addTimes = 15000
	wg.Add(addTimes)
	for i := 0; i < addTimes; i++ {
		go func(i int) {
			defer wg.Done()
			_, _ = p.Do(context.Background(), func() Result {
				s[i] = struct{}{}
				return Result{Done: true}
			})
		}(i)
	}
	wg.Wait()

	assert.Equal(t, addTimes, len(s), "All values must be stored")
	for i := 0; i < addTimes; i++ {
		_, ok := s[i]
		assert.True(t, ok, "Value '%d' must be in the map", i)
	}
}

func TestStop(t *testing.T) {
	p := Open(Settings{Id: "pipe"})
	Stop()
	assert.Panics(t, func() {
		_, _ = p.Do(context.Background(), func() Result {
			return Result{}
		})
	})
}
