package gc

import (
	"zkpass-node/app/typings"
	u "zkpass-node/app/utils"
)

type Evaluator struct {
	Steps        int
	Circuit      []*typings.Circuit
	GarbledTable [][]byte
}

func Evaluate(c *typings.Circuit, wireLabels *[][]byte, garbledTable *[]byte) []byte {
	counter := 0

	for i := 0; i < len(c.Gates); i++ {
		g := c.Gates[i]
		if g.Type == typings.GATE_TYPE_ADD {
			evaluateAndGate(g, wireLabels, garbledTable, counter)
			counter += 1
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

func evaluateAndGate(g typings.Gate, wireLabels *[][]byte, garbledTable *[]byte, andGateIdx int) {

	idx_a := g.InputWires[0]
	idx_b := g.InputWires[1]
	idx_c := g.OutputWire

	a := (*wireLabels)[idx_a]
	b := (*wireLabels)[idx_b]

	sa := getColor(a)
	sb := getColor(b)

	offset := andGateIdx * 32
	tg := (*garbledTable)[offset : offset+16]

	offset = offset + 16
	te := (*garbledTable)[offset : offset+16]

	hash_a := u.CCRHash(16, a)

	//generator half gate
	wg := hash_a
	if sa == 1 {
		wg = u.XorBytes(wg, tg)
	}

	//evaluator half gate
	we := u.CCRHash(16, b)
	if sb == 1 {
		we = u.XorBytes(we, u.XorBytes(te, hash_a))
	}

	//two halves make a whole
	c := u.XorBytes(wg, we)

	(*wireLabels)[idx_c] = c
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
