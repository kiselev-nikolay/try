package try_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/kiselev-nikolay/try"
)

func TestTry(t *testing.T) {
	t.Run("no action", func(t *testing.T) {
		err := try.Try(context.Background(), func(tc try.TryContext) {})
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("nil err", func(t *testing.T) {
		err := try.Try(context.Background(), func(tc try.TryContext) {
			tc.Catch(nil)
			if tc.Err() != nil {
				t.Fail()
			}
		})
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("cancel parent", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			<-time.After(1 * time.Second)
			cancel()
		}()
		happen := false
		err := try.Try(ctx, func(tc try.TryContext) {
			<-time.After(2 * time.Second)
			tc.Catch(nil)
			happen = true
		})
		if err.Error() != "parent context error: context canceled" {
			t.Error(err)
		}
		if happen {
			t.Error("must not happen")
		}
	})
	t.Run("err", func(t *testing.T) {
		err := try.Try(context.Background(), func(tc try.TryContext) {
			var v interface{}
			err := json.Unmarshal([]byte(""), &v)
			tc.Catch(err)
		})
		if err == nil {
			t.Error("error must be not nil")
		}
	})
	t.Run("deadline", func(t *testing.T) {
		deadline := time.Now().Add(time.Hour)
		ctx, cancel := context.WithDeadline(context.Background(), deadline)
		defer cancel()
		err := try.Try(ctx, func(tc try.TryContext) {
			d, ok := tc.Deadline()
			if !ok || d != deadline {
				tc.Catch(fmt.Errorf("deadline value is unexpected"))
			}
		})
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("value", func(t *testing.T) {
		const key = "test"
		ctx := context.WithValue(context.Background(), key, "test")
		err := try.Try(ctx, func(tc try.TryContext) {
			v := tc.Value(key)
			if v != "test" {
				tc.Catch(fmt.Errorf("context value is unexpected"))
			}
		})
		if err != nil {
			t.Error(err)
		}
	})
}
