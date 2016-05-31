package scheduler

import (
	"sync"
	"container/list"
	"crypto/md5"
)

type SchedulerQueue struct {
	locker *sync.Mutex
	rm bool
	queue *list.List
	rmKey map[[md5.Size]byte]*list.Element
}


// Interface
func ( this *SchedulerQueue ) Push(){
}

// Interface
func ( this *SchedulerQueue ) Poll(){

}

// Interface
func ( this *SchedulerQueue ) Count() int{
	this.locker.Lock()
	len := this.queue.Len()
	this.locker.Unlock()
	return len
}

func NewSchedulerQUeue ( rmDuplicate bool ) *SchedulerQueue {
	queue := list.New()
	rmKey := make( map[ [md5.Size]byte ]*list.Element )
	locker := new( sync.Mutex )
	return &SchedulerQueue{ rm: rmDuplicate, queue: queue, rmKey: rmKey, locker: locker  }
}