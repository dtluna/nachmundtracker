package presentation

import (
	"fmt"
	"strconv"

	"github.com/dtluna/nachmundtracker/business"
	"github.com/dtluna/nachmundtracker/model"
	"github.com/fatih/color"
	"github.com/rodaine/table"
)

func PrintResults(results business.Results, phase string, alliance string) {
	switch phase {
	case "a", "all":
		for _, phaseInt := range []model.Phase{model.Phase1, model.Phase2, model.Phase3} {
			printPhase(int(phaseInt), results[phaseInt], alliance)
			fmt.Println("")
		}
	default:
		phaseInt, _ := strconv.Atoi(phase)
		printPhase(phaseInt, results[model.Phase(phaseInt)], alliance)
	}
}

var headers = []any{"Alliance", "Bastion", "Battery", "Spaceport", "Tower", "Total"}

func printPhase(phase int, phaseResults business.PhaseResults, alliance string) {
	color.New(color.FgHiGreen, color.Underline).Printf("Phase %v\n", phase)

	bpTable := table.New(headers...)
	bpTable.WithHeaderFormatter(
		color.New(color.FgRed, color.Underline).SprintfFunc(),
	).WithFirstColumnFormatter(
		color.New(color.FgYellow).SprintfFunc(),
	)

	sapTable := table.New(headers...)
	sapTable.WithHeaderFormatter(
		color.New(color.FgBlue, color.Underline).SprintfFunc(),
	).WithFirstColumnFormatter(
		color.New(color.FgYellow).SprintfFunc(),
	)

	switch alliance {
	case "a", "all":
		for alliance, allianceResults := range map[string]business.AllianceResults{
			"Guardians":  phaseResults.Guardians,
			"Despoilers": phaseResults.Despoilers,
			"Marauders":  phaseResults.Marauders,
		} {
			bpResults := allianceResults.BPAllocation
			sapResults := allianceResults.SAPAllocation

			addResults(bpTable, sapTable, bpResults, sapResults, alliance)
		}

	case "g", "guardians":
		bpResults := phaseResults.Guardians.BPAllocation
		sapResults := phaseResults.Guardians.SAPAllocation

		addResults(bpTable, sapTable, bpResults, sapResults, "Guardians")

	case "d", "despoilers":
		bpResults := phaseResults.Despoilers.BPAllocation
		sapResults := phaseResults.Despoilers.SAPAllocation

		addResults(bpTable, sapTable, bpResults, sapResults, "Despoilers")

	case "m", "marauders":
		bpResults := phaseResults.Marauders.BPAllocation
		sapResults := phaseResults.Marauders.SAPAllocation

		addResults(bpTable, sapTable, bpResults, sapResults, "Marauders")
	}

	color.New(color.FgRed).Println("BP Allocation")
	bpTable.Print()

	fmt.Println("")

	color.New(color.FgBlue).Println("SAP Allocation")
	sapTable.Print()
}

func addResults(bpTable, sapTable table.Table, bpResults business.LocationBPs, sapResults business.LocationSAPs, alliance string) {
	bpTable.AddRow(alliance, bpResults.Bastion, bpResults.Battery, bpResults.Spaceport, bpResults.Tower, bpResults.Total())
	sapTable.AddRow(alliance, sapResults.Bastion, sapResults.Battery, sapResults.Spaceport, sapResults.Tower, sapResults.Total())
}
