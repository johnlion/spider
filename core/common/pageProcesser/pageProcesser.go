package pageProcesser

import "github.com/johnlion/spider/core/common/page"

type PageProcesser interface{
	Processer(p *page.Page )
	Finish()
}
