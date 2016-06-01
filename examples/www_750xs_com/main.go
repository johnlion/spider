package main
import(
	"github.com/johnlion/spider/core/common/xlog"


	"fmt"
	"github.com/johnlion/spider/core/common/request"
	"github.com/johnlion/spider/core/common/page"
	"github.com/johnlion/spider/core/common/downloader"
	"github.com/PuerkitoBio/goquery"
	"github.com/johnlion/spider/core/common/spider"
	"github.com/johnlion/spider/core/common/pipeline"
)

//
//
//
type MyPageProcesser struct {

}


// Interface
func NewMyPageProcesser () *MyPageProcesser{
	return &MyPageProcesser{}
}

// Interface
func ( this *MyPageProcesser ) Processer( p *page.Page ){

}

// Interface
func ( this *MyPageProcesser ) Finish(){
	fmt.Println( "TODO:before end spider \r\n" )
}




func main() {
	xlog.LogInst().Open()
	xlog.LogInst().LogInfo("Spider is running ... ")
	xlog.LogInst().Close()

	spider.NewSpider(NewMyPageProcesser(), "sina_stock_news").
	AddUrl("http://live.sina.com.cn/zt/api/l/get/finance/globalnews1/index.htm?format=json&id=63621&pagesize=10&dire=f", "json"). // start url, html is the responce type ("html" or "json" or "jsonp" or "text")
	AddPipeline(pipeline.NewPipelineConsole()).                                                                                   // Print result to std output
	AddPipeline(pipeline.NewPipelineFile("/tmp/sinafile")).                                                                       // Print result in file
	OpenFileLog("/tmp").                                                                                                          // Error info or other useful info in spider will be logged in file of defalt path like "WD/log/log.2014-9-1".
	SetSleepTime("rand", 1000, 3000).                                                                                             // Sleep time between 1s and 3s.
	Run()

}