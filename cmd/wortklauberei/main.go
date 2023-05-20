package main

import (
	"fmt"
	"os"

	"github.com/yawn77/wortklauberei/controllers"
)

var version string

func main() {
	gc, err := controllers.NewGameController(false, version)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	gc.Run()
}
