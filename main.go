package main

import (
	"fmt"
	"log"

	"github.com/alecthomas/kong"
	"github.com/dtluna/nachmundtracker/business"
)

var CLI struct {
	Default DefaultCommand `cmd:"" default:"withargs"`
}

type DefaultCommand struct {
	CampaignYAML string `arg:"" name:"campaign_yaml" help:"Path to the campaign YAML file." type:"path"`
}

func (com *DefaultCommand) Run(ctx *kong.Context) error {
	games, err := business.DecodeData(com.CampaignYAML)
	if err != nil {
		log.Fatalf("decoding data: %v\n", err)
	}

	fmt.Printf("%+v\n", games)
	return nil
}

func main() {
	ctx := kong.Parse(&CLI)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
