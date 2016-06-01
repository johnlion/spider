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


// Deal with one url and return the PageItems.
func (this *Spider) Get(url string, respType string) *pageItems.PageItems {
	req := request.NewRequest(url, respType, "", "GET", "", nil, nil, nil, nil)
	return this.GetByRequest(req)
}

// Deal with several urls and return the PageItems slice.
func (this *Spider) GetAll(urls []string, respType string) []*pageItems.PageItems {
	for _, u := range urls {
		req := request.NewRequest(u, respType, "", "GET", "", nil, nil, nil, nil)
		this.AddRequest(req)
	}

	pip := pipeline.NewCollectPipelinePageItems()
	this.AddPipeline(pip)

	this.Run()

	return pip.GetCollected()
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
		xlog.StraceInst().Println( " sleep pageProcess ..." )
		p = this.pDownloader.Download(  req )
	}

	if !p.IsSucc() { // if fail do not need process
		return
	}

	this.pPageProcesser.Processer( p )
	for _, req := range p.GetTargetRequests() {
		this.AddRequest(req)
	}

	// output
	if !p.GetSkip() {
		for _, pip := range this.pPipelines {
			fmt.Println("%v",p.GetPageItems().GetAll())
			pip.Process(p.GetPageItems(), this)
		}
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

func (this *Spider) SetThreadnum(i uint) *Spider {
	this.threadnum = i
	return this
}

func (this *Spider) GetThreadnum() uint {
	return this.threadnum
}

// If exit when each crawl task is done.
// If you want to keep spider in memory all the time and add url from outside, you can set it true.
func (this *Spider) SetExitWhenComplete(e bool) *Spider {
	this.exitWhenComplete = e
	return this
}

func (this *Spider) GetExitWhenComplete() bool {
	return this.exitWhenComplete
}

// The OpenFileLog initialize the log path and open log.
// If log is opened, error info or other useful info in spider will be logged in file of the filepath.
// Log command is mlog.LogInst().LogError("info") or mlog.LogInst().LogInfo("info").
// Spider's default log is closed.
// The filepath is absolute path.
func (this *Spider) OpenFileLog(filePath string) *Spider {
	xlog.InitFileLog(true, filePath)
	return this
}

// OpenFileLogDefault open file log with default file path like "WD/log/log.2014-9-1".
func (this *Spider) OpenFileLogDefault() *Spider {
	xlog.InitFileLog(true, "")

	return this
}

// The CloseFileLog close file log.
func (this *Spider) CloseFileLog() *Spider {
	xlog.InitFileLog(false, "")
	return this
}

// The OpenStrace open strace that output progress info on the screen.
// Spider's default strace is opened.
func (this *Spider) OpenStrace() *Spider {
	xlog.StraceInst().Open()
	return this
}

// The CloseStrace close strace.
func (this *Spider) CloseStrace() *Spider {
	xlog.StraceInst().Close()
	return this
}

// The SetSleepTime set sleep time after each crawl task.
// The unit is millisecond.
// If sleeptype is "fixed", the s is the sleep time and e is useless.
// If sleeptype is "rand", the sleep time is rand between s and e.
func (this *Spider) SetSleepTime(sleeptype string, s uint, e uint) *Spider {
	this.sleeptype = sleeptype
	this.startSleeptime = s
	this.endSleeptime = e
	if this.sleeptype == "rand" && this.startSleeptime >= this.endSleeptime {
		panic("startSleeptime must smaller than endSleeptime")
	}

}

func ( this *Spider ) sleep(){
	if  this.sleeptype == "fixed" {
		time.Sleep( time.Duration( this.startSleeptime )  * time.Millisecond )
	}else if this.sleeptype == "rand"{
		sleeptime := rand.Intn( int( this.endSleeptime-this.startSleeptime ) ) + int( this.startSleeptime )
		time.Sleep(  time.Duration(  sleeptime ) * time.Millisecond  )
	}
}

func (this *Spider) AddUrl(url string, respType string) *Spider {
	req := request.NewRequest(url, respType, "", "GET", "", nil, nil, nil, nil)
	this.AddRequest(req)
	return this
}

func (this *Spider) AddUrlEx(url string, respType string, headerFile string, proxyHost string) *Spider {
	req := request.NewRequest(url, respType, "", "GET", "", nil, nil, nil, nil)
	this.AddRequest(req.AddHeaderFile(headerFile).AddProxyHost(proxyHost))
	return this
}

func (this *Spider) AddUrlWithHeaderFile(url string, respType string, headerFile string) *Spider {
	req := request.NewRequestWithHeaderFile(url, respType, headerFile)
	this.AddRequest(req)
	return this
}

func (this *Spider) AddUrls(urls []string, respType string) *Spider {
	for _, url := range urls {
		req := request.NewRequest(url, respType, "", "GET", "", nil, nil, nil, nil)
		this.AddRequest(req)
	}
	return this
}

func (this *Spider) AddUrlsWithHeaderFile(urls []string, respType string, headerFile string) *Spider {
	for _, url := range urls {
		req := request.NewRequestWithHeaderFile(url, respType, headerFile)
		this.AddRequest(req)
	}
	return this
}

func (this *Spider) AddUrlsEx(urls []string, respType string, headerFile string, proxyHost string) *Spider {
	for _, url := range urls {
		req := request.NewRequest(url, respType, "", "GET", "", nil, nil, nil, nil)
		this.AddRequest(req.AddHeaderFile(headerFile).AddProxyHost(proxyHost))
	}
	return this
}

