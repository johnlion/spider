package xlog

import "runtime"

type xlog struct {
	isOpen bool
}

func ( *xlog ) getCaller()  ( string,int ){
	_, file, line, ok := runtime.Caller( 3 )
	if !ok{
		file = "???"
		line = 0
	}

	return file, line
}

func ( this *xlog ) Open(){
	this.isOpen = true
}

func ( this *xlog ) Close(){
	this.isOpen = false
}

