package commands

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/dtluna/nachmundtracker/business"
)

type Default struct {
	CampaignYAML string `arg:"" name:"campaign_yaml" help:"Path to the campaign YAML file." type:"existingfile"`
	Phase        string `enum:"1,2,3,all,a" default:"all" short:"p" help:"Phase to filter by. One of: 1, 2, 3, all (short: a)."`
	Alliance     string `enum:"all,a,despoilers,d,guardians,g,marauders,m" default:"all" short:"a" help:"Alliance to filter by. One of: all (short: a), despoilers (short: d), guardians (short: g), marauders (short: m)."`
}

func (com *Default) Run(ctx *kong.Context) error {
	games, err := business.DecodeData(com.CampaignYAML)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: validation errors present, run the validate command to see details\n")
	}

	fmt.Printf("%+v\n", games)
	return nil
}
