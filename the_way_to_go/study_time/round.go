package main

import (
	gtime "github.com/og/x/time"
	"log"
	"time"
)

func main()  {
	// why
	t,_:= time.ParseInLocation(gtime.Second, "2016-06-13 15:23:36", time.Local)
	log.Print(t.Format(gtime.Second))
	for i:=1;i<=10;i++ {
		log.Print(i ," : ", t.Round(time.Duration(i)*time.Hour).Format(gtime.Second) )
	}

}

