package downloader

import (
	"github.com/johnlion/spider/core/common/request"
	"github.com/johnlion/spider/core/common/page"
)

type Downloader interface {
	Download( req *request.Request ) *page.Page
}