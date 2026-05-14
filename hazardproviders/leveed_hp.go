package hazardproviders

import (
	"errors"
	"time"

	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazards"
	"github.com/dewberry/gdal"
)

type leveedHazardProvider struct {
	depthcr cogReader
	bbox    geography.BBox
	shp     gdal.Geometry
	process HazardFunction
}

type LeveeInfo struct {
	DepthFP        string
	BoundaryDriver string
	BoundaryFP     string
	BoundaryLayer  string
	BoundaryName   string
}

func InitLeveedHP(li LeveeInfo) (leveedHazardProvider, error) {

	// Process boundary Layer
	d := gdal.OGRDriverByName(li.BoundaryDriver)
	ds, dsok := d.Open(li.BoundaryFP, int(gdal.ReadOnly))
	if !dsok {
		return leveedHazardProvider{}, errors.New("Error opening boundary layer " + li.BoundaryLayer)
	}

	// var layer gdal.Layer
	layer := ds.LayerByName(li.BoundaryLayer)
	var bb geography.BBox
	var geom gdal.Geometry

	fc, _ := layer.FeatureCount(true)
	for range fc {
		feat := layer.NextFeature()

		leveed_area := feat.FieldAsString(feat.Definition().FieldIndex("leveed_area"))
		if leveed_area == li.BoundaryName {
			geom = feat.Geometry()
			env := geom.Envelope()
			xmin := env.MinX()
			xmax := env.MaxX()
			ymin := env.MinY()
			ymax := env.MaxY()
			bb = geography.BBox{Bbox: []float64{xmin, ymin, xmax, ymax}}
		}
	}
	dr, ed := initCR(li.DepthFP)
	return leveedHazardProvider{depthcr: dr, bbox: bb, shp: geom, process: DepthHazardFunction()}, ed

}

func (lhp leveedHazardProvider) Close() {
	lhp.depthcr.Close()
}

func (lhp leveedHazardProvider) Hazard(l geography.Location) (hazards.HazardEvent, error) {
	var h hazards.HazardEvent
	test_geom := gdal.Create(gdal.GT_Point)
	test_geom.AddPoint(l.X, l.Y, 0)

	if lhp.shp.Contains(test_geom) {
		d, err := lhp.depthcr.ProvideValue(l)
		if err != nil {
			return h, err
		}
		hd := hazards.HazardData{
			Depth:       d,
			Velocity:    0,
			ArrivalTime: time.Time{},
			Erosion:     0,
			Duration:    0,
			WaveHeight:  0,
			Salinity:    false,
			Qualitative: "",
		}
		return lhp.process(hd, h)
	}
	return h, errors.New("Location outside of leveed area")
}

func (lhp leveedHazardProvider) HazardBoundary() (geography.BBox, error) {
	return lhp.bbox, nil
}
