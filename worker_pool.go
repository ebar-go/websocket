package websocket

import (
	"sync"
)

type pool struct {
	// worker数量
	workers int
	// 最大任务数量,也就是预估最大连接数
	maxTasks int
	// 任务channel
	taskQueue chan Connection
	// 并发锁
	mu sync.Mutex
	// 是否关闭
	closed bool
	// 关闭channel
	done chan struct{}
}

func newPool(w int, t int) *pool {
	return &pool{
		workers:   w,
		maxTasks:  t,
		taskQueue: make(chan Connection, t),
		done:      make(chan struct{}),
	}
}

func (p *pool) Close() {
	// 加锁
	p.mu.Lock()
	p.closed = true
	close(p.done)
	close(p.taskQueue)
	p.mu.Unlock()
}

func (p *pool) addTask(conn Connection) {
	p.mu.Lock()
	// 不用defer解锁，因为在for循环里defer影响性能
	if p.closed {
		p.mu.Unlock()
		return
	}
	p.mu.Unlock()
	p.taskQueue <- conn
}

// start 启动多协程
func (p *pool) start(engine *Engine) {
	// 通过协程去开启多个worker去处理connection
	for i := 0; i < p.workers; i++ {
		go p.work(engine)
	}
}
// work
func (p *pool) work(engine *Engine) {
	for {
		select {
		case <-p.done: // 当pool.Close()后，关闭所有的worker
			return
		case conn := <- p.taskQueue: // 有新的connection进来后，分配给worker处理
			if conn != nil {
				ctx, err := conn.context()
				if err != nil {
					break
				}
				// process conn
				engine.handle(ctx)
			}
		}
	}
}
