package xlog
//
// xlog.LogInst().Open()
// xlog.LogInst().LogInfo("Spider is running ... ")
// xlog.LogInst().Close()
//

// Example
// xlog.LogInst().Open()
// xlog.LogInst().LogInfo("Spider is running ... ")
// xlog.LogInst().LogError("xxx")
// xlog.StraceInst().Println("This is a error!!! ...")
// ..........................................................................

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/johnlion/spider/core/common/config"
	"time"
	"strconv"
	"fmt"
)

type fileLog struct{
	xlog
	loginst *log.Logger
}

var flog *fileLog  //global var

func LogInst() *fileLog{
	if flog == nil {
		InitFileLog(false, "")
	}
	return flog
}


func InitFileLog( isopen bool, fp string ){
	if !isopen{
		flog = &fileLog{}
		flog.loginst = nil
		flog.isOpen = isopen
	}

	if fp == ""{
		wd := config.XLOG_DIR

		if wd ==""{
			file, _ := exec.LookPath( os.Args[0] )
			path := filepath.Dir( file )
			wd = path
		}
		if wd == ""{
			panic ( "GOPATH is not setted in env or can not get exe path." )
		}
		fp = wd + "/" + "log"
	}
	flog = newFileLog( isopen, fp )
}


func newFileLog( isopen bool, logpath string ) *fileLog {
	year, month, day := time.Now().Date()
	filename := "log." + strconv.Itoa( year ) + "-" + strconv.Itoa( int( month) ) + "-" + strconv.Itoa( day )
	err := os.MkdirAll( logpath, 0755 )

	if err != nil{
		panic( "logpath error : " + logpath + "\n" )
	}

	f, err := os.OpenFile( logpath+"/"+filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644 )
	if err != nil{
		panic("log file open error : " + logpath + "/" + filename + "\n")
	}
	fmt.Println(  logpath + "/" + filename  )

	pfilelog := &fileLog{}
	pfilelog.loginst = log.New( f, "", log.LstdFlags  )
	return pfilelog

}

func ( this *fileLog ) log ( lable string, str string ){
	if !this.isOpen{
		return
	}
	file , line := this.getCaller()
	this.loginst.Printf(  "%s:%d: %s %s\n", file, line, lable, str )
}

func ( this *fileLog ) LogInfo( str string ){
	this.log( "[INFO]", str )
}

func ( this *fileLog ) LogError( str string ){
	this .log( "[ERROR]", str )
}

