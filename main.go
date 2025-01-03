package main

import (
	"fmt"
	"time"

	"github.com/USACE/go-consequences/compute"
	"github.com/USACE/go-consequences/hazardproviders"
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/resultswriters"
	"github.com/USACE/go-consequences/structureprovider"
)

func main() {
	//initialize the NSI API structure provider
	nsp := structureprovider.InitNSISP()
	now := time.Now()
	fmt.Println(now)
	//nsp.SetDeterministic(true)
	//identify the depth grid to apply to the structures.
	root := "/workspaces/go-consequences/data/burlington-davenport-100yr"
	filepath := root + ".tif"
	w, _ := resultswriters.InitSpatialResultsWriter(root+"_consequencesGHG.gpkg", "results", "GPKG")
	//w := consequences.InitSummaryResultsWriterFromFile(root + "_consequences_SUMMARY.json")
	//create a result writer based on the name of the depth grid.
	//w, _ := resultswriters.InitGpkResultsWriter(root+"_consequences_nsi.gpkg", "nsi_result")
	defer w.Close()
	//initialize a hazard provider based on the depth grid.
	dfr, _ := hazardproviders.Init_CustomFunction(filepath, func(valueIn hazards.HazardData, hazard hazards.HazardEvent) (hazards.HazardEvent, error) {
		if valueIn.Depth == 0 {
			return hazard, hazardproviders.NoHazardFoundError{}
		}
		process := hazardproviders.DepthHazardFunction()
		return process(valueIn, hazard)
	})
	//compute consequences.
	compute.StreamAbstract(dfr, nsp, w)
	fmt.Println(time.Since(now))
}
