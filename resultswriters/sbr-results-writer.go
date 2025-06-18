package resultswriters

import (
	"fmt"
	"io"
	"math"
	"os"
	"regexp"

	"github.com/USACE/go-consequences/consequences"
)

const (
	NONE       string = "none"           // groupby none -> aggregate to whole grid
	DAMCAT     string = "damcat"         // RES, COM, IND, etc.
	OCCTYPE    string = "occupancy type" // RES1-1SNB, RES2, RES3A, etc
	COUNTY     string = "county"
	TRACT      string = "tract"      // tract fips code (substring of cb_id from nsi)
	BLOCKGROUP string = "blockgroup" // block group fips code (substring of cb_id from nsi)
	BLOCK      string = "block"      // cb_id
	GRID       string = "grid"       // aggregate by rounding structure x,y coords down to nearest whole degree
	GRID10     string = "grid10"     // aggregate by rounding structure x,y coords down to neareast 0.1 degree
)

// type csvResultsWriter struct {
// 	filepath    string
// 	w           csv.Writer
// 	hasColnames bool
// }

// func InitCSVResultsWriterFromFile(filepath string) *csvResultsWriter {
// 	w, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return &csvResultsWriter{filepath: filepath, w: w}
// }

// func (crw *csvResultsWriter) Write(r consequences.Result) {
// 	if !crw.hasColnames {
// 		crw.w.Write(r.Headers)
// 		crw.hasColnames = true
// 	}
// 	// crw.Write(r.Result)
// }

type hucMatchWriter struct {
	filepath string
	w        io.Writer
	tx       string
	index    int
}

func InitHucMatchWriter(filepath string) (*hucMatchWriter, error) {
	w, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return &hucMatchWriter{}, err
	}
	tx := ""
	return &hucMatchWriter{filepath: filepath, w: w, tx: tx, index: 0}, nil
}

func (hw *hucMatchWriter) Write(r consequences.Result) {
	if hw.index == 0 {
		hw.tx = "fid,huc08"
	}

	hw.index++
}

func (hw *hucMatchWriter) Commit() {
	fmt.Fprintf(hw.w, "%s", hw.tx)
	hw.tx = ""
}

type csvSummaryResultsWriter struct {
	filepath    string
	w           io.Writer
	groupby     string // aggvar could be a list of levels (e.g. group_by(tract, occtype))
	Colnames    []string
	hasColnames bool
	Rows        map[string]map[string]float64
}

func Init_csvSummaryResultsWriterFromFile(filepath string, groupby string) (*csvSummaryResultsWriter, error) {
	w, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return &csvSummaryResultsWriter{}, err
	}
	c := []string{}
	r := make(map[string]map[string]float64)
	return &csvSummaryResultsWriter{filepath: filepath, w: w, groupby: groupby, Colnames: c, hasColnames: false, Rows: r}, nil
}

func (csrw *csvSummaryResultsWriter) Close() {
	csrw.calc_moe()

	switch csrw.groupby {
	case OCCTYPE:
		fmt.Fprint(csrw.w, "occtype") // OCCTYPE const == "occupancy type", but want "occtype" as column name
	default: // add cases for grid options eventually
		fmt.Fprint(csrw.w, csrw.groupby)
	}

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
	csrw.Rows = make(map[string]map[string]float64)
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

	agg := "total"
	this_aggvar := ""
	use_fips := false
	fips_len := 15

	// occtype := "occtype not found somehow"
	// var resmap = make(map[string]interface{})
	cols := make(map[string]float64)

	rp_mean_pattern := regexp.MustCompile(`^\d+.*M$`)
	ead_mean_pattern := regexp.MustCompile(`M_EAD$`)

	rp_sd_pattern := regexp.MustCompile(`^\d+.+\D+S$`)
	ead_sd_pattern := regexp.MustCompile(`.S_EAD$`)

	// other idea for aggregating on fips:
	// if so, can we just find cb_id value and substr without iterating through all headers?
	switch csrw.groupby {
	case DAMCAT:
		use_fips = false
		this_aggvar = DAMCAT
	case OCCTYPE:
		use_fips = false
		this_aggvar = OCCTYPE
	case COUNTY:
		use_fips = true
		fips_len = 5
		this_aggvar = "cb_id"
	case TRACT:
		use_fips = true
		fips_len = 11
		this_aggvar = "cb_id"
	case BLOCKGROUP:
		use_fips = true
		fips_len = 12
		this_aggvar = "cb_id"
	case BLOCK:
		use_fips = true
		fips_len = 15
		this_aggvar = "cb_id"
	case GRID:
		panic("grid-level aggregation not available")
	case GRID10:
		panic("grid-level aggregation not available")
	case NONE:
		use_fips = false
	default:
		panic("no groupby variable supplied")
	}

	for i, h := range r.Headers {

		if h == this_aggvar {
			// fmt.Println("AGGREGATING on " + this_aggvar)
			if use_fips {
				cb_id := r.Result[i].(string)
				fmt.Println("cb_id: " + cb_id)
				agg = cb_id[:fips_len]
			} else {
				agg = r.Result[i].(string)
			}
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
	existing_cols, ok := csrw.Rows[agg]
	if ok {
		// results have already been added for this occtype.
		// need to add current values to existing
		for i, v := range cols {
			existing_cols[i] += v
		}
	} else {
		csrw.Rows[agg] = cols

	}
}

type MultiFrequencyResultSetWriter struct{}
