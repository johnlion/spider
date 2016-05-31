package resourceManage

type ResourceMange interface {
	Getone()
	FreeOne()
	Has() uint
	Left() uint
}