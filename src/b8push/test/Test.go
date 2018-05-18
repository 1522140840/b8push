package main

import (
	"fmt"
)

type  A struct {
	Name int
	X string
	Y int
}


func main() {
	var m map[string]bool

	fmt.Println(len(m))
}

func test(a A){

	if a==(A{}){
		fmt.Print(a)
	}
}
