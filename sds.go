package sds

import (
	"fmt"
)

var Version string = "1.0.0"

func PrintVersion() {
	fmt.Println("[gosds] " + Version)
}
