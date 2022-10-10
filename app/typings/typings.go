package typings

const GATE_TYPE_XOR = 0
const GATE_TYPE_ADD = 1
const GATE_TYPE_INV = 2

type Gate struct {
	Id         uint32
	Type       uint8
	InputWires []uint32
	OutputWire uint32
}

type Circuit struct {
	WireCount    int
	VInputSize   int
	PInputSize   int
	OutputSize   int
	AndGateCount int
	Gates        []Gate
	OutputsSizes []int
}
