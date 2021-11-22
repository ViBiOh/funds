package model

import (
	"context"
	"sync"
)

// Group describes a task group with a fail-fast approach
type Group struct {
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	once    sync.Once
	done    chan struct{}
	limiter chan struct{}
	err     error
}

// NewGroup creates a Group with given concurrency limit
func NewGroup(limit uint64) *Group {
	return &Group{
		done:    make(chan struct{}),
		limiter: make(chan struct{}, limit),
	}
}

// WithContext make the given context cancelable for the group
func (g *Group) WithContext(ctx context.Context) context.Context {
	if g.cancel != nil {
		panic("cancelable context already set-up")
	}

	ctx, g.cancel = context.WithCancel(ctx)
	return ctx
}

// Go run given function in a goroutine according to limiter and current status
func (g *Group) Go(f func() error) {
	select {
	case <-g.done:
	case g.limiter <- struct{}{}:
		g.wg.Add(1)

		go func() {
			defer g.wg.Done()
			if g.limiter != nil {
				defer func() { <-g.limiter }()
			}

			if err := f(); err != nil {
				g.close(err)
			}
		}()
	}
}

// Wait for Group to end
func (g *Group) Wait() error {
	defer g.close(nil)

	g.wg.Wait()
	return g.err
}

func (g *Group) close(err error) {
	g.once.Do(func() {
		close(g.done)

		if g.cancel != nil {
			g.cancel()
		}

		g.err = err
	})
}
