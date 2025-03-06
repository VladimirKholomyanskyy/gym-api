package common

type WeightUnit string
type HeightUnit string
type MeasurementUnit string

const (
	WeightUnitKg        WeightUnit      = "kg"
	WeightUnitLbs       WeightUnit      = "lbs"
	HeightUnitMeters    HeightUnit      = "m"
	HeightUnitFeet      HeightUnit      = "ft"
	MeasurementMetric   MeasurementUnit = "metric"
	MeasurementImperial MeasurementUnit = "imperial"
)

func (w WeightUnit) IsValid() bool {
	return w == WeightUnitKg || w == WeightUnitLbs
}

func (h HeightUnit) IsValid() bool {
	return h == HeightUnitMeters || h == HeightUnitFeet
}

func (m MeasurementUnit) IsValid() bool {
	return m == MeasurementMetric || m == MeasurementImperial
}
