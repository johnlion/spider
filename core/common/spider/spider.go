package spider

import (
	"github.com/johnlion/spider/core/common/pageProcesser"
	"github.com/johnlion/spider/core/common/xlog"
	"github.com/johnlion/spider/core/common/scheduler"
	"github.com/johnlion/spider/core/common/downloader"
	//"github.com/johnlion/spider/core/common/pageItems"
	"github.com/johnlion/spider/core/common/request"
	//"github.com/johnlion/spider/core/common/pipeline"
	"github.com/johnlion/spider/core/common/pageItems"
	"fmt"
)

type Spider struct{
	tastname string
	threadnum int
	exitWhenComplete bool
	pPageProcesser  pageProcesser.PageProcesser
	pScheduler scheduler.Scheduler
	pDownloader downloader.Downloader



}

func NewSpider( taskname string , pageinst pageProcesser.PageProcesser  ) *Spider{
	xlog.LogInst().Open()
	app := &Spider{ tastname: taskname , pPageProcesser: pageinst }
	// init fileLog
	xlog.LogInst().LogInfo( "tastname " + taskname + " is processing ... " )
	xlog.StraceInst().Println( "[tastname] " + taskname + " is processing ... " )
	xlog.LogInst().Close()

	//init spider
	if app.pScheduler == nil{
		app.SetScheduler( scheduler.NewSchedulerQUeue( false ) )
	}

	if app.pDownloader == nil {
		app.SetDownloader( downloader.NewDownloaderHttp() )
	}


	xlog.StraceInst().Println( "**** start spider **** ")
	return app
}

//func ( this *Spider ) Get( url string, respType string ) *pageItems.PageItems{
//	req := request.NewRequest( url, respType, "", "GET", "", nil, nil, nil, nil )
//	return this.GetByRequest(req)
//}
//
func ( this *Spider ) GetByRequest( req *request.Request ) *pageItems.PageItems{
	var reqs []*request.Request
	reqs = append( reqs, req )
	fmt.Printf( "%s\n", reqs )
	//items := this.GetAllByRequest( reqs )
	//if len(items) != 0 {
	//	return items[0]
	//}
	return nil

}

//func ( this *Spider ) GetAllByRequest( reqs []*request.Request ) []*pageItems.PageItems{
//	// push url
//	for _, req := range reqs {
//		//req := request.NewRequest(u, respType, urltag, method, postdata, header, cookies)
//		this.AddRequest( req )
//	}
//
//	pip := pipeline.NewCollectPipelinePageItems()
//	this.AddPipeline(pip)
//
//	this.Run()
//
//	return pip.GetCollected()
//
//}

func ( this *Spider ) AddRequest( req *request.Request ) *Spider{
	if req == nil{
		xlog.LogInst().LogError( "request is nil" )
	}else if req.GetUrl() == ""{
		xlog.LogInst().LogError( "request is empty" )
	}
	this.pScheduler.Push()
	return this
}

func ( this *Spider ) AddRequests ( reqs []*request.Request ) *Spider{
	for _, req := range reqs{
		this.AddRequest( req )
	}
	return this
}

func ( this *Spider ) SetScheduler( s scheduler.Scheduler ) *Spider{
	this.pScheduler = s
	return this
}

func ( this *Spider ) GetScheduler() scheduler.Scheduler {
	return this.pScheduler
}

func ( this *Spider ) SetDownloader( d downloader.Downloader ) *Spider{
	this.pDownloader = d
	return this
}

func ( this *Spider ) GetDownloader() downloader.Downloader {
	return this.pDownloader
}

//func ( this *Spider ) AddUrl( url string, respType string ) *Spider{
//	req := request.NewRequest(url, respType, "", "GET", "", nil, nil, nil, nil)
//	this.AddRequest(req.AddHeaderFile(headerFile).AddProxyHost(proxyHost))
//return this
//}

