package business

import (
	"errors"
	"fmt"
	"io"
	"maps"
	"os"
	"slices"
	"strings"

	"github.com/dtluna/nachmundtracker/model"
	"github.com/goccy/go-yaml"
)

func DecodeData(filename string) (validGames []model.GameRecord, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("reading file data: %w", err)
	}

	var games []model.GameRecord

	if err := yaml.Unmarshal(data, &games); err != nil {
		return nil, fmt.Errorf("decoding yaml: %w", err)
	}

	validGames, err = validateGames(games)
	if err != nil {
		return validGames, fmt.Errorf("validation errors: %w", err)
	}

	return validGames, nil
}

type GameErrors struct {
	game   model.GameRecord
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

func validateGame(game model.GameRecord) *GameErrors {
	ge := GameErrors{game: game}

	if game.Phase > 3 {
		ge.AddError("phase cannot be more than 3")
	}

	switch game.Scale {
	case model.ScaleIncursion, model.ScaleStrikeForce, model.ScaleOnslaught:
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
		case model.AllianceGuardians, model.AllianceMarauders, model.AllianceDespoilers:
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
	case model.VictorDraw, model.VictorGuardians, model.VictorDespoilers, model.VictorMarauders:
	default:
		ge.AddError(fmt.Sprintf("invalid victor: %q", game.Victor))
	}

	if game.Victor != model.VictorDraw && !slices.Contains(game.Alliances, model.Alliance(game.Victor)) {
		ge.AddError(fmt.Sprintf("victor must be one of %v or %s", game.Alliances, model.VictorDraw))
	}

	if len(game.BPAllocation) != 2 {
		ge.AddError("bp_allocation must have 2 entries")
	}

	for alliance, location := range game.BPAllocation {
		if !slices.Contains(game.Alliances, alliance) {
			ge.AddError("bp_allocation should only mention alliances present in the game")
		}

		switch location {
		case model.LocationTower, model.LocationBastion, model.LocationBattery, model.LocationSpaceport:
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
		case model.LocationTower, model.LocationBastion, model.LocationBattery, model.LocationSpaceport:
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

func validateGames(games []model.GameRecord) (validGames []model.GameRecord, err error) {
	var ce CampaignErrors

	for _, game := range games {
		if err := validateGame(game); err != nil {
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
