package main

import (

)
import (
	"os"
	"strconv"
	"time"
	"fmt"
	"github.com/johnlion/spider/core/common/config"
)

func newFileLog( isopen bool, logpath string )  {
	year, month, day := time.Now().Date()
	filename := "log." + strconv.Itoa( year ) + "-" + strconv.Itoa( int( month) ) + "-" + strconv.Itoa( day )
	err := os.MkdirAll( logpath, 0755 )


	if err != nil{
		panic( "logpath error : " + logpath + "\n" )
	}

	_, err = os.OpenFile( logpath+"/"+filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644 )
	if err != nil{
		panic("log file open error : " + logpath + "/" + filename + "\n")
	}
	fmt.Println(  logpath + "/" + filename  )


}


func main(){

	newFileLog( false, config.XLOG_DIR + "/" + "log" )

}
