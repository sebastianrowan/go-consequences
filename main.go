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

type leveeConfig struct {
	Year         string `json:"year"`
	SSP          string `json:"ssp"`
	ScenarioList string `json:"scenarios"`
	DataDir      string `json:"data_dir"`
	ResultDir    string `json:"results_dir"`
}

func compute_LeveeMultiFrequency(scenario string, conf leveeConfig) {

	nsp, err := structureprovider.InitStructureProvider("/workspaces/go-consequences/data/nsi/nsi_levee_setback.gpkg", "nsi_sp", "GPKG")
	if err != nil {
		log.Fatal(err)
	}

	year_ssp := conf.Year
	if conf.SSP != "" {
		year_ssp = fmt.Sprintf("%s-%s", conf.Year, conf.SSP)
	}

	results_filename := fmt.Sprintf("%s_%s_consequences.parquet", year_ssp, scenario)
	results_path := fmt.Sprintf("%s/%s", conf.ResultDir, results_filename)

	w, _ := resultswriters.InitSpatialResultsWriter(results_path, "result", "Parquet")

	rps := []int{1, 2, 5, 10, 20, 50, 100, 200, 500}
	frequencies := make([]float64, len(rps))
	hazardProviders := make([]hazardproviders.HazardProvider, len(rps))

	for i, r := range rps {
		frequencies[i] = 1.0 / float64(r)
		file := fmt.Sprintf("%s/%dyr_%s_%s_Depth.tif", conf.DataDir, r, year_ssp, scenario)

		hp, err := hazardproviders.Init(file)
		if err != nil {
			log.Fatal("Failed to get hazard provider for file: ", file, "\n", err)
		}
		hazardProviders[i] = hp
	}

	compute.StreamAbstract_MultiFreq_MultiVar(hazardProviders, frequencies, nsp, w)
	w.Close()

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

func merge2(conf leveeConfig, cs ...<-chan string) <-chan string {
	var wg sync.WaitGroup
	out := make(chan string)

	output := func(c <-chan string) {
		for filename := range c {
			ts := time.Now()
			compute_LeveeMultiFrequency(filename, conf)
			// compute_FathomMultiFrequency(filename, conf)
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

func run_with_channels(conf leveeConfig) {

	content, err := os.ReadFile(conf.ScenarioList)
	// content, err := os.ReadFile(conf.FileList)
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
	// c6 := process_file2(c)
	// c7 := process_file2(c)
	// c8 := process_file2(c)
	// c9 := process_file2(c)
	// c10 := process_file2(c)
	// c11 := process_file2(c)
	// c12 := process_file2(c)
	// c13 := process_file2(c)
	// c14 := process_file2(c)
	// c15 := process_file2(c)
	// c16 := process_file2(c)

	// , c2, c3, c4, c5, c6, c7, c8, c9, c10, c11, c12, c13, c14, c15, c16
	for i := range merge2(conf, c1, c2, c3, c4, c5) {
		fmt.Println(i)
	}

	te := time.Since(ts)
	fmt.Printf("All files completed in %s\n", te)
}

func main() {

	// when running with data from external hard drive, analysis took 5.5 hours vs 1.5 with data on internal solid state drive
	// fp := os.Args[1]
	fp := "config_levee.json"
	b, err := os.ReadFile(fp)
	if err != nil {
		log.Fatal(err)
	}
	var conf leveeConfig
	json.Unmarshal(b, &conf)
	run_with_channels(conf)
}
