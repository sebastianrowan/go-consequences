package resultswriters

import (
	"fmt"
	"io"
	"os"

	"github.com/USACE/go-consequences/consequences"
)

type jsonResultsWriter struct {
	filepath             string
	w                    io.Writer
	headerHasBeenWritten bool
}

func InitJsonResultsWriterFromFile(filepath string) *jsonResultsWriter {
	w, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}
	return &jsonResultsWriter{filepath: filepath, w: w}
}
func InitJsonResultsWriter(w io.Writer) *jsonResultsWriter {
	return &jsonResultsWriter{filepath: "not applicapble", w: w}
}

func (srw *jsonResultsWriter) Write(r consequences.Result) {
	if !srw.headerHasBeenWritten {
		fmt.Fprintf(srw.w, "{\"consequences\":[")
		srw.headerHasBeenWritten = true
	}
	b, _ := r.MarshalJSON()
	s := string(b) + ","
	fmt.Fprint(srw.w, s)
}

func (srw *jsonResultsWriter) Close() {
	fmt.Fprintf(srw.w, "]}")
	w2, ok := srw.w.(io.WriteCloser)
	if ok {
		w2.Close()
	}
}

// func (srw *jsonResultsWriter) Write(r consequences.Result) {
// 	if !srw.headerHasBeenWritten {
// 		srw.r = "{\"consequences\":["
// 		srw.headerHasBeenWritten = true
// 	}
// 	b, _ := r.MarshalJSON()
// 	s := string(b) + ","
// 	srw.r = srw.r + s
// }
// func (srw *jsonResultsWriter) Close() {
// 	srw.r = srw.r + "]}"
// 	fmt.Fprintf(srw.w, srw.r)
// 	w2, ok := srw.w.(io.WriteCloser)
// 	if ok {
// 		w2.Close()
// 	}
// }
