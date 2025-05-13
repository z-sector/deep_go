package main

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Group struct {
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	errOnce  sync.Once
	firstErr error
}

func NewErrGroup(ctx context.Context) (*Group, context.Context) {
	newCtx, cancel := context.WithCancel(ctx)
	return &Group{ctx: newCtx, cancel: cancel}, newCtx
}

func (g *Group) Go(action func() error) {
	if action == nil || g.ctx.Err() != nil {
		return
	}

	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		g.execAction(action)
	}()
}

func (g *Group) Wait() error {
	g.wg.Wait()
	g.cancel()
	return g.firstErr
}

func (g *Group) execAction(action func() error) {
	select {
	case <-g.ctx.Done():
		g.errOnce.Do(func() {
			g.firstErr = g.ctx.Err()
			g.cancel()
		})
	default:
		if err := action(); err != nil {
			g.errOnce.Do(func() {
				g.firstErr = err
				g.cancel()
			})
		}
	}
}

func TestErrGroupWithoutError(t *testing.T) {
	var counter atomic.Int32
	group, _ := NewErrGroup(context.Background())

	for i := 0; i < 5; i++ {
		group.Go(func() error {
			time.Sleep(time.Second)
			counter.Add(1)
			return nil
		})
	}

	err := group.Wait()
	assert.Equal(t, int32(5), counter.Load())
	assert.NoError(t, err)
}

func TestErrGroupWithError(t *testing.T) {
	var counter atomic.Int32
	group, ctx := NewErrGroup(context.Background())

	for i := 0; i < 5; i++ {
		group.Go(func() error {
			timer := time.NewTimer(time.Second)
			defer timer.Stop()

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-timer.C:
				counter.Add(1)
				return nil
			}
		})
	}

	group.Go(func() error {
		return errors.New("error")
	})

	err := group.Wait()
	assert.Equal(t, int32(0), counter.Load())
	assert.Error(t, err)
}
