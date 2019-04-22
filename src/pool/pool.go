package pool

import "sync"

type JobFunc func()

type Job struct {
	id       int
	executor JobFunc
	status   string
}

type routinePool struct {
	PoolSize  int
	JobSize   int
	QueueSize int
	LastJobID int
	Jobs      chan Job
}

var mutex = &sync.Mutex{}

func NewRoutinePool(size int) routinePool {
	newPool := routinePool{
		PoolSize:  size,
		QueueSize: 0,
		LastJobID: 0,
	}
	newPool.Jobs = make(chan Job)
	go newPool.LaunchJob()
	return newPool
}

func (p *routinePool) Submit(jobFunc JobFunc) {
	job := Job{id: p.LastJobID + 1, executor: jobFunc, status: "ready"}
	p.queueSizeInc()
	p.Jobs <- job
	p.LastJobID = p.LastJobID + 1
}

func (p *routinePool) LaunchJob() {
	for {
		println("Queue size ", p.QueueSize)
		if p.QueueSize < p.PoolSize {
			select {
			case jobToLaunch := <-p.Jobs:
				go func() {
					jobToLaunch.executor()
					p.queueSizeDec()
				}()
			default:

			}
		} else {
			println("Pool doesn't have any available workers, queue size is ", p.QueueSize)
		}
	}
}

func (p *routinePool) queueSizeInc() {
	mutex.Lock()
	p.QueueSize = p.QueueSize + 1
	mutex.Unlock()
}

func (p *routinePool) queueSizeDec() {
	mutex.Lock()
	p.QueueSize = p.QueueSize - 1
	mutex.Unlock()
}
