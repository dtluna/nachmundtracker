package business

import (
	"fmt"
	"io"
	"maps"
	"os"
	"slices"

	"github.com/dtluna/nachmundtracker/model"
	"github.com/goccy/go-yaml"
)

func DecodeData(filename string) ([]model.GameRecord, error) {
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

	if err := validate(games); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	return games, nil
}

type GameError struct {
	game        model.GameRecord
	description string
}

func (ge GameError) Error() string {
	return fmt.Sprintf("%v: %s", ge.game, ge.description)
}

func validate(games []model.GameRecord) error {
	for _, game := range games {
		if game.Phase > 3 {
			return GameError{
				game:        game,
				description: "phase cannot be more than 3",
			}
		}

		switch game.Scale {
		case model.ScaleIncursion, model.ScaleStrikeForce, model.ScaleOnslaught:
		default:
			return GameError{
				game:        game,
				description: fmt.Sprintf("incorrect scale %q", game.Scale),
			}
		}

		if len(game.Players) != 2 {
			return GameError{
				game:        game,
				description: "game must have 2 players",
			}
		}

		if len(game.Alliances) != 2 {
			return GameError{
				game:        game,
				description: "game must have 2 alliances",
			}
		}

		for _, alliance := range game.Alliances {
			switch alliance {
			case model.AllianceGuardians, model.AllianceMarauders, model.AllianceDespoilers:
			default:
				return GameError{
					game:        game,
					description: fmt.Sprintf("invalid alliance entry: %q", alliance),
				}
			}
		}

		if game.Alliances[0] == game.Alliances[1] {
			return GameError{
				game:        game,
				description: "game must have a 2 different alliances",
			}
		}

		if game.Victor == "" {
			return GameError{
				game:        game,
				description: "victor is required",
			}
		}

		switch game.Victor {
		case model.VictorDraw, model.VictorGuardians, model.VictorDespoilers, model.VictorMarauders:
		default:
			return GameError{
				game:        game,
				description: fmt.Sprintf("invalid victor: %q", game.Victor),
			}
		}

		if game.Victor != model.VictorDraw && !slices.Contains(game.Alliances, model.Alliance(game.Victor)) {
			return GameError{
				game:        game,
				description: fmt.Sprintf("victor must be one of %v or %s", game.Alliances, model.VictorDraw),
			}
		}

		if len(game.BPAllocation) != 2 {
			return GameError{
				game:        game,
				description: "bp_allocation must have 2 entries",
			}
		}

		for alliance, location := range game.BPAllocation {
			if !slices.Contains(game.Alliances, alliance) {
				return GameError{
					game:        game,
					description: fmt.Sprintf("bp_allocation should only mention alliances present in the game"),
				}
			}

			switch location {
			case model.LocationTower, model.LocationBastion, model.LocationBattery, model.LocationSpaceport:
			default:
				return GameError{
					game:        game,
					description: fmt.Sprintf("invalid location for bp_allocation for alliance %s: %q", alliance, location),
				}
			}
		}

		for alliance := range game.SAPGain {
			if !slices.Contains(game.Alliances, alliance) {
				return GameError{
					game:        game,
					description: fmt.Sprintf("sap_gain should only mention alliances present in the game"),
				}
			}
		}

		sapGainAlliances := slices.Sorted(maps.Keys(game.SAPGain))

		for alliance, location := range game.SAPAllocation {
			if !slices.Contains(sapGainAlliances, alliance) {
				return GameError{
					game:        game,
					description: fmt.Sprintf("sap_allocation should only mention alliances present in sap_gain"),
				}
			}

			switch location {
			case model.LocationTower, model.LocationBastion, model.LocationBattery, model.LocationSpaceport:
			default:
				return GameError{
					game:        game,
					description: fmt.Sprintf("invalid location for sap_allocation for alliance %s: %q", alliance, location),
				}
			}
		}

		for _, sapGainAlliance := range sapGainAlliances {
			if _, found := game.SAPAllocation[sapGainAlliance]; !found {
				return GameError{
					game:        game,
					description: fmt.Sprintf("sap gained are unallocated for alliance %v", sapGainAlliance),
				}
			}
		}
	}

	return nil
}
