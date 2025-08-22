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
	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazardproviders"
	"github.com/USACE/go-consequences/resultswriters"
	"github.com/USACE/go-consequences/structureprovider"
)

type fathomConfig struct {
	Year      string `json:"year"`
	SSP       string `json:"ssp"`
	Scenario  string `json:"scenario"`
	FileList  string `json:"filelist"`
	DataDir   string `json:"data_dir"`
	ResultDir string `json:"results_dir"`
}

func compute_FathomMultiFrequency(filename string, conf fathomConfig) {

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

	year_ssp := conf.Year
	if conf.SSP != "" {
		year_ssp = fmt.Sprintf("%s/%s", conf.Year, conf.SSP)
	}

	result_dir := fmt.Sprintf("%s/%s/%s", conf.ResultDir, year_ssp, conf.Scenario)

	err2 := os.MkdirAll(result_dir, os.ModePerm)
	if err2 != nil {
		log.Fatal(err2)
	}

	result_file := fmt.Sprintf("%v_consequences.parquet", dataset)
	path := fmt.Sprintf("%s/%s", result_dir, result_file)

	// w, _ := resultswriters.Init_csvSummaryResultsWriterFromFile(path, "occupancy type")
	w, _ := resultswriters.InitSpatialResultsWriter(path, "result", "Parquet")
	// w := resultswriters.InitJsonResultsWriterFromFile(path)
	// defer w.Close()

	//initialize a set of frequencies
	rps := []int{5, 10, 20, 50, 100, 200, 500, 1000}
	frequencies := []float64{.2, .1, .05, .02, .01, .005, .002, .001}

	// root := "/workspaces/go-consequences/data/fathom"
	//identify the depth grids to represent the frequencies.
	hazardProviders := make([]hazardproviders.HazardProvider, len(rps))

	if conf.SSP != "" {
		year_ssp = fmt.Sprintf("%s-%s", conf.Year, conf.SSP)
	}

	for i, r := range rps {
		file := fmt.Sprintf("%s/%s/FLOOD_MAP-1_3ARCSEC-NW_OFFSET-1in%d-%s-DEPTH-%s-PERCENTILE50-v3.1/%s.tif", conf.DataDir, conf.Year, r, conf.Scenario, year_ssp, dataset)
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

func merge2(conf fathomConfig, cs ...<-chan string) <-chan string {
	var wg sync.WaitGroup
	out := make(chan string)

	output := func(c <-chan string) {
		for filename := range c {
			ts := time.Now()
			compute_FathomMultiFrequency(filename, conf)
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

func run_with_channels(conf fathomConfig) {

	content, err := os.ReadFile(conf.FileList)
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
	// N := 15
	// chans := make([]<-chan string, N)
	// for i := 0; i < N; i++ {
	// 	ci := process_file2(c)
	// 	chans = append(chans, ci)
	// }
	// for i := range merge2(conf, chans...) {
	// 	fmt.Println(i)
	// }

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
	c13 := process_file2(c)
	c14 := process_file2(c)
	c15 := process_file2(c)
	c16 := process_file2(c)

	for i := range merge2(conf, c1, c2, c3, c4, c5, c6, c7, c8, c9, c10, c11, c12, c13, c14, c15, c16) {
		fmt.Println(i)
	}

	te := time.Since(ts)
	fmt.Printf("All files completed in %s\n", te)
}

func run_with_wgs(conf fathomConfig) {

	content, err := os.ReadFile(conf.FileList)
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
					compute_FathomMultiFrequency(file, conf)
				}(file_list[i])
				i++
			}
		}
		wg.Wait()
	}
}

func mainjunk() {

	// when running with data from external hard drive, analysis took 5.5 hours vs 1.5 with data on internal solid state drive
	fp := os.Args[1]
	b, err := os.ReadFile(fp)
	if err != nil {
		log.Fatal(err)
	}
	var conf fathomConfig
	json.Unmarshal(b, &conf)

	run_with_channels(conf)
}

func main() {

	fmt.Println("Begin")
	nsp, err := structureprovider.InitStructureProvider("/workspaces/go-consequences/data/nsi/nsi_2022.gpkg", "nsi", "GPKG")
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now()
	fmt.Println(now)
	//nsp.SetDeterministic(true)
	//identify the depth grid to apply to the structures.
	// filepath := "/workspaces/go-consequences/data/testraster2.tif"
	// w, _ := resultswriters.InitSpatialResultsWriter("/workspaces/go-consequences/data/test3.gpkg", "results", "GPKG")
	//w := consequences.InitSummaryResultsWriterFromFile(root + "_consequences_SUMMARY.json")
	//create a result writer based on the name of the depth grid.
	//w, _ := resultswriters.InitGpkResultsWriter(root+"_consequences_nsi.gpkg", "nsi_result")
	// defer w.Close()
	//initialize a hazard provider based on the depth grid.

	//compute consequences.
	fmt.Println("Starting BBOX stream...")

	nsp.ByBbox(geography.BBox{Bbox: []float64{-72.0, 45.0, -72.1, 45.1}}, func(f consequences.Receptor) {
		fmt.Println(f)
	})

	fmt.Println(time.Since(now))
}
