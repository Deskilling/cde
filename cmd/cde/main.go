package main

import (
	"fmt"

	"cde/internal/core"
	"cde/internal/editor/zed"
)

func init() {
	core.InitLogger(-4)
}

func main() {
	p, _ := zed.ZedExtractWorkspacePath()
	fmt.Println(p)
}
