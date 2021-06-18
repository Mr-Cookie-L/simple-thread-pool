package threadpool

type ThreadPool struct {
	maxProcs int

	goCh    chan struct{}
	taskQue chan struct {
		name string
		f    func() error
	}
}

// maxProcs sets the maximum number of goroutine
// queBuf sets the buffer size of Task channel
func New(maxProcs, queBuf int) *ThreadPool {
	return &ThreadPool{
		maxProcs: maxProcs,
		goCh:     make(chan struct{}, maxProcs),
		taskQue: make(chan struct {
			name string
			f    func() error
		}, queBuf),
	}
}

func (tp *ThreadPool) Run() {
	for i := 0; i < tp.maxProcs; i++ {
		tp.worker()
	}

	go func() {
		for {
			<-tp.goCh
			tp.worker()
		}
	}()
}

func (tp *ThreadPool) worker() {
	defer func() {
		if err := recover(); err != nil {
			// TODO 处理日志
		}
		tp.goCh <- struct{}{}
	}()

	for {
		task := <-tp.taskQue
		err := task.f()
		if err != nil {
			// TODO 处理日志
		}
	}
}

func (tp *ThreadPool) Push(name string, f func() error) {
	tp.taskQue <- struct {
		name string
		f    func() error
	}{name: name, f: f}
}
