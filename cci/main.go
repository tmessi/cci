package main

import (
	"os"

	"github.com/tmessi/cci/internal/command"
)

func main() {
	app := command.App()
	app.Run(os.Args)
}
