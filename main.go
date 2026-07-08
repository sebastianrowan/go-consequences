package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/USACE/go-consequences/compute"
	"github.com/USACE/go-consequences/cropprovider"
	"github.com/USACE/go-consequences/hazardproviders"
	"github.com/USACE/go-consequences/resultswriters"
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
func main_conf() {
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

func main() {
	filter := make([]string, 11)
	filter[0] = "1" //filter to corn
	filter[1] = "5"
	filter[2] = "6"
	filter[3] = "22"
	filter[4] = "23"
	filter[5] = "24"
	filter[6] = "28"
	filter[7] = "36"
	filter[8] = "42"
	filter[9] = "52"
	filter[10] = "21"
	// nassSp := InitNassCropProvider("2024", filter)
	scenarios := []string{
		"HarrisonvilleRemoval",
		"ColumbiaRemoval",
		"PrairieDuPontRemoval",
		"MESDRemoval",
		"WoodRemoval",
		"1000mSetback",
		"3000mSetback",
		// "1000mSetbackLeveed",
		// "3000mSetbackLeveed",
	}

	rps := []int{2, 5, 10, 20, 50, 100, 200, 500}
	gpkg := "/workspaces/go-consequences/data/levee/leveed_areas.gpkg"
	at := time.Date(2026, time.Month(7), 1, 0, 0, 0, 0, time.UTC)

	now := time.Now()

	for _, scenario := range scenarios {
		cropFile := fmt.Sprintf("/workspaces/go-consequences/data/levee/%s-30m-CDLS.tif", scenario)
		for _, rp := range rps {
			now2 := time.Now()
			floodMap := fmt.Sprintf("/mnt/dlevee/data/depth_v3/ssp245-2020-2059/clipped/SSP245_2020-2059_%s_MaxDepth_%vyr.tif", scenario, rp)
			resFile := fmt.Sprintf("/mnt/dlevee/results/crops/ssp245-2020-2059/cropLoss-%s-%vyr.parquet", scenario, rp)

			crops := cropprovider.InitTiffCropProvider(cropFile, filter)
			fmt.Println(crops.TifReader.FilePath)
			hp, _ := hazardproviders.InitShpDAHP(gpkg, scenario, floodMap, at, 14.0)

			rw, _ := resultswriters.InitSpatialResultsWriter_EPSG_Projected(resFile, "agdamage", "Parquet", 5070) // testing data output
			compute.StreamAbstract(hp, crops, rw)
			hp.Close()
			rw.Close()
			fmt.Printf("%s %vyr finished in: %v\n", scenario, rp, time.Since(now2))

		}
	}
	fmt.Printf("All compute finished in: %v\n", time.Since(now))
}
