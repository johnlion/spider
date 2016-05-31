package resourceManage

type ResourceMange interface {
	GetOne()
	FreeOne()
	Has() uint
	Left() uint
}