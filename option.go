package websocket

type options struct {
	// worker数量
	workers int
	// 最大任务数量,也就是预估最大连接数
	tasks int
}

// Option server option interface
type Option interface {
	apply(opt *options)
}

// workerNumberOption worker数量选项
type workerNumberOption int

// apply implements of Option
func (o workerNumberOption) apply(opt *options) {
	opt.workers = int(o)
}

// WithWorkerNumber 设置worker数量
func WithWorkerNumber(n int) Option {
	return workerNumberOption(n)
}

// taskNumberOption task数量选项
type taskNumberOption int

// apply implements of Option
func (o taskNumberOption) apply(opt *options) {
	opt.tasks = int(o)
}

// WithTaskNumber 设置task数量
func WithTaskNumber(n int) Option {
	return taskNumberOption(n)
}
