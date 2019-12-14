package script

import (
	"fmt"
	"sync"

	"go.starlark.net/starlark"
)

var threadPool = sync.Pool{
	// TODO: Update print function in threadpool to use provided shell stdout.
	New: func() interface{} {
		return &starlark.Thread{
			Name:  "executor",
			Print: func(_ *starlark.Thread, msg string) { fmt.Println(msg) },
		}
	},
}

func getThread() *starlark.Thread {
	return threadPool.Get().(*starlark.Thread)
}

func putThread(t *starlark.Thread) {
	threadPool.Put(t)
}
