package scheduler

type Scheduler interface {
	Push()
	Poll()
	Count() int
}