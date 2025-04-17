package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/USACE/go-consequences/compute"
	"github.com/USACE/go-consequences/hazardproviders"
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

func get_files(file_list []string) <-chan string {

	out := make(chan string)

	go func() {
		for _, f := range file_list {
			out <- f
		}
		close(out)
	}()
	return out
}

func process_file(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		for i := range in {
			ts := time.Now()
			compute_FathomMultiFrequency(i, "2020", "FLUVIAL-DEFENDED_KNOWN")
			te := time.Since(ts)
			out_str := fmt.Sprintf("Processed file: %s in %s", i, te)
			out <- out_str
		}
		close(out)
	}()
	return out
}

func merge(cs ...<-chan string) <-chan string {
	var wg sync.WaitGroup
	out := make(chan string)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan string) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	content, err := os.ReadFile("/workspaces/go-consequences/data/testgrids.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var file_list []string
	err = json.Unmarshal(content, &file_list)
	if err != nil {
		log.Fatal("Error during Unmarshal():", err)
	}

	c := get_files(file_list)

	ts := time.Now()

	c1 := process_file(c)
	c2 := process_file(c)
	c3 := process_file(c)
	c4 := process_file(c)
	c5 := process_file(c)

	for i := range merge(c1, c2, c3, c4, c5) {
		fmt.Println(i)
	}

	te := time.Since(ts)
	fmt.Printf("All files completed in %s\n", te)
}

func compute_FathomMultiFrequency(filename string, year string, scenario string) {

	// year := "2020"
	// year :- 2050-SSP5_8.5
	// scenario := "FLUVIAL-DEFENDED_KNOWN"

	dataset := filename[:len(filename)-4]

	//initialize the NSI API structure provider
	nsp, err := structureprovider.InitStructureProvider("/workspaces/go-consequences/data/nsi/nsi_2022.gpkg", "nsi(shape)", "GPKG")
	if err != nil {
		log.Fatal(err)
	}

	// nsp := structureprovider.InitNSISP()

	//initialize a set of frequencies
	rps := []int{5, 10, 20, 50, 100, 200, 500, 1000}
	frequencies := []float64{.2, .1, .05, .02, .01, .005, .002, .001}

	root := "/workspaces/go-consequences/data/fathom/2020/"
	//identify the depth grids to represent the frequencies.
	hazardProviders := make([]hazardproviders.HazardProvider, len(rps))

	for i, r := range rps {
		file := fmt.Sprintf("%sFLOOD_MAP-1_3ARCSEC-NW_OFFSET-1in%d-%s-DEPTH-%s-PERCENTILE50-v3.1/%s.tif", root, r, scenario, year, dataset)
		// fmt.Println(file)
		hp, err := hazardproviders.Init(file)
		if err != nil {
			log.Fatal("Failed to get hazard provider for file: ", file, "\n", err)
		}
		hazardProviders[i] = hp
	}

	//create a result writer based on the name of the depth grid.
	//write local
	path := fmt.Sprintf("/workspaces/go-consequences/data/results/%s/%s/%v_consequences_summary.csv", year, scenario, dataset)
	w, _ := resultswriters.Init_csvSummaryResultsWriterFromFile(path)
	defer w.Close()
	//compute consequences.
	compute.StreamAbstract_MultiFreq_MultiVar(hazardProviders, frequencies, nsp, w)

}

// func main_old() {
// 	start := time.Now()
// 	fp := os.Args[1]
// 	b, err := os.ReadFile(fp)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	var config compute.Config
// 	json.Unmarshal(b, &config)
// 	computable, err := config.CreateComputable()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer computable.ResultsWriter.Close()
// 	defer computable.HazardProvider.Close()
// 	err = computable.Compute()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	elapsed := time.Since(start)
// 	fmt.Println("Execution time:", elapsed)
// }

// func main2() {
// 	//initialize the NSI API structure provider
// 	// nsp := structureprovider.InitNSISP()
// 	nsp, _ := structureprovider.InitStructureProvider("/workspaces/go-consequences/data/burlington-davenport-nsi.gpkg", "nsi", "GPKG")
// 	nsp.SetDeterministic(true)
// 	now := time.Now()
// 	fmt.Println(now)
// 	//nsp.SetDeterministic(true)
// 	//identify the depth grid to apply to the structures.
// 	root := "/workspaces/go-consequences/data/burlington-davenport-100yr"
// 	filepath := root + ".tif"
// 	w, _ := resultswriters.InitSpatialResultsWriter(root+"_consequencesGHG.gpkg", "results", "GPKG")
// 	//w := consequences.InitSummaryResultsWriterFromFile(root + "_consequences_SUMMARY.json")
// 	//create a result writer based on the name of the depth grid.
// 	//w, _ := resultswriters.InitGpkResultsWriter(root+"_consequences_nsi.gpkg", "nsi_result")
// 	defer w.Close()
// 	//initialize a hazard provider based on the depth grid.
// 	dfr, _ := hazardproviders.Init_CustomFunction(filepath, func(valueIn hazards.HazardData, hazard hazards.HazardEvent) (hazards.HazardEvent, error) {
// 		if valueIn.Depth == 0 {
// 			return hazard, hazardproviders.NoHazardFoundError{}
// 		}
// 		process := hazardproviders.DepthHazardFunction()
// 		return process(valueIn, hazard)
// 	})
// 	//compute consequences.
// 	fmt.Println("running compute.StreamAbstract")
// 	compute.StreamAbstract(dfr, nsp, w)
// 	// compute.StreamAbstractMultiVariate(dfr, nsp, w)
// 	fmt.Println(time.Since(now))
// }
