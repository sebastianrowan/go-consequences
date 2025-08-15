package main

import "C"

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/USACE/go-consequences/compute"
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

//go build -o pyconsequences/pyconsequences.dll -buildmode=c-shared main.go
//go build -o pyconsequences/pyconsequences.so -buildmode=c-shared main.go

//export RunFromConfigFile
func RunFromConfigFile(fp *C.char) {
	//TODO: when users run go-consequences through a compiled dll with python,
	//	will we run into issues if users pass a relative file path?
	//	or will the dll run in the same working directory as the python
	//	code that imported it?

	config_path := C.GoString(fp)

	b, err := os.ReadFile(config_path)
	if err != nil {
		log.Fatal(err)
	}

	var conf compute.Config
	json.Unmarshal(b, &conf)

	fmt.Println(conf)
}
