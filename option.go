package websocket

type workerPoolOption struct {
	// worker数量
	workers int
	// 最大任务数量,也就是预估最大连接数
	tasks int
}

type Option interface {
	apply(opt *workerPoolOption)
}

type workerNumberOption int

func (o workerNumberOption) apply(opt *workerPoolOption) {
	opt.workers = int(o)
}

// 设置worker数量
func WithWorkerNumber(n int) Option {
	return workerNumberOption(n)
}

type taskNumberOption int

func (o taskNumberOption) apply(opt *workerPoolOption) {
	opt.tasks = int(o)
}

// 设置task数量
func WithTaskNumber(n int) Option {
	return taskNumberOption(n)
}