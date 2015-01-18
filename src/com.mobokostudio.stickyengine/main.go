package main

import (
	"fmt"
	//	"com.mobokostudio.engine/util"
	//	"com.mobokostudio.engine/samples"
	//	"github.com/op/go-logging"
)

//var log = logging.MustGetLogger("example")

func init() {
	start_process()
}

func main() {
	start_process()

}

func start_process() {
	defer func() {
		if str := recover(); str != nil {
			fmt.Println(str)
		}
	}()

	startHTTP()

	//	util.PrintDir(".")
	//
	//	fmt.Println(util.GetFileData("src/com.mobokostudio.engine/app.yaml"))
	//	util.WriteFile("src/com.mobokostudio.engine/test.txt","Hello, Wordfdfsdfld!!")
	//	fmt.Println(util.GetFileData("src/com.mobokostudio.engine/test.txt"))
	//
	//	samples.GoroutineSample()
	//	samples.WebSample()

	//	fmt.Println("finish")
	//	fmt.Println("finish3")

}
