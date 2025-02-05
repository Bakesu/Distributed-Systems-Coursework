package ownRSA

import (
	"crypto/rand"
	"crypto/sha256"
	"dissy2021-group1/Hand-in9/peer/utility"
	"math/big"
)

type SecretKey struct {
	N *big.Int
	D *big.Int
}

type PublicKey struct {
	N *big.Int
	E *big.Int
}

func KeyGen(k int) (SecretKey, PublicKey) {
	e := big.NewInt(3)
	n := big.NewInt(0)
	p := big.NewInt(0)
	q := big.NewInt(0)
	d := big.NewInt(0)
	order := big.NewInt(0)
	one := big.NewInt(1)
	pMinusOne := big.NewInt(0)
	qMinusOne := big.NewInt(0)
	gcd1 := big.NewInt(0)
	gcd2 := big.NewInt(0)

	for {
		p, _ = rand.Prime(rand.Reader, k/2+k%2)
		q, _ = rand.Prime(rand.Reader, k/2)

		pMinusOne.Sub(p, one)
		qMinusOne.Sub(q, one)

		if gcd1.GCD(nil, nil, pMinusOne, e).Cmp(one) == 0 && gcd2.GCD(nil, nil, qMinusOne, e).Cmp(one) == 0 {
			break
		}
	}

	order.Mul(pMinusOne, qMinusOne)

	n.Mul(p, q)

	if n.BitLen() == k {
		d.ModInverse(e, order)
	} else {
		panic(-1) //abort
	}

	secretKey := SecretKey{N: n, D: d}
	publicKey := PublicKey{N: n, E: e}
	return secretKey, publicKey
}

func Encrypt(msg *big.Int, pk PublicKey) *big.Int {
	encodedmessage := big.NewInt(0)
	encodedmessage = encodedmessage.Exp(msg, pk.E, pk.N)
	return encodedmessage
}

func Decrypt(encodedMessage *big.Int, sk SecretKey) *big.Int {
	decryptedMessage := big.NewInt(0)
	decryptedMessage = decryptedMessage.Exp(encodedMessage, sk.D, sk.N)
	return decryptedMessage
}

func GenerateHash(message string) []byte {
	messageAsByteArray := []byte(message)
	h := sha256.New()
	h.Write(messageAsByteArray)
	return h.Sum(nil)
}

func AssignSignature(msg string, sk SecretKey) *big.Int {
	hashedMsg := GenerateHash(msg)
	hashAsBigInt := utility.ConvertByteArrayToBigInt(hashedMsg)
	signature := big.NewInt(0)
	signature = signature.Exp(hashAsBigInt, sk.D, sk.N)
	return signature
}

func VerifySignature(msg string, signature *big.Int, pk PublicKey) bool {
	temp := big.NewInt(0)
	verifySignature := temp.Exp(signature, pk.E, pk.N)
	hashedMsg := GenerateHash(msg)
	hashAsBigInt := utility.ConvertByteArrayToBigInt(hashedMsg)
	if hashAsBigInt.Cmp(verifySignature) != 0 {
		return false
	} else {
		return true
	}
}
