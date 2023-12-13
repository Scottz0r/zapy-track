package main

import (
	"fmt"
	"go-zero-apr-mgr/mvc"
	za "go-zero-apr-mgr/zero-apr-lib"
)

func main() {
	da, err := za.Connect("test.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	// err = consoleLoop(da)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	mvc.ServerMain(da)
}
