package model

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"
)

type GameRecord struct {
	Date          string
	Phase         Phase
	Scale         Scale
	Mission       string
	Players       []string
	Alliances     []Alliance
	Victor        Victor                // either a participating alliance or DRAW
	SAPGain       map[Alliance]uint     `yaml:"sap_gain"`
	BPAllocation  map[Alliance]Location `yaml:"bp_allocation"`
	SAPAllocation map[Alliance]Location `yaml:"sap_allocation"`
}

func (gr GameRecord) String() string {
	return fmt.Sprintf("game on %v, mission %s between %v", gr.Date, gr.Mission, gr.Players)
}

type Phase int

const (
	Phase1 Phase = 1
	Phase2 Phase = 2
	Phase3 Phase = 3
)

type Alliance string

const (
	AllianceGuardians  Alliance = "guardians"
	AllianceDespoilers Alliance = "despoilers"
	AllianceMarauders  Alliance = "marauders"
)

type Scale string

const (
	ScaleIncursion   Scale = "incursion"
	ScaleStrikeForce Scale = "strike_force"
	ScaleOnslaught   Scale = "onslaught"
)

type Location string

const (
	LocationBastion   Location = "bastion"
	LocationTower     Location = "tower"
	LocationBattery   Location = "battery"
	LocationSpaceport Location = "spaceport"
)

type Victor string

const (
	VictorDraw       Victor = "DRAW"
	VictorGuardians  Victor = Victor(AllianceGuardians)
	VictorDespoilers Victor = Victor(AllianceDespoilers)
	VictorMarauders  Victor = Victor(AllianceMarauders)
)

type GameErrors struct {
	game   GameRecord
	errors []error
}

func (ge GameErrors) Error() string {
	return fmt.Sprintf("%v: %v", ge.game, ge.errors)
}

func (ge *GameErrors) AddError(s string) {
	ge.errors = append(ge.errors, errors.New(s))
}

func (ge GameErrors) HasErrors() bool {
	return len(ge.errors) > 0
}

func (game GameRecord) Validate() *GameErrors {
	ge := GameErrors{game: game}

	if game.Phase > 3 {
		ge.AddError("phase cannot be more than 3")
	}

	switch game.Scale {
	case ScaleIncursion, ScaleStrikeForce, ScaleOnslaught:
	default:
		ge.AddError(fmt.Sprintf("incorrect scale %q", game.Scale))
	}

	if len(game.Players) != 2 {
		ge.AddError("game must have only 2 players")
	}

	if len(game.Alliances) != 2 {
		ge.AddError("game must have only 2 alliances")
	}

	for _, alliance := range game.Alliances {
		switch alliance {
		case AllianceGuardians, AllianceMarauders, AllianceDespoilers:
		default:
			ge.AddError(fmt.Sprintf("invalid alliance entry: %q", alliance))
		}
	}

	if game.Alliances[0] == game.Alliances[1] {
		ge.AddError("game must have a 2 different alliances")
	}

	if game.Victor == "" {
		ge.AddError("victor is required")
	}

	switch game.Victor {
	case VictorDraw, VictorGuardians, VictorDespoilers, VictorMarauders:
	default:
		ge.AddError(fmt.Sprintf("invalid victor: %q", game.Victor))
	}

	if game.Victor != VictorDraw && !slices.Contains(game.Alliances, Alliance(game.Victor)) {
		ge.AddError(fmt.Sprintf("victor must be one of %v or %s", game.Alliances, VictorDraw))
	}

	if len(game.BPAllocation) != 2 {
		ge.AddError("bp_allocation must have 2 entries")
	}

	for alliance, location := range game.BPAllocation {
		if !slices.Contains(game.Alliances, alliance) {
			ge.AddError("bp_allocation should only mention alliances present in the game")
		}

		switch location {
		case LocationTower, LocationBastion, LocationBattery, LocationSpaceport:
		default:
			ge.AddError(fmt.Sprintf("invalid location for bp_allocation for alliance %s: %q", alliance, location))
		}
	}

	for alliance := range game.SAPGain {
		if !slices.Contains(game.Alliances, alliance) {
			ge.AddError("sap_gain should only mention alliances present in the game")
		}
	}

	sapGainAlliances := slices.Sorted(maps.Keys(game.SAPGain))

	for alliance, location := range game.SAPAllocation {
		if !slices.Contains(sapGainAlliances, alliance) {
			ge.AddError("sap_allocation should only mention alliances present in the game")
		}

		switch location {
		case LocationTower, LocationBastion, LocationBattery, LocationSpaceport:
		default:
			ge.AddError(fmt.Sprintf("invalid location for sap_allocation for alliance %s: %q", alliance, location))
		}
	}

	for _, sapGainAlliance := range sapGainAlliances {
		if _, found := game.SAPAllocation[sapGainAlliance]; !found {
			ge.AddError(fmt.Sprintf("sap gained are unallocated for alliance %v", sapGainAlliance))
		}
	}

	if ge.HasErrors() {
		return &ge
	}

	return nil

}

type CampaignErrors []GameErrors

func (ce CampaignErrors) Error() string {
	builder := strings.Builder{}

	for _, ge := range ce {
		fmt.Fprintf(&builder, "%v\n", ge)
	}

	return builder.String()
}

func ValidateGames(games []GameRecord) (validGames []GameRecord, err error) {
	var ce CampaignErrors

	for _, game := range games {
		if err := game.Validate(); err != nil {
			ce = append(ce, *err)
		} else {
			validGames = append(validGames, game)
		}
	}

	if len(ce) > 0 {
		return validGames, ce
	}

	return validGames, nil
}
