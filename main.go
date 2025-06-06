package main

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"

	// _ "net/http/pprof"
	"os"
	"sync"
	"time"

	"github.com/USACE/go-consequences/compute"
	"github.com/USACE/go-consequences/hazardproviders"
	"github.com/USACE/go-consequences/resultswriters"
	"github.com/USACE/go-consequences/structureprovider"
)

func compute_FathomMultiFrequency(filename string, year string, scenario string) {

	// year := "2020"
	// year := 2050-SSP5_8.5
	// scenario := "FLUVIAL-DEFENDED_KNOWN"
	// fmt.Println("Active Goroutines:", runtime.NumGoroutine())

	dataset := filename[:len(filename)-4]

	//initialize the NSI API structure provider
	nsp, err := structureprovider.InitStructureProvider("/workspaces/go-consequences/data/nsi/nsi_2022.gpkg", "nsi", "GPKG")
	if err != nil {
		log.Fatal(err)
	}

	result_dir := fmt.Sprintf("/workspaces/go-consequences/data/results/%s/%s", year, scenario)
	result_file := fmt.Sprintf("%v_consequences.parquet", dataset)
	path := fmt.Sprintf("%s/%s", result_dir, result_file)

	// w, _ := resultswriters.Init_csvSummaryResultsWriterFromFile(path, "occupancy type")
	w, _ := resultswriters.InitSpatialResultsWriter(path, "result", "Parquet")
	// w := resultswriters.InitJsonResultsWriterFromFile(path)
	// defer w.Close()

	//initialize a set of frequencies
	rps := []int{5, 10, 20, 50, 100, 200, 500, 1000}
	frequencies := []float64{.2, .1, .05, .02, .01, .005, .002, .001}

	root := "/workspaces/go-consequences/data/fathom"
	//identify the depth grids to represent the frequencies.
	hazardProviders := make([]hazardproviders.HazardProvider, len(rps))

	for i, r := range rps {
		file := fmt.Sprintf("%s/%s/FLOOD_MAP-1_3ARCSEC-NW_OFFSET-1in%d-%s-DEPTH-%s-PERCENTILE50-v3.1/%s.tif", root, year, r, scenario, year, dataset)
		// fmt.Println(file)
		hp, err := hazardproviders.Init(file)
		if err != nil {
			log.Fatal("Failed to get hazard provider for file: ", file, "\n", err)
		}
		hazardProviders[i] = hp
	}

	//compute consequences.
	compute.StreamAbstract_MultiFreq_MultiVar(hazardProviders, frequencies, nsp, w)
	w.Close()
	w = nil
	runtime.GC()
}

func get_files(file_list []string) <-chan string {
	// func gen() in go.dev example

	out := make(chan string, len(file_list))

	for _, f := range file_list {
		out <- f
	}
	close(out)
	return out
}

func process_file2(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for i := range in {
			out <- i
		}
	}()
	return out
}

func merge2(year string, scenario string, cs ...<-chan string) <-chan string {
	var wg sync.WaitGroup
	out := make(chan string)

	output := func(c <-chan string) {
		for filename := range c {
			ts := time.Now()
			compute_FathomMultiFrequency(filename, year, scenario)
			te := time.Since(ts)
			out_str := fmt.Sprintf("Processed file: %s in %s", filename, te)
			out <- out_str
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

func run_with_channels(year string, scenario string, filelist string) {

	content, err := os.ReadFile(filelist)
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

	// could I make a list of chans and then specify N instead of manually defining each
	// var chans []chan string
	// N := 12
	// for i := 0; i < N; i++ {
	// 	chans = append(chans, process_file2(c))
	// }
	// for i := range merge2(year, scenario, chans) {
	// 		fmt.Println(i)
	//}
	c1 := process_file2(c)
	c2 := process_file2(c)
	c3 := process_file2(c)
	c4 := process_file2(c)
	c5 := process_file2(c)
	c6 := process_file2(c)
	c7 := process_file2(c)
	c8 := process_file2(c)
	c9 := process_file2(c)
	c10 := process_file2(c)
	c11 := process_file2(c)
	c12 := process_file2(c)

	for i := range merge2(year, scenario, c1, c2, c3, c4, c5, c6, c7, c8, c9, c10, c11, c12) {
		// for i := range merge(c1, c2, c3, c4, c5, c6) {
		fmt.Println(i)
	}

	te := time.Since(ts)
	fmt.Printf("All files completed in %s\n", te)
}

func run_with_wgs(year string, scenario string, filelist string) {

	content, err := os.ReadFile(filelist)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var file_list []string
	err = json.Unmarshal(content, &file_list)
	if err != nil {
		log.Fatal("Error during Unmarshal():", err)
	}

	i := 0
	max := 12

	for i < len(file_list) {
		var wg sync.WaitGroup
		for j := 0; j < max; j++ {
			if j == 0 {
				if i+max > len(file_list) {
					wg.Add(len(file_list) - i - 1)
				} else {
					wg.Add(max)
				}
			}
			if i < len(file_list) {
				go func(file string) {
					defer wg.Done()
					compute_FathomMultiFrequency(file, year, scenario)
				}(file_list[i])
				i++
			}
		}
		wg.Wait()
	}
}

type fathomConfig struct {
	Year     string `json:"year"`
	Scenario string `json:"scenario"`
	FileList string `json:"filelist"`
}

func main() {

	fp := os.Args[1]
	b, err := os.ReadFile(fp)
	if err != nil {
		log.Fatal(err)
	}
	var conf fathomConfig
	json.Unmarshal(b, &conf)

	// YEAR := "2100-SSP2_4.5"
	// // SSP := "SSP1_2.6"
	// SCENARIO := "FLUVIAL-UNDEFENDED"
	// FILELIST := "/workspaces/go-consequences/data/fathom/2020/files.json"

	// run_with_channels(YEAR, SCENARIO, FILELIST)
	// run_with_wgs(YEAR, SCENARIO, FILELIST)
	fmt.Println(conf.Year, conf.Scenario, conf.FileList)
}
