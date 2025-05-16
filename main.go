package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/dtluna/nachmundtracker/business"
)

var CLI struct {
	Default DefaultCommand `cmd:"" help:"Show BP and SAP allocation for all phases for all alliances." default:"withargs"`
}

type DefaultCommand struct {
	CampaignYAML string `arg:"" name:"campaign_yaml" help:"Path to the campaign YAML file." type:"path"`
}

func (com *DefaultCommand) Run(ctx *kong.Context) error {
	games, err := business.DecodeData(com.CampaignYAML)
	if err != nil {
		return fmt.Errorf("decoding data: %w", err)
	}

	fmt.Printf("%+v\n", games)
	return nil
}

func main() {
	ctx := kong.Parse(&CLI)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
