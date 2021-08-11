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
	done   chan struct{}
	err    error
}

func Try(parent context.Context, action func(TryContext)) error {
	ctx := &errContext{parent: parent, done: make(chan struct{})}
	actionComplete := make(chan bool)
	go func() {
		select {
		case <-ctx.Done():
			actionComplete <- false
			runtime.Goexit()
			return
		default:
			action(ctx)
			actionComplete <- true
			return
		}
	}()
	select {
	case <-parent.Done():
		close(ctx.done)
		return fmt.Errorf("parent context error: %v", parent.Err())
	case <-actionComplete:
		return ctx.err
	}
}

func (ctx *errContext) Catch(err error) {
	if ctx.err != nil {
		return
	}
	if err != nil {
		ctx.err = err
		close(ctx.done)
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
