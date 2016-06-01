package pipeline

import (
	"github.com/johnlion/spider/core/common/pageItems"
	"github.com/johnlion/spider/core/common/comInterfaces"
)

type PipelineConsole struct {
}

func NewPipelineConsole() *PipelineConsole {
	return &PipelineConsole{}
}

func (this *PipelineConsole) Process(items *pageItems.PageItems, t comInterfaces.Task) {
	println("----------------------------------------------------------------------------------------------")
	println("Crawled url :\t" + items.GetRequest().GetUrl() + "\n")
	println("Crawled result : ")
	for key, value := range items.GetAll() {
		println(key + "\t:\t" + value)
	}
}
