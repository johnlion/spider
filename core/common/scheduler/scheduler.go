package scheduler

import "github.com/johnlion/spider/core/common/request"

type Scheduler interface {
	Push( requ *request.Request )
	Poll() *request.Request
	Count() int
}