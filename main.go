//go:build !gui && !terminal
// +build !gui,!terminal

package main

import (
	"fmt"
)

func main() {

	fmt.Println("Please specify a build tag (terminal | gui)")
}
