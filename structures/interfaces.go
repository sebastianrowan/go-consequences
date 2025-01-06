package structures

import (
	"github.com/USACE/go-consequences/hazards"
)

// Common interface for all damage function providers
type IPrototype interface {
	DamageFunction(component string) DamageFunction // component = structure, content, vehicle, etc.
}

type OccupancyTypeProvider interface {
	OccupancyTypeMap() map[string]OccupancyTypeStochastic
	OccupancyTypeMapSBR() map[string]OccupancyTypeSBR
	Write(outputpath string) error
}

// Main container that includes all info associated with a structure prototype
type Prototype struct {
	Name                        string         `json:"Name"`
	Description                 string         `json:"Description"`
	DamageCategory              DamageCategory `json:"DamageCategory"`
	FoundationHeightUncertainty Uncertainty    `json:"FoundationHeightUncertainty"`
	StructureUncertainty        Uncertainty    `json:"StructureUncertainty"`
	ContentUncertainty          Uncertainty    `json:"ContentUncertainty"`
	OtherUncertainty            Uncertainty    `json:"OtherUncertainty"`
	VehicleUncertainty          Uncertainty    `json:"VehicleUncertainty"`
	StructureDD                 FunctionDD     `json:"StructureDD"`
	ContentDD                   FunctionDD     `json:"ContentDD"`
	OtherDD                     FunctionDD     `json:"OtherDD"`
	VehicleDD                   FunctionDD     `json:"VehicleDD"`
}

type RawDFStruct struct {
	OccTypes struct {
		Prototypes []Prototype `json:"OccupancyType"`
	} `json:"OccTypes"`
}

// Map with OccupancyType string as index
type DFStore map[string]Prototype

type DamageCategory struct {
	Name        string      `json:"Name"`
	Description interface{} `json:"Description"`
	Rebuild     string      `json:"Rebuild"`
	CostFactor  string      `json:"CostFactor"`
}

type Uncertainty struct {
	None struct {
		Value string `json:"_value"`
	} `json:"None"`
}

type Ordinate struct {
	X     float64 `json:"X,string"`
	Value float64 `json:"_value,string"`
}

type FunctionDD struct {
	CalculateDamage       bool `json:"CalculateDamage,string"`
	MonotonicCurveUSingle struct {
		UncertaintyType string     `json:"UncertaintyType"`
		Ordinates       []Ordinate `json:"Ordinate"`
	} `json:"MonotonicCurveUSingle"`
}

type DamageVector struct {
	Intercept    float64 `json:"intercept"`
	Depth        float64 `json:"depth"`
	Sqft         float64 `json:"sqft"`
	N_bed        float64 `json:"n_bed"`
	N_bath       float64 `json:"n_bath"`
	N_car        float64 `json:"n_car"`
	Depth_sqft   float64 `json:"depth_sqft"`
	Depth_n_bed  float64 `json:"depth_n_bed"`
	Depth_n_bath float64 `json:"depth_n_bath"`
	Depth_n_car  float64 `json:"depth_n_car"`
}

type DamageFunctionSBR struct {
	Source           string            `json:"source"`
	DamageDriver     hazards.Parameter `json:"damagedriver"`
	DamageVectorMean DamageVector      `json:"damagevectormean"`
	DamageVectorSD   DamageVector      `json:"damagevectorsd"`
}

type DamageFunctionFamilySBR struct {
	DamageFunctions map[string]DamageFunctionSBR `json:"damagefunctions"`
}

type OccupancyTypeSBR struct {
	Name                     string                             `json:"name"`
	ComponentDamageFunctions map[string]DamageFunctionFamilySBR `json:"componentdamagefunctions"`
}

type OccupancyTypesContainerSBR struct {
	OccupancyTypes map[string]OccupancyTypeSBR `json:"occupancytypes"`
}

//////////////////////////////////////
//  END DepthDFProvider
//////////////////////////////////////
