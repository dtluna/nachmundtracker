package model

import "fmt"

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
