package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/USACE/go-consequences/compute"
	"github.com/USACE/go-consequences/hazardproviders"
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/resultswriters"
	"github.com/USACE/go-consequences/structureprovider"
)

/*
//Config describes the configuration settings for go-consequences.

	type Config struct {
		SkipJWT       bool
		LambdaContext bool
		DBUser        string
		DBPass        string
		DBName        string
		DBHost        string
		DBSSLMode     string
	}
*/
func main() {
	fp := os.Args[1]
	b, err := os.ReadFile(fp)
	if err != nil {
		log.Fatal(err)
	}
	var config compute.Config
	json.Unmarshal(b, &config)
	computable, err := config.CreateComputable()
	if err != nil {
		log.Fatal(err)
	}
	defer computable.ResultsWriter.Close()
	defer computable.HazardProvider.Close()
	err = computable.Compute()
	if err != nil {
		log.Fatal(err)
	}
}

func main2() {
	//initialize the NSI API structure provider
	// nsp := structureprovider.InitNSISP()
	nsp, _ := structureprovider.InitStructureProvider("/workspaces/go-consequences/data/burlington-davenport-nsi.gpkg", "nsi-clipped", "GPKG")
	nsp.SetDeterministic(true)
	now := time.Now()
	fmt.Println(now)
	//nsp.SetDeterministic(true)
	//identify the depth grid to apply to the structures.
	root := "/workspaces/go-consequences/data/burlington-davenport-100yr-clipped"
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
	fmt.Println("running compute.StreamAbstract")
	compute.StreamAbstract(dfr, nsp, w)
	// compute.StreamAbstractMultiVariate(dfr, nsp, w)
	fmt.Println(time.Since(now))
}
