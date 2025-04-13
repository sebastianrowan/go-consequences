package resultswriters

import (
	"io"
	"os"

	"github.com/USACE/go-consequences/consequences"
)

// Results should be produced as a DataFrame which has the sum of each impact category for each occtype
//		In R with tidyverse syntax, the code below would produce the desired result if starting with a
//		file produced from a spatialResultsWriter
//
// | occtype |

type csvSummaryResultsWriter struct {
	filepath string
	w        io.Writer
	totals   map[string]float64
}

func Init_csvSummaryResultsWriterFromFile(filepath string) (*csvSummaryResultsWriter, error) {
	w, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return &csvSummaryResultsWriter{}, err
	}

	totals := make(map[string]float64, 1)

	return &csvSummaryResultsWriter{filepath: filepath, w: w, totals: totals}, nil
}
func Init_csvSummaryResultsWriter(w io.Writer) *csvSummaryResultsWriter {
	totals := make(map[string]float64, 1)

	return &csvSummaryResultsWriter{filepath: "not applicable", w: w, totals: totals}
}

func (csrw *csvSummaryResultsWriter) Write(r consequences.Result) {
	var occtype = ""
	// Is there a way to get all value "columns" using a pattern (regex?) rather than hardcoding each one?
	//	I can write this function to work for the specific return periods and frequencies of my analysis, but 
	//	ideally it would be flexible to different frequencies (or no frequencies) and different impact calculations
	//		for header in
	structure_damage := 0.0
	content_damage := 0.0
	h = r.Headers
	for i, v := range h {
		if v == "occupancy type" {
			occtype = r.Result[i].(string)
		}
		if 
	}
	var struct_dam_ead = "structure damage"
	var cont_dam_ead =  "content damage"
	var dmg_mean = "dmg_mean" // this is the mean value from the multivariate damage calculation
	var dmg_sd = "dmg_sd" // this is the standard deviation from the multivariate damage calculation
	var ghg_mean = "ghg_mean" // this is the mean value from the multivariate GHG calculation
	var ghg_sd =  "ghg_sd" // this is the standard deviation value from the multivariate GHG calculation

	// Will calculate low and high (95% CI) for each impact variable 
	// 		>>> CI = mean_i * 1.96 * sd_i
	// Then will have running total for mean, low, and high
}


type MultiFrequencyResultSetWriter struct{
	
}