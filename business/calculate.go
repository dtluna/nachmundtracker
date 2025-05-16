package business

import "github.com/dtluna/nachmundtracker/model"

type LocationBPs struct {
	Tower     uint
	Spaceport uint
	Battery   uint
	Bastion   uint
}

func (lb *LocationBPs) update(location model.Location, bpGain uint) {
	switch location {
	case model.LocationBastion:
		lb.Bastion += bpGain
	case model.LocationBattery:
		lb.Battery += bpGain
	case model.LocationTower:
		lb.Tower += bpGain
	case model.LocationSpaceport:
		lb.Spaceport += bpGain
	}
}

func (lb LocationBPs) Total() uint {
	return lb.Bastion + lb.Battery + lb.Spaceport + lb.Tower
}

type LocationSAPs struct {
	Tower     uint
	Spaceport uint
	Battery   uint
	Bastion   uint
}

func (ls *LocationSAPs) update(location model.Location, saps uint) {
	switch location {
	case model.LocationBastion:
		ls.Bastion += saps
	case model.LocationBattery:
		ls.Battery += saps
	case model.LocationTower:
		ls.Tower += saps
	case model.LocationSpaceport:
		ls.Spaceport += saps
	}
}

func (ls LocationSAPs) Total() uint {
	return ls.Bastion + ls.Battery + ls.Spaceport + ls.Tower
}

type AllianceResults struct {
	BPAllocation  LocationBPs
	SAPAllocation LocationSAPs
}

type PhaseResults struct {
	Guardians  AllianceResults
	Despoilers AllianceResults
	Marauders  AllianceResults
}

func (pr PhaseResults) update(game model.GameRecord) PhaseResults {
	for _, alliance := range game.Alliances {
		bpLocation := game.BPAllocation[alliance]
		bpGain := calculateBPGain(alliance, game.Victor, game.Scale)
		sapLocation := game.SAPAllocation[alliance]

		switch alliance {
		case model.AllianceGuardians:
			pr.Guardians.BPAllocation.update(bpLocation, bpGain)
			pr.Guardians.SAPAllocation.update(sapLocation, game.SAPGain[model.AllianceGuardians])
		case model.AllianceDespoilers:
			pr.Despoilers.BPAllocation.update(bpLocation, bpGain)
			pr.Despoilers.SAPAllocation.update(sapLocation, game.SAPGain[model.AllianceDespoilers])
		case model.AllianceMarauders:
			pr.Marauders.BPAllocation.update(bpLocation, bpGain)
			pr.Marauders.SAPAllocation.update(sapLocation, game.SAPGain[model.AllianceMarauders])
		}

	}

	return pr
}

var (
	winBP = map[model.Scale]uint{
		model.ScaleIncursion:   2,
		model.ScaleStrikeForce: 3,
		model.ScaleOnslaught:   4,
	}

	drawBP = map[model.Scale]uint{
		model.ScaleIncursion:   2,
		model.ScaleStrikeForce: 2,
		model.ScaleOnslaught:   3,
	}

	lossBP = map[model.Scale]uint{
		model.ScaleIncursion:   1,
		model.ScaleStrikeForce: 1,
		model.ScaleOnslaught:   2,
	}
)

func calculateBPGain(alliance model.Alliance, victor model.Victor, scale model.Scale) uint {
	if victor == model.VictorDraw {
		return drawBP[scale]
	}

	if alliance == model.Alliance(victor) {
		return winBP[scale]
	}

	return lossBP[scale]
}

type Results map[model.Phase]PhaseResults

func CalculateResults(games []model.GameRecord) Results {
	results := map[model.Phase]PhaseResults{
		model.Phase1: {},
		model.Phase2: {},
		model.Phase3: {},
	}

	for _, game := range games {
		phaseResults := results[game.Phase]
		results[game.Phase] = phaseResults.update(game)
	}

	return results
}
