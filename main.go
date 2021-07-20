package main

import (
	//"dataflowobfuscation"
	"dataflowobfuscation"
	//"local2state"
	"fmt"
	"os"
	//"strconv"
)

func main() {
	//fmt.Println("Hello world")
	if len(os.Args) != 3 {
		fmt.Println("Wrong parameters!")
		return
	}
	var args [2]string
	for i, arg := range os.Args[1:] {
		args[i] = arg
	}

	dfo := dataflowobfuscation.NewDfo(args[0], args[1], "Configuration.json")
	//fmt.Printf("%+v\n",dfo)
	var _ = dfo.Run()
}
