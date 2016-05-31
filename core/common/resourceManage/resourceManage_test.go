package resourceManage_test
import(
	"testing"
	"github.com/johnlion/spider/core/common/resourceManage"
)

func TestResourceManage( t *testing.T ){
	var mc *resourceManage.ResourceManageChan
	mc = resourceManage.NewResourceManageChan(10)
	for{
	mc.GetOne()
	println(mc.Has())
	mc.FreeOne()
	println(mc.Has())
	}

}
