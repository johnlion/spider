package main
import(
	"github.com/johnlion/spider/core/common/xlog"

	"github.com/johnlion/spider/core/common/spider"
	"fmt"
	"github.com/johnlion/spider/core/common/request"
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
func ( this *MyPageProcesser ) Processer(){}

// Interface
func ( this *MyPageProcesser ) Finish(){
	fmt.Println( "TODO:before end spider \r\n" )
}




func main() {
	xlog.LogInst().Open()
	xlog.LogInst().LogInfo("Spider is running ... ")
	xlog.LogInst().Close()
	// Spider input
	// PageProcesser
	// Task name used in Pipeline for record;
	sp := spider.NewSpider( "www.75xs.com",  NewMyPageProcesser() )

	req := request.NewRequest( "http://baike.baidu.com/view/1628025.htm?fromtitle=http&fromid=243074&type=syn",
		"html", "", "GET", "", nil, nil, nil, nil)
	sp.GetByRequest( req )



	fmt.Printf("%s", "---------------------------------------------------------------------------------------------------\r\n")
	//fmt.Printf( "%$v\n", req )
	//fmt.Printf( "%v\n", sp )
}