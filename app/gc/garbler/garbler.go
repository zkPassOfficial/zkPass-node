package gc

import (
	"zkpass-node/app/typings"
	u "zkpass-node/app/utils"
)

type Garbler struct {
	Circuits []CircuitData
}

type CircuitData struct {
	InputLabels []byte
	InputBits   []int
	Masks       [][]byte
	Circuit     *typings.Circuit
}

func (g *Garbler) Garble(c *typings.Circuit) (*[]byte, *[]byte, *[]byte) {
	Δ := u.GenRandom(16)
	Δ[15] = Δ[15] | 0x01

	inputSize := c.PInputSize + c.VInputSize
	wireLabels := make([][][]byte, c.WireCount)
	copy(wireLabels, *randomInputLabels(inputSize, Δ))

	garbledTable := make([]byte, c.AndGateCount*32)

	offset := 0

	for i := 0; i < len(c.Gates); i++ {
		gate := c.Gates[i]
		if gate.Type == typings.GATE_TYPE_ADD {
			encryptedGate := garbleAndGate(gate, &wireLabels, &Δ)
			copy((garbledTable)[offset*32:(offset+1)*32], encryptedGate[0:32])
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

// P&P + FreeXor + Half-Gate
func garbleAndGate(g typings.Gate, wireLabels *[][][]byte, R *[]byte) []byte {
	idx_a := g.InputWires[0]
	idx_b := g.InputWires[1]
	idx_c := g.OutputWire

	a0 := (*wireLabels)[idx_a][0]
	a1 := (*wireLabels)[idx_a][1]

	b0 := (*wireLabels)[idx_b][0]
	b1 := (*wireLabels)[idx_b][1]

	pa := getColor(a0)
	pb := getColor(b0)

	hash_a0 := u.CCRHash(16, a0)
	hash_a1 := u.CCRHash(16, a1)

	hash_b0 := u.CCRHash(16, b0)
	hash_b1 := u.CCRHash(16, b1)

	// generator half gate
	tg := u.XorBytes(hash_a0, hash_a1)
	if pb == 1 {
		tg = u.XorBytes(tg, *R)
	}

	wg0 := hash_a0
	if pa == 1 {
		wg0 = u.XorBytes(wg0, tg)
	}

	//evaluator half gate
	te := u.XorBytes(u.XorBytes(hash_b0, hash_b1), a0)
	we0 := hash_b0
	if pb == 1 {
		we0 = u.XorBytes(we0, u.XorBytes(te, a0))
	}

	truthTable := make([][]byte, 2)
	truthTable[0] = tg
	truthTable[1] = te

	//two halves make a whole
	c0 := u.XorBytes(wg0, we0)
	c1 := u.XorBytes(c0, *R)

	(*wireLabels)[idx_c] = [][]byte{c0, c1}

	return u.Flatten(truthTable)
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
