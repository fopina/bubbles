//go:generate go run -tags=dev data/assets_generate.go
package main

import (
	"github.com/fopina/bubbles/cmd"
)

func main() {
	cmd.Run()
}
