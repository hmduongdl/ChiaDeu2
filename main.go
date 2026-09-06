package main

import (
	"fmt"

	"chiadeu/internal/cli"
)

func main() {
	fmt.Println("Chạy Chương Trình")
	app := cli.NewCLIApp()
	app.Run()
}
