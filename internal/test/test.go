package test

import (
    "testing"
    "time"
)

// Completes executes f (in a goroutine), and blocks until either f returns,
// or the provided duration has elapsed. In the latter case, calls t.Errorf to
// fail the test. Provide optional format string and arguments to add
// context to the test error message.
func Completes(t *testing.T, duration time.Duration, f func(), args ... interface{}) {
    done := make(chan struct{}, 1)
    timeout := time.After(duration)
    go func() {
        f()
        done <- struct{}{}
    }()

    select {
        case <-done: // OK
        case <-timeout:
            if len(args) > 0 {
                t.Errorf("test timed out after "+duration.String()+": " + args[0].(string), args[1:]...)
            } else {
                t.Errorf("test timed out after %s", duration.String())
            }
    }
}
