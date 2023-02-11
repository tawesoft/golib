package reflect_test

import (
    "fmt"
    "io"

    "github.com/tawesoft/golib/v2/fun/slices"
    "github.com/tawesoft/golib/v2/must"
    "github.com/tawesoft/golib/v2/operator/reflect"
)

type Handle struct{
    closed bool
}

func (h *Handle) Close() error {
    if h.closed { return fmt.Errorf("handle already closed") }
    h.closed = true
    return nil
}

func OpenHandle() *Handle {
    return &Handle{}
}

func ExampleCast() {
    handles := []*Handle{
        OpenHandle(),
        OpenHandle(),
        OpenHandle(),
    }

    closeEverything := func(closers ... io.Closer) {
        for _, closer := range closers {
            must.Check(closer.Close())
        }
    }

    closeables := slices.Map(reflect.Cast[*Handle, io.Closer], handles)
    closeEverything(closeables...)

    checkClosed := func(handles ... *Handle) {
        for _, handle := range handles {
            must.True(handle.closed)
        }
    }

    handlesAgain := slices.Map(reflect.Cast[io.Closer, *Handle], closeables)
    checkClosed(handlesAgain...)

    // Output:
}
