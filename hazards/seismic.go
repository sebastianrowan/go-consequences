package hazards

type SeismicEvent struct {
	magnitute float64
}

func (h SeismicEvent) Depth() float64 {
	return h.magnitute //
}
