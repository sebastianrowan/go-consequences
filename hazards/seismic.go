package hazards

import (
	"fmt"
	"time"
)

// Notes from Hazus Earthquake Model Technical Manual:
// https://www.fema.gov/sites/default/files/2020-10/fema_hazus_earthquake_technical_manual_4-2.pdf
// https://www.fema.gov/sites/default/files/documents/fema_hazus-earthquake-model-technical-manual-6-1.pdf

// TODO: notes below are from version 4.2 of the Hazus Earthquake Model, update
// to version 6.1

// "Hazus earthquake building damage functions are in the form of
// lognormal fragility curves that relate the probability of being in,
// or exceeding, a damage state for a given Potential Earthquake Hazard (PEH)
// demand parameter (e.g. response spectrum displacement)."

// Section 5.4 - Building Damage Due to Ground Shaking
// Section 5.4.2 - Fragility Curves
// * Spectral displacement is the PEH parameter used for structural and
//   non-structural damage to "drift-sensitive" components.
// * Spectral acceleration is the PEH parameter used for non-structural
//   damage to "acceleration-sensitive" components.

// Section 5.4.3 - Structural Fragility Curves - Equivalent Peak Ground Acceleration
// * Equivalent Peak ground acceleration (PGA) is the PEH parameter used for structural
//   damage to buildings that are components of utility and transportation
//   systems

// Section 5.5 - Building Damage Due to Ground Failure
// Section 5.5.1 - Fragility Curves - Peak Ground Displacement
// * Separate fragility curves distinguish between ground failure due to
//   lateral spreading and ground failure due to ground settlement, and
//   between shallow and deep foundations.
// * By default, Hazus assumes all buildings are on shallow foundations.

// Section 5.5.
type SeismicEvent struct {
	peakGroundAcceleration     float64 // peak ground motion (acceleration, as percent of g)
	peakGroundVelocity         float64 // peak ground motion (velocity in cm/s)
	peakSpectralAcceleration03 float64 // spectral acceleration at 0.3 s period, 5% damping (percent of g)
	peakSpectralAcceleration10 float64 // spectral acceleration at 1.0 s period, 5% damping (percent of g)
	peakSpectralAcceleration30 float64 // spectral acceleration at 3.0 s period, 5% damping (percent of g)
	intensity                  float64 // estimated instrumental intensity
}

func (h SeismicEvent) Depth() float64 {
	return -901.0
}
func (h SeismicEvent) Velocity() float64 {
	return -901.0
}
func (h SeismicEvent) ArrivalTime() time.Time {
	return time.Time{}
}
func (h SeismicEvent) Erosion() float64 {
	return -901.0
}
func (h SeismicEvent) Duration() float64 {
	return -901.0
}
func (h SeismicEvent) WaveHeight() float64 {
	return -901.0
}
func (h SeismicEvent) Salinity() bool {
	return false
}
func (h SeismicEvent) Qualitative() string {
	return ""
}
func (h SeismicEvent) DV() float64 {
	return -901.0
}

// Parameters implements the HazardEvent interface
func (h SeismicEvent) Parameters() Parameter {
	dp := Default
	dp = SetHasDepth(dp)
	return dp
}

// Has implements the HazardEvent Interface
func (h SeismicEvent) Has(p Parameter) bool {
	dp := h.Parameters()
	return dp&p != 0
}

func (s SeismicEvent) MarshalJSON() ([]byte, error) {
	j := fmt.Sprintf("{\"seismicevent\":{\"pga\":%f, \"pgv\":%f}", s.PeakGroundAcceleration(), s.PeakGroundVelocity())
	// TODO: when functions are implemented, add peak spectral accelerations and intensity to json string
	return []byte(j), nil
}

func (s SeismicEvent) PeakGroundAcceleration() float64 {
	return s.peakGroundAcceleration
}

func (s *SeismicEvent) SetPeakGroundAcceleration(pga float64) {
	s.peakGroundAcceleration = pga
}

func (s SeismicEvent) PeakGroundVelocity() float64 {
	return s.peakGroundVelocity
}

func (s *SeismicEvent) SetPeakGroundVelocity(pgv float64) {
	s.peakGroundVelocity = pgv
}

func (s SeismicEvent) PeakSpectralAcceleration03() float64 {
	return s.peakSpectralAcceleration03
}

func (s *SeismicEvent) SetSpectralAcceleration03(psa03 float64) {
	s.peakSpectralAcceleration03 = psa03
}

func (s SeismicEvent) PeakSpectralAcceleration10() float64 {
	return s.peakSpectralAcceleration10
}

func (s *SeismicEvent) SetSpectralAcceleration10(psa10 float64) {
	s.peakSpectralAcceleration03 = psa10
}

func (s SeismicEvent) PeakSpectralAcceleration30() float64 {
	return s.peakSpectralAcceleration30
}

func (s *SeismicEvent) SetSpectralAcceleration30(psa30 float64) {
	s.peakSpectralAcceleration30 = psa30
}

func (s SeismicEvent) Intensity() float64 {
	return s.intensity
}

func (s *SeismicEvent) SetIntensity(i float64) {
	s.intensity = i
}
