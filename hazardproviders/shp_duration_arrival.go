package hazardproviders

import (
	"errors"
	"time"

	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazards"
	"github.com/dewberry/gdal"
)

type shpDurationAndArrivalHazardProvider struct {
	arrival  time.Time
	duration float64
	extentCR cogReader
	bbox     geography.BBox
	geom     gdal.Geometry
	process  HazardFunction
}

func InitShpDAHP(gpkg string, layername string, extentfp string, arrivalDate time.Time, durationDays float64) (shpDurationAndArrivalHazardProvider, error) {
	d := gdal.OGRDriverByName("GPKG")
	ds, dsok := d.Open(gpkg, int(gdal.ReadOnly))
	if !dsok {
		return shpDurationAndArrivalHazardProvider{}, errors.New("Could not open geopackage " + gpkg)
	}

	ex, err := initCR(extentfp)
	if err != nil {
		panic(err)
	}

	layer := ds.LayerByName(layername)
	geom := layer.NextFeature().Geometry()
	env := geom.Envelope()
	xmin := env.MinX()
	xmax := env.MaxX()
	ymin := env.MinY()
	ymax := env.MaxY()
	bb := geography.BBox{Bbox: []float64{xmin, ymin, xmax, ymax}}

	ret := shpDurationAndArrivalHazardProvider{
		arrival:  arrivalDate,
		duration: durationDays,
		extentCR: ex,
		bbox:     bb,
		geom:     geom,
		process:  ArrivalAndDurationHazardFunction(),
	}

	return ret, nil
}

func (shp shpDurationAndArrivalHazardProvider) Close() {
	shp.extentCR.Close()
}

func (shp shpDurationAndArrivalHazardProvider) Hazard(l geography.Location) (hazards.HazardEvent, error) {
	var h hazards.HazardEvent
	hd := hazards.HazardData{
		Depth:       0,
		Velocity:    0,
		ArrivalTime: shp.arrival,
		Erosion:     0,
		Duration:    shp.duration,
		WaveHeight:  0,
		Salinity:    false,
		Qualitative: "",
	}
	he, err := shp.process(hd, h)
	if err != nil {
		panic(err)
	}

	check_geom := gdal.Create(gdal.GT_Point)
	check_geom.AddPoint(l.X, l.Y, 0)

	d, err := shp.extentCR.ProvideValue(l)
	if err != nil {
		return he, err
	}

	if shp.geom.Contains(check_geom) && (d > 0) {
		return he, nil
	}
	return he, errors.New("Provided hazard location is outside the provided boundary")

}

func (shp shpDurationAndArrivalHazardProvider) HazardBoundary() (geography.BBox, error) {
	return shp.bbox, nil
}
