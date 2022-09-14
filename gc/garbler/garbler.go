package gc

import (
	"zkpass-node/typings"
	u "zkpass-node/utils"
)

type Garbler struct {
	Steps    int
	Circuits []CircuitData
}

type CircuitData struct {
	InputLabels []byte
	InputBits   []int
	Masks       [][]byte
	Circuit     *typings.Circuit
}

func (g *Garbler) Init(inputLabels [][][]byte, circuits []*typings.Circuit, steps int) {
	g.Steps = steps
	g.Circuits = make([]CircuitData, len(circuits))
	for i := 1; i < len(g.Circuits); i++ {
		g.Circuits[i].InputLabels = u.Concat(inputLabels[i]...)
		g.Circuits[i].Circuit = circuits[i]
	}

	g.Circuits[1].Masks = make([][]byte, 2)
	g.Circuits[1].Masks[1] = u.GenRandom(32)

	g.Circuits[2].Masks = make([][]byte, 2)
	g.Circuits[2].Masks[1] = u.GenRandom(32)

	g.Circuits[3].Masks = make([][]byte, 5)
	g.Circuits[3].Masks[1] = u.GenRandom(16)
	g.Circuits[3].Masks[2] = u.GenRandom(16)
	g.Circuits[3].Masks[3] = u.GenRandom(4)
	g.Circuits[3].Masks[4] = u.GenRandom(4)

	g.Circuits[4].Masks = make([][]byte, 3)
	g.Circuits[4].Masks[1] = u.GenRandom(16)
	g.Circuits[4].Masks[2] = u.GenRandom(16)

	g.Circuits[5].Masks = make([][]byte, 3)
	g.Circuits[5].Masks[1] = u.GenRandom(16)
	g.Circuits[5].Masks[2] = u.GenRandom(16)

	g.Circuits[7].Masks = make([][]byte, 2)
	g.Circuits[7].Masks[1] = u.GenRandom(16)

}

func (g *Garbler) Garble(c *typings.Circuit) (*[]byte, *[]byte, *[]byte) {
	Δ := u.GenRandom(16)
	Δ[15] = Δ[15] | 0x01

	inputSize := c.PInputSize + c.VInputSize
	wireLabels := make([][][]byte, c.WireCount)
	copy(wireLabels, *randomInputLabels(inputSize, Δ))

	garbledTable := make([]byte, c.AndGateCount*48)

	offset := 0

	for i := 0; i < len(c.Gates); i++ {
		gate := c.Gates[i]
		if gate.Type == typings.GATE_TYPE_ADD {
			encryptedGate := garbleAndGate(gate, &wireLabels, &Δ)
			copy((garbledTable)[offset*48:(offset+1)*48], encryptedGate[0:48])
			offset += 1
		} else if gate.Type == typings.GATE_TYPE_XOR {
			garbleXorGate(gate, &wireLabels, &Δ)
		} else if gate.Type == typings.GATE_TYPE_INV {
			garbleInvGate(gate, &wireLabels)
		}
	}

	if len(wireLabels) != c.WireCount {
		panic("len(wireLabels) != c.WireCount")
	}

	inputLabels := make([]byte, inputSize*32)
	for i := 0; i < inputSize; i++ {
		copy(inputLabels[i*32:i*32+16], wireLabels[i][0])
		copy(inputLabels[i*32+16:i*32+32], wireLabels[i][1])
	}
	lsb := make([]int, c.OutputSize)
	for i := 0; i < c.OutputSize; i++ {
		lsb[i] = int(wireLabels[c.WireCount-c.OutputSize+i][0][15]) & 1
	}
	decodingTable := u.BitsToBytes(lsb)
	return &inputLabels, &garbledTable, &decodingTable
}

func randomInputLabels(count int, Δ []byte) *[][][]byte {
	labels := make([][][]byte, count)
	for i := 0; i < count; i++ {
		label_0 := u.GenRandom(16)
		label_1 := u.XorBytes(label_0, Δ)
		labels[i] = [][]byte{label_0, label_1}
	}
	return &labels
}

// get the color bit
func getColor(label []byte) int {
	return int(label[15]) & 0x01
}

func garbleAndGate(g typings.Gate, wireLabels *[][][]byte, R *[]byte) []byte {
	in_a := g.InputWires[0]
	in_b := g.InputWires[1]
	out := g.OutputWire

	a0 := (*wireLabels)[in_a][0]
	a1 := (*wireLabels)[in_a][1]
	b0 := (*wireLabels)[in_b][0]
	b1 := (*wireLabels)[in_b][1]

	var c_0, c_1 []byte
	var rows = [4][3]*[]byte{
		{&a0, &b0, &c_0},
		{&a0, &b1, &c_0},
		{&a1, &b0, &c_0},
		{&a1, &b1, &c_1},
	}

	red_red := -1
	for i := 0; i < len(rows); i++ {
		a := *rows[i][0]
		b := *rows[i][1]

		// GRR3
		if getColor(a) == 1 && getColor(b) == 1 {
			c := make([]byte, 16)
			outWire := u.Encrypt(a, b, g.Id, c)

			if i == 3 { //TRUE outWire means true label
				c_1 = outWire
				c_0 = u.XorBytes(outWire, *R)
			} else { // FALSE  outWire means false label
				c_0 = outWire
				c_1 = u.XorBytes(outWire, *R)
			}

			red_red = i
			break
		}
	}
	(*wireLabels)[out] = [][]byte{c_0, c_1}
	if red_red == -1 {
		panic(red_red == -1)
	}

	garbledGate := make([][]byte, 3)

	for i := 0; i < len(rows); i++ {
		a := *rows[i][0]
		b := *rows[i][1]
		c := *rows[i][2]

		if i == red_red {
			continue
		}
		value := u.Encrypt(a, b, g.Id, c) // H(A, B) + C
		row := 2*getColor(a) + getColor(b)
		garbledGate[row] = value
	}
	return u.Flatten(garbledGate)
}

func garbleXorGate(g typings.Gate, wireLabels *[][][]byte, R *[]byte) {
	in_a := g.InputWires[0]
	in_b := g.InputWires[1]
	out := g.OutputWire

	c0 := u.XorBytes((*wireLabels)[in_a][0], (*wireLabels)[in_b][0])
	c1 := u.XorBytes(u.XorBytes((*wireLabels)[in_a][1], (*wireLabels)[in_b][1]), *R)
	(*wireLabels)[out] = [][]byte{c0, c1}
}

func garbleInvGate(g typings.Gate, wireLabels *[][][]byte) {
	in := g.InputWires[0]
	out := g.OutputWire
	(*wireLabels)[out] = [][]byte{(*wireLabels)[in][1], (*wireLabels)[in][0]}
}
