package main

import (
	"fmt"
	"go-zero-apr-mgr/mvc"
	za "go-zero-apr-mgr/zero-apr-lib"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Invalid arguments. Expected json file path.")
		return
	}

	dataFile := os.Args[1]

	da, err := za.Connect(dataFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	mvc.ServerMain(da)
}
