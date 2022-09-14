package gc

import (
	"zkpass-node/typings"
	u "zkpass-node/utils"
)

type Evaluator struct {
	Steps        int
	Circuit      []*typings.Circuit
	GarbledTable [][]byte
}

func (e *Evaluator) Init(circuits []*typings.Circuit, steps int) {
	e.Steps = steps
	e.Circuit = circuits
	e.GarbledTable = make([][]byte, len(e.Circuit))
}

func (e *Evaluator) Evaluate(no int, vLabels, pLabels,
	garbledTable []byte) []byte {
	type chunk struct {
		wireLabels   *[][]byte
		garbledTable *[]byte
	}

	c := (e.Circuit)[no]
	vLabelsChunks := u.Slice(vLabels, c.VInputSize*16)
	pLabelsChunks := u.Slice(pLabels, c.PInputSize*16)
	garbledTableChunks := u.Slice(garbledTable, c.AndGateCount*48)

	steps := []int{0, 1, 1, 1, 1, 1, e.Steps, 1}[no]
	chunks := make([]chunk, steps)
	for r := 0; r < steps; r++ {
		wireLabels := make([][]byte, c.WireCount)
		copy(wireLabels, u.Slice(u.Concat(vLabelsChunks[r], pLabelsChunks[r]), 16))
		chunks[r] = chunk{&wireLabels, &garbledTableChunks[r]}
	}

	encodedOutput := make([][]byte, steps)
	for r := 0; r < steps; r++ {
		encodedOutput[r] = evaluate(c, chunks[r].wireLabels, chunks[r].garbledTable)
	}
	return u.Concat(encodedOutput...)
}

func evaluate(c *typings.Circuit, wireLabels *[][]byte, garbledTable *[]byte) []byte {
	offset := 0

	for i := 0; i < len(c.Gates); i++ {
		g := c.Gates[i]
		if g.Type == typings.GATE_TYPE_ADD {
			evaluateAndGate(g, wireLabels, garbledTable, offset)
			offset += 1
		} else if g.Type == typings.GATE_TYPE_XOR {
			evaluateXorGate(g, wireLabels)
		} else if g.Type == typings.GATE_TYPE_INV {
			evaluateInvGate(g, wireLabels)
		} else {
			panic("Unknown gate")
		}
	}

	lsb := make([]int, c.OutputSize)
	for i := 0; i < c.OutputSize; i++ {
		lsb[i] = int((*wireLabels)[c.WireCount-c.OutputSize+i][15]) & 1
	}
	return u.BitsToBytes(lsb)
}

func evaluateAndGate(g typings.Gate, wireLabels *[][]byte, garbledTable *[]byte, offset int) {

	in_a := g.InputWires[0]
	in_b := g.InputWires[1]
	out := g.OutputWire

	a := (*wireLabels)[in_a]
	b := (*wireLabels)[in_b]

	var cipher []byte
	index := 2*getColor(a) + getColor(b)

	if index == 3 { // GRR3
		cipher = make([]byte, 16)
	} else { // P&P
		offset := offset*48 + 16*index
		cipher = (*garbledTable)[offset : offset+16]
	}
	(*wireLabels)[out] = u.Decrypt(a, b, g.Id, cipher)
}

func evaluateXorGate(g typings.Gate, wireLabels *[][]byte) {
	in_a := g.InputWires[0]
	in_b := g.InputWires[1]
	out := g.OutputWire

	a := (*wireLabels)[in_a]
	b := (*wireLabels)[in_b]

	(*wireLabels)[out] = u.XorBytes(a, b)
}

func evaluateInvGate(g typings.Gate, wireLabels *[][]byte) {
	in := g.InputWires[0]
	out := g.OutputWire
	(*wireLabels)[out] = (*wireLabels)[in]
}

func getColor(label []byte) int {
	return int(label[15]) & 0x01
}
