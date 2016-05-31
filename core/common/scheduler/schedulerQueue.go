package scheduler

import (
	"sync"
	"container/list"
	"crypto/md5"
	"github.com/johnlion/spider/core/common/request"
)

type SchedulerQueue struct {
	locker *sync.Mutex
	rm bool
	queue *list.List
	rmKey map[[md5.Size]byte]*list.Element
}


// Interface
func ( this *SchedulerQueue ) Push( requ *request.Request){
	this.locker.Lock()
	var key [md5.Size]byte
	if this.rm{
		key = md5.Sum( []byte( requ.GetUrl() ) )
		if _, ok := this.rmKey[key]; ok{
			this.locker.Unlock()
			return
		}
	}
	e := this.queue.PushBack( requ )
	if this.rm{
		this.rmKey[key] = e
	}
	this.locker.Unlock()
}

// Interface
func ( this *SchedulerQueue ) Poll() *request.Request{
	this.locker.Lock()
	if this.queue.Len() <= 0 {
		this.locker.Unlock()
		return nil
	}

	e := this.queue.Front()
	requ := e.Value.( *request.Request )
	key := md5.Sum( []byte( requ.GetUrl() ) )
	this.queue.Remove( e )
	if this.rm {
		delete( this.rmKey, key )
	}
	this.locker.Unlock()
	return requ

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