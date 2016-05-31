package resourceManage_test
import(
	"testing"
	"github.com/johnlion/spider/core/common/resourceManage"
)

func TestResourceManage( t *testing.T ){
	var mc *resourceManage.ResourceManageChan
	mc = resourceManage.NewResourceManageChan(1)
	mc.GetOne()
	println(mc.Has())
	mc.FreeOne()
	println(mc.Has())
	mc.GetOne()
	println(mc.Has())
	mc.FreeOne()
	println(mc.Has())
	mc.GetOne()
	println(mc.Has())
	mc.FreeOne()
	println(mc.Has())
}
