package try

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

type TryContext interface {
	Catch(error)
	context.Context
}

type errContext struct {
	parent context.Context
	closed bool
	done   chan struct{}
	err    error
}

func Try(parent context.Context, action func(TryContext)) error {
	ctx := &errContext{
		parent: parent,
		closed: false,
		done:   make(chan struct{}),
	}
	go func() {
		action(ctx)
		if !ctx.closed {
			close(ctx.done)
			ctx.closed = true
		}
	}()
	select {
	case <-parent.Done():
		if !ctx.closed {
			close(ctx.done)
			ctx.closed = true
		}
		ctx.err = fmt.Errorf("parent context error: %v", parent.Err())
		return ctx.err
	case <-ctx.done:
		return ctx.err
	}
}

func (ctx *errContext) Catch(err error) {
	if ctx.err != nil {
		runtime.Goexit()
		return
	}
	if err != nil {
		ctx.err = err
		if !ctx.closed {
			close(ctx.done)
			ctx.closed = true
		}
		runtime.Goexit()
	}
}

func (ctx *errContext) Deadline() (deadline time.Time, ok bool) {
	return ctx.parent.Deadline()
}

func (ctx *errContext) Done() <-chan struct{} {
	return ctx.done
}

func (ctx *errContext) Err() error {
	if err := ctx.parent.Err(); err != nil {
		return err
	}
	if ctx.err != nil {
		return ctx.err
	}
	return nil
}

func (ctx *errContext) Value(key interface{}) interface{} {
	return ctx.parent.Value(key)
}
