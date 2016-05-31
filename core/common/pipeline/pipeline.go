package pipeline

import (
	"github.com/johnlion/spider/core/common/pageItems"
	"github.com/johnlion/spider/core/common/comInterfaces"
)

type Pipeline interface {
	// The Process implements result persistent.
	// The items has the result be crawled.
	// The t has informations of this crawl task.
	Process(items *pageItems.PageItems, t comInterfaces.Task)
}

// The interface CollectPipeline recommend result in process's memory temporarily.
type CollectPipeline interface {
	Pipeline

	// The GetCollected returns result saved in in process's memory temporarily.
	GetCollected() []*pageItems.PageItems
}


