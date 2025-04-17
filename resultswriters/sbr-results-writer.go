package resultswriters

import (
	"fmt"
	"io"
	"math"
	"os"
	"regexp"

	"github.com/USACE/go-consequences/consequences"
)

// Results should be produced as a DataFrame which has the sum of each impact category for each occtype
//		In R with tidyverse syntax, the code below would produce the desired result if starting with a
//		file produced from a spatialResultsWriter
//
// | occtype |

type csvSummaryResultsWriter struct {
	filepath    string
	w           io.Writer
	Colnames    []string
	hasColnames bool
	Rows        map[string]map[string]float64
}

func Init_csvSummaryResultsWriterFromFile(filepath string) (*csvSummaryResultsWriter, error) {
	w, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return &csvSummaryResultsWriter{}, err
	}
	c := []string{}
	r := make(map[string]map[string]float64)
	return &csvSummaryResultsWriter{filepath: filepath, w: w, Colnames: c, hasColnames: false, Rows: r}, nil
}

func (csrw *csvSummaryResultsWriter) Close() {
	csrw.calc_moe()
	fmt.Fprint(csrw.w, "occtype")
	for _, c := range csrw.Colnames {
		fmt.Fprintf(csrw.w, ",%s", c)
	}

	for occtype, row_vals := range csrw.Rows {
		fmt.Fprintf(csrw.w, "\n%s", occtype)
		for _, c := range csrw.Colnames {
			fmt.Fprintf(csrw.w, ",%f", row_vals[c])
		}
	}

	w2, ok := csrw.w.(io.WriteCloser)
	if ok {
		w2.Close()
	}
}

func (csrw *csvSummaryResultsWriter) calc_moe() {
	// Find standard deviation columns and take their square root
	//	values added before this is called are the sums of variance (SD^2)
	rp_sd_pattern := regexp.MustCompile(`^\d+.*S$`)
	ead_sd_pattern := regexp.MustCompile(`.S_EAD$`)
	var occtypes = []string{}
	var sd_cols = []string{}

	for _, colname := range csrw.Colnames {
		sd_match := rp_sd_pattern.MatchString(colname) || ead_sd_pattern.MatchString(colname)
		if sd_match {
			sd_cols = append(sd_cols, colname)
		}
	}
	for _, o := range occtypes {
		for _, col := range sd_cols {
			// this needs to happen for each occtype. The data is nested, not tabular right now
			moe := csrw.Rows[o][col]
			sd := math.Sqrt(moe)
			csrw.Rows[o][col] = sd

		}
	}

}

func (csrw *csvSummaryResultsWriter) Write(r consequences.Result) {
	occtype := "occtype not found somehow"
	// var resmap = make(map[string]interface{})
	cols := make(map[string]float64)

	rp_mean_pattern := regexp.MustCompile(`^\d+.*M$`)
	ead_mean_pattern := regexp.MustCompile(`M_EAD$`)

	rp_sd_pattern := regexp.MustCompile(`^\d+.+\D+S$`)
	ead_sd_pattern := regexp.MustCompile(`.S_EAD$`)

	for i, h := range r.Headers {

		if h == "occupancy type" {
			occtype = r.Result[i].(string)
		}
		if h == "structure damage" {
			cols[h] = r.Result[i].(float64)
			if !csrw.hasColnames {
				csrw.Colnames = append(csrw.Colnames, h)
			}
		}
		if h == "content damage" {
			cols[h] = r.Result[i].(float64)
			if !csrw.hasColnames {
				csrw.Colnames = append(csrw.Colnames, h)
			}
		}
		if h == "S_EAD" {
			cols[h] = r.Result[i].(float64)
			if !csrw.hasColnames {
				csrw.Colnames = append(csrw.Colnames, h)
			}
		}
		if h == "C_EAD" {
			cols[h] = r.Result[i].(float64)
			if !csrw.hasColnames {
				csrw.Colnames = append(csrw.Colnames, h)
			}
		}
		rp_mean_match := rp_mean_pattern.MatchString(h)
		ead_mean_match := ead_mean_pattern.MatchString(h)
		rp_sd_match := rp_sd_pattern.MatchString(h)
		ead_sd_match := ead_sd_pattern.MatchString(h)
		if rp_mean_match {
			cols[h] = r.Result[i].(float64)
			if !csrw.hasColnames {
				csrw.Colnames = append(csrw.Colnames, h)
			}
		}
		if rp_sd_match {
			sd := r.Result[i].(float64)
			moe := 1.96 * sd / math.Sqrt(500) // 95% CI, used N=500 because underlying simulations included 500 MCS runs per building per depth
			cols[h] = moe                     // Add MOE in each iteration and take sqrt at the end
			if !csrw.hasColnames {
				csrw.Colnames = append(csrw.Colnames, h)
			}
		}
		if ead_mean_match {
			cols[h] = r.Result[i].(float64)
			if !csrw.hasColnames {
				csrw.Colnames = append(csrw.Colnames, h)
			}
		}
		if ead_sd_match {
			sd := r.Result[i].(float64)
			moe := 1.96 * sd / math.Sqrt(500) // 95% CI, used N=500 because underlying simulations included 500 MCS runs per building per depth
			cols[h] = moe                     // Add MOE in each iteration and take sqrt at the end
			if !csrw.hasColnames {
				csrw.Colnames = append(csrw.Colnames, h)
			}
		}

	}
	if !csrw.hasColnames {
		csrw.hasColnames = true
	}
	existing_cols, ok := csrw.Rows[occtype]
	if ok {
		// results have already been added for this occtype.
		// need to add current values to existing
		for i, v := range cols {
			existing_cols[i] += v
		}
	} else {
		csrw.Rows[occtype] = cols

	}
}

type MultiFrequencyResultSetWriter struct {
}
