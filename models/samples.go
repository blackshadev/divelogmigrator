package models

type TankPressure struct {
	Tank     uint8
	Pressure float64
}
type Event struct {
	Type  string
	Value uint
	Flag  uint
}

type DiveSample struct {
	Time        Duration
	Pressure    []TankPressure
	Depth       *float64
	Temperature *float64
	RBT         *uint
	Events      []Event
}
