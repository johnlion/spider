package pipeline

import (
	"github.com/johnlion/spider/core/common/pageItems"
	"github.com/johnlion/spider/core/common/comInterfaces"
)

type CollectPipelinePageItems struct {
	collector []*pageItems.PageItems
}

func NewCollectPipelinePageItems() *CollectPipelinePageItems{
	collector := make ( []*pageItems.PageItems, 0 )
	return &CollectPipelinePageItems{ collector: collector }
}

func ( this *CollectPipelinePageItems ) Process( items *pageItems.PageItems, t comInterfaces.Task ){
	this.collector = append( this.collector, items )
}

func ( this *CollectPipelinePageItems ) GetCollected() []*pageItems.PageItems{
	return this.collector
}






