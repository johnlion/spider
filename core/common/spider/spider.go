package spider

import (
	"github.com/johnlion/spider/core/common/pageProcesser"
	"github.com/johnlion/spider/core/common/xlog"
	"github.com/johnlion/spider/core/common/scheduler"
	"github.com/johnlion/spider/core/common/downloader"
	"github.com/johnlion/spider/core/common/request"
	"github.com/johnlion/spider/core/common/pageItems"
	"github.com/johnlion/spider/core/common/pipeline"
	"github.com/johnlion/spider/core/common/resourceManage"
	"time"
	"os"
	"fmt"
	"math/rand"
	"github.com/johnlion/spider/core/common/page"
)

type Spider struct{
	tastname string
	threadnum uint
	exitWhenComplete bool
	pPageProcesser  pageProcesser.PageProcesser
	pScheduler scheduler.Scheduler
	pDownloader downloader.Downloader
	pPipelines []pipeline.Pipeline
	mc resourceManage.ResourceMange

	startSleeptime uint
	endSleeptime uint
	sleeptype string






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

func ( this *Spider ) Taskname() string{
	return this.tastname
}

//func ( this *Spider ) Get( url string, respType string ) *pageItems.PageItems{
//	req := request.NewRequest( url, respType, "", "GET", "", nil, nil, nil, nil )
//	return this.GetByRequest(req)
//}
//
func ( this *Spider ) GetByRequest( req *request.Request ) *pageItems.PageItems{
	var reqs []*request.Request
	reqs = append( reqs, req )
	items := this.GetAllByRequest( reqs )
	if len(items) != 0 {
		return items[0]
	}

	return nil

}

func ( this *Spider ) GetAllByRequest( reqs []*request.Request ) []*pageItems.PageItems{
	// push url
	for _, req := range reqs {
		//req := request.NewRequest(u, respType, urltag, method, postdata, header, cookies)
		this.AddRequest( req )
	}

	pip := pipeline.NewCollectPipelinePageItems()
	this.AddPipeline( pip )

	this.Run()

	fmt.Printf( "%v", pip.GetCollected() )
	os.Exit(1)
	return pip.GetCollected()

}

func ( this *Spider ) Run(){
	if this.threadnum == 0 {
		this.threadnum = 1
	}

	this.mc = resourceManage.NewResourceManageChan( this.threadnum )

	for {
		req := this.pScheduler.Poll()
		if this.mc.Has() == 0 && req == nil && this.exitWhenComplete{
			xlog.StraceInst().Println( "**** executed callback ****" )
			this.pPageProcesser.Finish()
			xlog.StraceInst().Println( "**** end spider ****" )
			break
		}else if req == nil {
			time.Sleep( 500 * time.Millisecond  )
			continue
		}
		this.mc.GetOne()

		go func ( req *request.Request ){
			defer this.mc.FreeOne()
			xlog.StraceInst().Println( "start crawl: " + req.GetUrl() )
			this.pageProcess( req )
		}( req )
	}
	this.close()
}

func ( this *Spider ) close(){
	this.SetScheduler( scheduler.NewSchedulerQUeue( false ) )
	this.SetDownloader( downloader.NewDownloaderHttp() )
	this.exitWhenComplete = true
}

func ( this *Spider ) AddPipeline( p pipeline.Pipeline ) *Spider{
	this.pPipelines = append( this.pPipelines, p )
	return this
}

func ( this *Spider ) AddRequest( req *request.Request ) *Spider{
	if req == nil{
		xlog.LogInst().LogError( "request is nil" )
	}else if req.GetUrl() == ""{
		xlog.LogInst().LogError( "request is empty" )
	}
	this.pScheduler.Push( req )
	return this
}

func ( this *Spider ) AddRequests ( reqs []*request.Request ) *Spider{
	for _, req := range reqs{
		this.AddRequest( req )
	}
	return this
}

//core processer
func ( this *Spider ) pageProcess( req *request.Request ){
	var p *page.Page

	defer func(){
		if err := recover(); err != nil{ // do not affect other
			if strerr , ok := err.(string); ok{
				xlog.LogInst().LogError( strerr )
			}else{
				xlog.LogInst().LogError( "pageProcess error" )
			}
		}
	}()

	//download page
	for i:=0; i<3; i++{
		this.sleep()
		p = this.pDownloader.Download(  req )
	}


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

func ( this *Spider ) sleep(){
	if  this.sleeptype == "fixed" {
		time.Sleep( time.Duration( this.startSleeptime )  * time.Millisecond )
	}else if this.sleeptype == "rand"{
		sleeptime := rand.Intn( int( this.endSleeptime-this.startSleeptime ) ) + int( this.startSleeptime )
		time.Sleep(  time.Duration(  sleeptime ) * time.Millisecond  )
	}
}