package websocket

import (
	"github.com/ebar-go/websocket/epoll"
	"log"
	"net/http"
	"sync"
)
// worker pool
type workerPool struct {
	// worker数量
	workers int
	// 最大任务数量,也就是预估最大连接数
	maxTasks int
	// 任务channel
	taskQueue chan int
	// 并发锁
	mu sync.Mutex
	// 是否关闭
	closed bool
	callback func(fd int)
	// 关闭channel
	done chan struct{}
}

func newPool(w int, t int) *workerPool {
	return &workerPool{
		workers:   w,
		maxTasks:  t,
		taskQueue: make(chan int, t),
		done:      make(chan struct{}),
	}
}

func (p *workerPool) Close() {
	// 加锁
	p.mu.Lock()
	p.closed = true
	close(p.done)
	close(p.taskQueue)
	p.mu.Unlock()
}

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
	for i := 0; i < p.workers; i++ {
		go p.work()
	}
}
// work
func (p *workerPool) work() {
	for {
		select {
		case <-p.done: // 当workerPool.Close()后，关闭所有的worker
			return
		case fd := <- p.taskQueue: // 有新的connection进来后，分配给worker处理
			p.callback(fd)
		}
	}
}



type workerPoolServer struct {
	server
	// epoll
	epoller epoll.Epoll

	// worker pool
	workerPool *workerPool
}

func (srv *workerPoolServer) HandleRequest(w http.ResponseWriter, r *http.Request) {
	// 获取socket连接
	conn, err := newConnection(w, r)
	if err != nil {
		// do something..
		return
	}
	if err := srv.epoller.Add(conn.fd()); err != nil {
		log.Println("unable to add connection:", err.Error())
		_ = conn.close()
		return
	}
	srv.AddConnection(conn)
}

// Close connection
func (srv *workerPoolServer) Close(conn Connection) {
	if err := srv.epoller.Remove(conn.fd()); err != nil {
		log.Println("unable to remove conn:", err.Error())
	}
	srv.RemoveConnection(conn)
}

func (srv *workerPoolServer) Start() {
	srv.workerPool.start()
	srv.workerPool.callback = func(fd int) {
		conn , exist := srv.GetConnection(fd)
		if !exist {
			return
		}
		ctx, err := conn.context()
		if err != nil {
			srv.Close(conn)
			return
		}
		srv.engine.handle(ctx)
	}

	go func() {
		defer srv.workerPool.Close()
		for {
			fds, err := srv.epoller.Wait()
			if err != nil {
				log.Println("unable to get active socket connection from epoll:", err)
				continue
			}
			for _, fd := range fds {
				srv.workerPool.addTask(fd)
			}
		}

	}()
}

// WorkerPoolServer 多进程server，相比epoll的单进程，降低了延迟
//
// workers 进程数
//
// maxConnections 最大并发连接数
func WorkerPoolServer(workers , maxConnections int) Server {
	e, err := epoll.Create()
	if err != nil {
		log.Fatalf("create epoll:%v\n", err)
	}
	return &workerPoolServer{
		server: base(),
		epoller: e ,
		workerPool: newPool(workers, maxConnections),
	}
}
