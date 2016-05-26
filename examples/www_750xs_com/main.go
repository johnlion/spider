package main
import(
	"github.com/johnlion/spider/core/common/xlog"


)




func main() {
	xlog.LogInst().Open()
	xlog.LogInst().LogInfo("Spider is running ... ")

	xlog.LogInst().Close()
}