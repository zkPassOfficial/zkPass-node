package ectf

import (
	ec "crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"math/big"
	u "zkpass-node/app/utils"

	paillier "github.com/roasbeef/go-go-gadget-paillier"
)

// Point computation in EC =》computation in F:
// At zkpass protocol, TLS session keys are divided to share of each party which are then
// used by MPC to generate an encrypted and authenticated request.
// The process of compute pre-master secret in TLS is hard to be described with bitwise operations of boolean circuits,
// so we give a simple and effitive protocol to compute shares of the pre-master secret.
// The core of the protocol is to convert shares of EC points in EC(F) to shares of coordinates in F.
type ECtF struct {
	p256 ec.Curve
	// share of EC private key
	ecPrivKey *big.Int
	// (x, y) point shares of EC public key
	ecPubKeyX, ecPubKeyY *big.Int
	// Paillier Homomorphic Encryption key
	paillierPrivKey *paillier.PrivateKey
	// P-256's Field prime
	P *big.Int
}

// Final result = A * B + C
// A = (Yn^2 − 2 * Yn * Yc + Yc^2)
// B = (Xn − Xc)^(-2) mod p= (Xn − Xc)^(p - 3)
// C = −Xn − Xc

func (p *ECtF) Init() {
	p.p256 = ec.P256()
	p.P = p.p256.Params().P

	// int in range [0, max)
	randInt, err := rand.Int(rand.Reader, u.Sub(p.p256.Params().N, big.NewInt(1)))
	if err != nil {
		panic("crypto random error")
	}

	// Node choose a random private key share ecPrivKey
	p.ecPrivKey = u.Add(randInt, big.NewInt(1))

	// Node computes a public key share ecPubKeyX = ecPrivKey * G
	p.ecPubKeyX, p.ecPubKeyY = p.p256.ScalarBaseMult(p.ecPrivKey.Bytes())

	p.paillierPrivKey = u.PaillierGenerateKey()
}

// send paillier params(n, g), server public key, E(Yn^2), E(−2*Yn)
func (p *ECtF) PreComputeA(payload []byte) ([]byte, []byte) {
	type ServerPubKey struct {
		pub_x string
		pub_y string
	}
	var srvPubKey ServerPubKey
	json.Unmarshal([]byte(string(payload)), &srvPubKey)

	// Client passes server public key to node
	serverX := u.H2Bi(srvPubKey.pub_x)
	serverY := u.H2Bi(srvPubKey.pub_y)
	serverPubkey := u.Concat([]byte{0x04}, u.To32Bytes(serverX), u.To32Bytes(serverY))

	// Node computes an EC point (xn, yn) = ecPrivKey * ServerPubKey
	xn, yn := p.p256.ScalarMult(serverX, serverY, p.ecPrivKey.Bytes())

	// encrypted xn
	eXn := u.PaillierEncrypt(p.paillierPrivKey, xn.Bytes())

	// encrypted negative xn
	eNegXn := u.PaillierEncrypt(p.paillierPrivKey, u.Sub(p.P, xn).Bytes())

	// yn^2
	pow2Yn := u.Exp(yn, big.NewInt(2), p.P)

	// encrypted yn^2
	ePow2Yn := u.PaillierEncrypt(p.paillierPrivKey, pow2Yn.Bytes())

	// -2 * yn mod p == p - 2 * yn
	neg2MulYn := u.Mod(u.Sub(p.P, u.Mul(yn, big.NewInt(2))), p.P)

	// encrypted -2 * yn mod p
	eNeg2MulYn := u.PaillierEncrypt(p.paillierPrivKey, neg2MulYn.Bytes())

	// send paillier params and enc values to client
	res := `{"e_xn":"` + hex.EncodeToString(eXn) + `",
		     "e_neg_xn":"` + hex.EncodeToString(eNegXn) + `",
			 "e_pow_2_yn":"` + hex.EncodeToString(ePow2Yn) + `",
			 "e_neg_2_mul_yn":"` + hex.EncodeToString(eNeg2MulYn) + `",
			 "n":"` + hex.EncodeToString(p.paillierPrivKey.PublicKey.N.Bytes()) + `",
			 "g":"` + hex.EncodeToString(p.paillierPrivKey.PublicKey.G.Bytes()) + `",
			 "pub_x":"` + hex.EncodeToString(p.ecPubKeyX.Bytes()) + `",
			 "pub_y":"` + hex.EncodeToString(p.ecPubKeyY.Bytes()) + `"}`

	return serverPubkey, []byte(res)
}

