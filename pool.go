package websocket

import (
	"sync"
)

// WorkerPool 通过pool模式提高worker的吞吐率
type WorkerPool struct {
	// worker数量
	workers int
	// 任务channel
	taskQueue chan Context
	// 并发锁
	mu sync.Mutex
	// 是否关闭
	closed bool
	// 处理器，也就是回调
	handler func(ctx Context)
	// 关闭channel
	done chan struct{}
}

// newWorkerPool return instance of WorkerPool
func newWorkerPool(workerNumber, taskNumber int) *WorkerPool {
	return &WorkerPool{
		workers: workerNumber,
		taskQueue: make(chan Context, taskNumber),
		done:      make(chan struct{}),
		handler:   func(ctx Context) {}, // default handler
	}
}

func (pool *WorkerPool) setHandler(handler func(ctx Context)) {
	if handler == nil {
		return
	}
	pool.handler = handler
}

// stop 停止所有协程的工作
func (pool *WorkerPool) stop() {
	// 加锁
	pool.mu.Lock()
	pool.closed = true
	close(pool.done)
	close(pool.taskQueue)
	pool.mu.Unlock()
}

// addTask 新工作
func (pool *WorkerPool) addTask(ctx Context) {
	pool.mu.Lock()
	// 不用defer解锁，因为在for循环里defer影响性能
	if pool.closed {
		pool.mu.Unlock()
		return
	}
	pool.mu.Unlock()
	pool.taskQueue <- ctx
}

// start 启动多协程
func (pool *WorkerPool) start() {
	// 通过协程去开启多个worker去处理connection
	for i := 0; i < pool.workers; i++ {
		go pool.work()
	}
}

// work 工作
func (pool *WorkerPool) work() {
	for {
		select {
		case <-pool.done: // 当workerPool.Close()后，关闭所有的worker
			return
		case ctx := <-pool.taskQueue: // 有新的connection进来后，分配给worker处理
			pool.handler(ctx)
		}
	}
}
