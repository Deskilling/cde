package main

import (
	"fmt"

	"cde/internal/core"
	"cde/internal/programs/zed"
)

func init() {
	core.InitLogger(-4)
}

func main() {
	p, _ := zed.ZedExtractWorkspacePath()
	fmt.Println(p)
}
