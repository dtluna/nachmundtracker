package model

type GameRecord struct {
	Phase         Phase
	Scale         Scale
	Mission       string
	Players       []string
	Alliances     []Alliance
	Victor        string                // either a participating alliance or DRAW
	SAPGain       SAPGain               `yaml:"sap_gain"`
	BPAllocation  map[Alliance]Location `yaml:"bp_allocation"`
	SAPAllocation map[Alliance]Location `yaml:"sap_allocation"`
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

type SAPGain struct {
	Guardians, Despoilers, Marauders int
}

type Location string

const (
	LocationBastion   Location = "bastion"
	LocationTower     Location = "tower"
	LocationBattery   Location = "battery"
	LocationSpaceport Location = "spaceport"
)
