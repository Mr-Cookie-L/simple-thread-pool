package pkg

type Worker struct {
	pool *Pool
	task chan func()
}

func NewWorker(pool *Pool) *Worker {
	return &Worker{
		pool: pool,
	}
}

func (w *Worker) run() {
	go func() {
		for f := range w.task {
			if f == nil {
				// TODO 退出
			}
			f()
			// TODO 回收worker
		}
	}()
}
