package downloader

import (
	"github.com/johnlion/spider/core/common/request"
	"github.com/johnlion/spider/core/common/page"
	"github.com/johnlion/spider/core/common/xlog"
	"github.com/PuerkitoBio/goquery"
	"github.com/bitly/go-simplejson"
	"net/http"
	"net/url"
	"strings"

	"io"
	"compress/gzip"
	"golang.org/x/net/html/charset"
	"io/ioutil"
	"bytes"
	"github.com/johnlion/spider/core/common/util"
)

type DownloaderHttp struct {

}

// The HttpDownloader download page by package net/http.
// The "html" content is contained in dom parser of package goquery.
// The "json" content is saved.
// The "jsonp" content is modified to json.
// The "text" content will save body plain text only.
// The page result is saved in Page.
func NewDownloaderHttp () *DownloaderHttp{
	return &DownloaderHttp{}
}

// Interface
func (this *DownloaderHttp )  Download( req *request.Request ) *page.Page{
	var mtype string
	var p = page.NewPage( req )
	mtype = req.GetResponceType()
	switch mtype {
	case "html":
		return this.downloadHtml( p, req )
		break
	case "json":
		return this.downloadJson( p, req )
		break
	case "jsonp":
		//return this.downloadJsonp( p , req )
		break
	case "text":
		return this.downloadText( p, req )
		break
	default:
		xlog.LogInst().LogError( "error request type:" + mtype )
	}

	return p
}

// Charset auto determine. Use golang.org/x/net/html/charset. Get page body and change it to utf-8
func (this *DownloaderHttp) changeCharsetEncodingAuto(contentTypeStr string, sor io.ReadCloser) string {
	var err error
	destReader, err := charset.NewReader(sor, contentTypeStr)

	if err != nil {
		xlog.LogInst().LogError(err.Error())
		destReader = sor
	}

	var sorbody []byte
	if sorbody, err = ioutil.ReadAll(destReader); err != nil {
		xlog.LogInst().LogError(err.Error())
		// For gb2312, an error will be returned.
		// Error like: simplifiedchinese: invalid GBK encoding
		// return ""
	}
	//e,name,certain := charset.DetermineEncoding(sorbody,contentTypeStr)
	bodystr := string(sorbody)

	return bodystr
}

func (this *DownloaderHttp) changeCharsetEncodingAutoGzipSupport(contentTypeStr string, sor io.ReadCloser) string {
	var err error
	gzipReader, err := gzip.NewReader(sor)
	if err != nil {
		xlog.LogInst().LogError(err.Error())
		return ""
	}
	defer gzipReader.Close()
	destReader, err := charset.NewReader(gzipReader, contentTypeStr)

	if err != nil {
		xlog.LogInst().LogError(err.Error())
		destReader = sor
	}

	var sorbody []byte
	if sorbody, err = ioutil.ReadAll(destReader); err != nil {
		xlog.LogInst().LogError(err.Error())
		// For gb2312, an error will be returned.
		// Error like: simplifiedchinese: invalid GBK encoding
		// return ""
	}
	//e,name,certain := charset.DetermineEncoding(sorbody,contentTypeStr)
	bodystr := string(sorbody)
	return bodystr
}

func connectByHttp( p *page.Page, req *request.Request ) ( *http.Response, error ){
	client := &http.Client{
		CheckRedirect: req.GetRedirectFunc(),
	}

	httpreq, err := http.NewRequest( req.GetMethod() , req.GetUrl(), strings.NewReader( req.GetPostdata() ) )
	if header := req.GetHeader() ; header != nil{
		httpreq.Header = req.GetHeader()
	}

	if cookies := req.GetCookies(); cookies != nil{
		for i := range cookies{
			httpreq.AddCookie( cookies[i] )
		}
	}

	var resp *http.Response
	if resp, err = client.Do( httpreq ); err != nil {
		if e, ok := err.( *url.Error ); ok && e.Err != nil && e.Err.Error() == "normail"{
			// normal
		}else{
			xlog.LogInst().LogError( err.Error() )
			p.SetStatus( true, err.Error() )
			return nil , err
		}
	}

	return resp, nil
}

func connectByHttpProxy( p *page.Page, in_req *request.Request ) ( *http.Response, error ){
	request , _ := http.NewRequest( "GET", in_req.GetUrl(), nil )
	proxy, err := url.Parse( in_req.GetProxyHost() )
	if err != nil {
		return nil , err
	}
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL( proxy ),
		},
	}
	resp, err := client.Do( request )
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func ( this *DownloaderHttp ) downloadFile( p *page.Page , req *request.Request )  ( *page.Page, string ){
	var err error
	var urlstr string
	if urlstr = req.GetUrl(); len( urlstr ) == 0{
		xlog.LogInst().LogError( "url is empty" )
		p.SetStatus( true, "url is empty" )
	}

	var resp *http.Response

	if proxystr := req.GetProxyHost(); len( proxystr ) !=0{
		resp, err = connectByHttpProxy( p, req )
	}else{
		resp, err = connectByHttp( p, req )
	}

	if err != nil{
		return  p, ""
	}

	p.SetHeader( resp.Header )
	p.SetCookies( resp.Cookies() )

	var bodyStr string
	if resp.Header.Get( "Content-Encoding" ) == "gzip"{
		bodyStr = this.changeCharsetEncodingAutoGzipSupport( resp.Header.Get("Content-Type"), resp.Body )
	}else{
		bodyStr = this.changeCharsetEncodingAuto(resp.Header.Get("Content-Type"), resp.Body)
	}
	defer resp.Body.Close()
	return p , bodyStr

}

func ( this *DownloaderHttp ) downloadHtml( p *page.Page, req * request.Request ) *page.Page{
	var err error
	p, destbody := this.downloadFile(p, req)
	//fmt.Printf("Destbody %v \r\n", destbody)
	if !p.IsSucc() {
		//fmt.Print("Page error \r\n")
		return p
	}
	bodyReader := bytes.NewReader([]byte(destbody))

	var doc *goquery.Document
	if doc, err = goquery.NewDocumentFromReader(bodyReader); err != nil {
		xlog.LogInst().LogError(err.Error())
		p.SetStatus(true, err.Error())
		return p
	}

	var body string
	if body, err = doc.Html(); err != nil {
		xlog.LogInst().LogError(err.Error())
		p.SetStatus(true, err.Error())
		return p
	}

	p.SetBodyStr(body).SetHtmlParser(doc).SetStatus(false, "")

	return p

}

func (this *DownloaderHttp) downloadJson(p *page.Page, req *request.Request) *page.Page {
	var err error
	p, destbody := this.downloadFile(p, req)
	if !p.IsSucc() {
		return p
	}

	var body []byte
	body = []byte(destbody)
	mtype := req.GetResponceType()
	if mtype == "jsonp" {
		tmpstr := util.JsonpToJson(destbody)
		body = []byte(tmpstr)
	}

	var r *simplejson.Json
	if r, err = simplejson.NewJson(body); err != nil {
		xlog.LogInst().LogError(string(body) + "\t" + err.Error())
		p.SetStatus(true, err.Error())
		return p
	}

	// json result
	p.SetBodyStr(string(body)).SetJson(r).SetStatus(false, "")

	return p
}

func (this *DownloaderHttp) downloadText(p *page.Page, req *request.Request) *page.Page {
	p, destbody := this.downloadFile(p, req)
	if !p.IsSucc() {
		return p
	}

	p.SetBodyStr(destbody).SetStatus(false, "")
	return p
}


