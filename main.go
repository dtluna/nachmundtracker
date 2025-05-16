package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/dtluna/nachmundtracker/business"
)

var CLI struct {
	Default  DefaultCommand  `cmd:"" help:"Show BP and SAP allocation for all phases for all alliances." default:"withargs"`
	Validate ValidateCommand `cmd:"" help:"Validate the campaign file."`
}

type DefaultCommand struct {
	CampaignYAML string `arg:"" name:"campaign_yaml" help:"Path to the campaign YAML file." type:"path"`
}

func (com *DefaultCommand) Run(ctx *kong.Context) error {
	games, err := business.DecodeData(com.CampaignYAML)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: validation errors present, run the validate command to see details\n")
	}

	fmt.Printf("%+v\n", games)
	return nil
}

type ValidateCommand struct {
	CampaignYAML string `arg:"" name:"campaign_yaml" help:"Path to the campaign YAML file." type:"path"`
}

func (val *ValidateCommand) Run(ctx *kong.Context) error {
	_, err := business.DecodeData(val.CampaignYAML)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("campaign file valid")
	}

	return nil
}

func main() {
	ctx := kong.Parse(&CLI)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
