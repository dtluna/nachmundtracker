package commands

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/dtluna/nachmundtracker/business"
)

type Validate struct {
	CampaignYAML string `arg:"" name:"campaign_yaml" help:"Path to the campaign YAML file." type:"path"`
}

func (val *Validate) Run(ctx *kong.Context) error {
	_, err := business.DecodeData(val.CampaignYAML)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("campaign file valid")
	}

	return nil
}
