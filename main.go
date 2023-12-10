package main

import (
	"fmt"
	zeroaprlib "go-zero-apr-mgr/zero-apr-lib"
	//"time"
)

func main() {
	da, err := zeroaprlib.Connect("test.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	// err = consoleLoop(da)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	serverMain(da)
}
