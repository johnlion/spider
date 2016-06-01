package downloader_test

import (
	"testing"
	"github.com/johnlion/spider/core/common/request"
	"github.com/johnlion/spider/core/common/downloader"
	"github.com/johnlion/spider/core/common/page"
	"github.com/PuerkitoBio/goquery"
	"fmt"
)

func TestDownloadHtml( t *testing.T ){
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
		t.Error("html parse failed!")
	}else{
		fmt.Printf( "%v", s )

	}
}


func TestDownloadJson(t *testing.T) {
	//return
	var req *request.Request
	req = request.NewRequest("http://live.sina.com.cn/zt/api/l/get/finance/globalnews1/index.htm?format=json&id=23521&pagesize=4&dire=f&dpc=1", "json", "", "GET", "", nil, nil, nil, nil)

	var dl downloader.Downloader
	dl = downloader.NewDownloaderHttp()

	var p *page.Page
	p = dl.Download(req)

	var jsonMap interface{}
	jsonMap = p.GetJson()
	fmt.Printf("%v", jsonMap)

	//fmt.Println(doc)
	//body := p.GetBodyStr()
	//fmt.Println(body)

}

func TestCharSetChange(t *testing.T) {
	var req *request.Request
	//req = request.NewRequest("http://stock.finance.sina.com.cn/usstock/api/jsonp.php/t/US_CategoryService.getList?page=1&num=60", "jsonp")
	req = request.NewRequest("http://soft.chinabyte.com/416/13164916.shtml", "html", "", "GET", "", nil, nil, nil, nil)

	var dl downloader.Downloader
	dl = downloader.NewDownloaderHttp()

	var p *page.Page
	p = dl.Download(req)

	//hp := p.GetHtmlParser()
	//fmt.Printf("%v", jsonMap)

	//fmt.Println(doc)
	p.GetBodyStr()
	body := p.GetBodyStr()
	fmt.Println(body)

}
