package async

import (
	"context"
)

// Future interface has the method signature for await
type Future interface {
	// Await for result
	Await() any

	// AwaitWithContext await for result with context
	// AwaitWithContext(ctx context.Context) any
}

type future struct {
	await func(ctx context.Context) any
}

func (f future) Await() any {
	return f.await(context.Background())
}

// Exec executes the async function
func Exec(f func() any) Future {
	var result any
	c := make(chan any)
	go func() {
		defer close(c) // 执行完后，关闭 chan
		result = f()   // 执行业务函数，拿到数据
	}()

	return future{
		await: func(ctx context.Context) any {
			select {
			case <-ctx.Done(): // 超时
				return ctx.Err()
			case <-c:
				return result
			}
		},
	}
}
