package hazards

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
	spectral_displacement       float64
	spectral_acceleration       float64
	eq_peak_ground_acceleration float64
}