func (p *ECtF) PreComputeB(payload []byte) []byte {
	type ClientValue struct {
		e_b_mul_Mb_plus_Nb string // E(b*Mb+Nb)
		Nb_mod_P           string // (Nb mod p)
	}
	var cValue ClientValue
	json.Unmarshal([]byte(string(payload)), &cValue)

	// Client computes E(b*Mb+Nb) and send to Node
	// Node decrypts and gets b*Mb+Nb
	bMulMbPlusNbBytes := u.PaillierDecrypt(p.paillierPrivKey, u.H2Bi(cValue.e_b_mul_Mb_plus_Nb).Bytes())
	bMulMbPlusNb := new(big.Int).SetBytes(bMulMbPlusNbBytes)

	// Computes (b * Mb) mod p = (b * Mb + Nb) mod p - Nb mod p
	bMulMb := u.Mod(u.Sub(bMulMbPlusNb, u.H2Bi(cValue.Nb_mod_P)), p.P)

	// Computes E((b * M_b)^(p-3) mod p)
	pSub3 := u.Sub(p.P, big.NewInt(3))          // p-3
	bMulMbPowPSub3 := u.Exp(bMulMb, pSub3, p.P) // (b * M_b)^(p-3)

	// Encrypts
	eBMulMbPowPSub3 := u.PaillierEncrypt(p.paillierPrivKey, bMulMbPowPSub3.Bytes())

	json := `{"e_b_mul_mb_pow_p_sub_3":"` + hex.EncodeToString(eBMulMbPowPSub3) + `"}`

	return []byte(json)
}

func (p *ECtF) ComputeABwithMask(payload []byte) []byte {
	type ClientValue struct {
		e_B_mul_MB_plus_NB string // E(B*MB+NB)
		NB_mod_p           string // (NB mod p)
		e_A_mul_MA_plus_NA string // E(A*MA+NA)
		NA_mod_p           string // (NA mod p)
	}
	var cValue ClientValue
	json.Unmarshal([]byte(string(payload)), &cValue)

	// Client computes E(B*MB+NB) and send to Node
	// Node decrypts
	BMulMBPlusNB := new(big.Int).SetBytes(u.PaillierDecrypt(p.paillierPrivKey, u.H2Bi(cValue.e_B_mul_MB_plus_NB).Bytes()))
	// Node now gets B*MB+NB and NB mod p, computes B*MB mod p
	BMulMB := u.Mod(u.Sub(BMulMBPlusNB, u.H2Bi(cValue.NB_mod_p)), p.P)

	// Client computes E(A*MA+NA) and send to Node
	// Node decrypts
	AMulMAPlusNA := new(big.Int).SetBytes(u.PaillierDecrypt(p.paillierPrivKey, u.H2Bi(cValue.e_A_mul_MA_plus_NA).Bytes()))
	// Node now gets A*MA+NA and NA mod p, computes A*MA mod p
	AMulMA := u.Mod(u.Sub(AMulMAPlusNA, u.H2Bi(cValue.NA_mod_p)), p.P)

	// Encrypts
	eAMulMAMulBMulNB := u.PaillierEncrypt(p.paillierPrivKey, u.Mod(u.Mul(BMulMB, AMulMA), p.P).Bytes())

	json := `{"e_A_mul_MA_mul_B_mul_MB":"` + hex.EncodeToString(eAMulMAMulBMulNB) + `"}`

	return []byte(json)
}

func (p *ECtF) ComputeNodePMSShare(payload []byte) []byte {
	type ClientValue struct {
		e_A_mul_B_plus_C_plus_Sq string
	}
	var cValue ClientValue

	json.Unmarshal([]byte(string(payload)), &cValue)

	preNodePMSShare := new(big.Int).SetBytes(u.PaillierDecrypt(p.paillierPrivKey, u.H2Bi(cValue.e_A_mul_B_plus_C_plus_Sq).Bytes()))
	nodePMSShare := u.To32Bytes(u.Mod(preNodePMSShare, p.P))

	return nodePMSShare
}
