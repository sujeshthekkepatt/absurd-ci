package main

import (
	"fmt"
	"gitcloner/gitutil"
)

func main() {

	fmt.Println("Git cloner version 1.0 loaded")

	gcloner := gitutil.NewClient("wasm")

	success, err := gcloner.Clone("https://github.com/go-git/go-git", "")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(success)
}
