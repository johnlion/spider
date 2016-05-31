package downloader
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
func (this *DownloaderHttp ) Download(){

}