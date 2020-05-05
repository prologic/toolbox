package rc_worker_impl

import (
	"github.com/watermint/toolbox/essentials/log/es_log"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/recipe/rc_worker"
	"github.com/watermint/toolbox/infra/util/ut_runtime"
	"sync"
)

func NewQueue(ctl app_control.Control, concurrency int) rc_worker.Queue {
	q := &Queue{
		ctl: ctl,
		wg:  sync.WaitGroup{},
		q:   make(chan rc_worker.Worker),
	}
	q.Launch(concurrency)
	return q
}

type Queue struct {
	ctl         app_control.Control
	wg          sync.WaitGroup
	q           chan rc_worker.Worker
	numWorker   int
	maxWorker   int
	workerMutex sync.Mutex
	lastErr     error
}

func (z *Queue) dequeue() {
	l := z.ctl.Log().With(es_log.String("Routine", ut_runtime.GetGoRoutineName()))
	jobId := 0

	for w := range z.q {
		ll := l.With(es_log.Int("Seq", jobId))
		jobId++

		ll.Debug("Run work")
		if err := w.Exec(); err != nil {
			z.lastErr = err
			ll.Debug("Failure", es_log.Error(err))
		} else {
			ll.Debug("Success")
		}
		z.wg.Done()
	}
	l.Debug("Shutdown")
}

func (z *Queue) Launch(concurrency int) {
	l := z.ctl.Log().With(es_log.String("Routine", ut_runtime.GetGoRoutineName()))
	if concurrency < 1 {
		l.Debug("RunConcurrently must positive number, use 1 as default", es_log.Int("concurrency", concurrency))
		concurrency = 1
	}

	l.Debug("Launch first worker", es_log.Int("concurrency", concurrency))
	z.maxWorker = concurrency
	z.numWorker = 1
	go z.dequeue()
}

func (z *Queue) Enqueue(w rc_worker.Worker) {
	l := z.ctl.Log().With(es_log.String("Routine", ut_runtime.GetGoRoutineName()))

	z.workerMutex.Lock()
	defer z.workerMutex.Unlock()

	if z.numWorker < z.maxWorker {
		l.Debug("Launch additional worker", es_log.Int("numWorker", z.numWorker), es_log.Int("maxWorker", z.maxWorker))
		go z.dequeue()
		z.numWorker++
	}

	z.wg.Add(1)
	z.q <- w
}

func (z *Queue) Wait() error {
	l := z.ctl.Log().With(es_log.String("Routine", ut_runtime.GetGoRoutineName()))
	l.Debug("Waiting for worker shutdown")
	z.wg.Wait()
	close(z.q)

	return z.lastErr
}
