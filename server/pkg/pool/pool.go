package pool

const (
	JobsBufferSize = 128
)

type Job func()

type Pool struct {
	Jobs chan Job
}

func NewPool(count int) *Pool {
	p := &Pool{
		Jobs: make(chan Job, JobsBufferSize),
	}

	for _ = range count {
		go p.worker()
	}

	return p
}

func (p *Pool) worker() {
	for job := range p.Jobs {
		job()
	}
}
