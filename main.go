package main

import (
	"github.com/alecthomas/kong"
	"github.com/dtluna/nachmundtracker/commands"
)

var CLI struct {
	Default  commands.Default  `cmd:"" help:"Show BP and SAP allocation. Output can be filtered by phase and alliance." default:"withargs"`
	Validate commands.Validate `cmd:"" help:"Validate the campaign file." aliases:"v"`
}

func main() {
	ctx := kong.Parse(&CLI)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
