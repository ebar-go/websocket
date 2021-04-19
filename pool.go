package websocket

import (
	"log"
	"sync"
)


// workerPool 通过pool模式提高worker的吞吐率
type workerPool struct {
	id int
	// 选项
	option workerPoolOption
	// 任务channel
	taskQueue chan int
	// 并发锁
	mu sync.Mutex
	// 是否关闭
	closed bool
	// 工作内容，也就是回调
	job func(fd int)
	// 关闭channel
	done chan struct{}
}

func newWorkerPool(opts ...Option) *workerPool {
	// default option
	option := workerPoolOption{
		workers: 50,
		tasks:   100000,
	}
	for _, opt := range opts {
		opt.apply(&option)
	}
	log.Println("start worker pool:", option.tasks, option.workers)

	return &workerPool{
		option: option,
		taskQueue: make(chan int, option.tasks),
		done:      make(chan struct{}),
		job: func(fd int) {
			
		},
	}
}

// stop 停止所有协程的工作
func (p *workerPool) stop() {
	// 加锁
	p.mu.Lock()
	p.closed = true
	close(p.done)
	close(p.taskQueue)
	p.mu.Unlock()
}

// addTask 新工作
func (p *workerPool) addTask(fd int) {
	p.mu.Lock()
	// 不用defer解锁，因为在for循环里defer影响性能
	if p.closed {
		p.mu.Unlock()
		return
	}
	p.mu.Unlock()
	p.taskQueue <- fd
}

// start 启动多协程
func (p *workerPool) start() {
	// 通过协程去开启多个worker去处理connection
	for i := 0; i < p.option.workers; i++ {
		go p.work(i)
	}
}
// work
func (p *workerPool) work(id int) {
	for {
		select {
		case <-p.done: // 当workerPool.Close()后，关闭所有的worker
			return
		case fd := <- p.taskQueue: // 有新的connection进来后，分配给worker处理
			p.job(fd)
		}
	}
}
