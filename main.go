package main

import (

)
import "github.com/johnlion/spider/core/common/resourceManage"

func main() {
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
