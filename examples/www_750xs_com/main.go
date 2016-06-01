package main
import(
	"github.com/johnlion/spider/core/common/xlog"


	"fmt"
	"github.com/johnlion/spider/core/common/request"
	"github.com/johnlion/spider/core/common/page"
	"github.com/johnlion/spider/core/common/downloader"
	"github.com/PuerkitoBio/goquery"
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


	var req *request.Request
	req = request.NewRequest("http://live.sina.com.cn/zt/l/v/finance/globalnews1/", "html", "", "GET", "", nil, nil, nil, nil)

	var dl downloader.Downloader
	dl = downloader.NewDownloaderHttp()

	var p *page.Page
	p = dl.Download( req )

	var doc *goquery.Document
	doc = p.GetHtmlParser()
	//fmt.Println(doc)
	//body := p.GetBodyStr()
	//fmt.Println(body)

	var s *goquery.Selection
	s = doc.Find("body")
	if s.Length() < 1 {
		xlog.StraceInst().Println("html parse failed!")
	}else{
		fmt.Printf( "%v", s )

	}
}