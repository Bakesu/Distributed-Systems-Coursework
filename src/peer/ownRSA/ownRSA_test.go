package ownRSA

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"
	"testing"
	"time"
)

func TestEncryptAndDecryptInt(t *testing.T) {
	secretKey, publicKey := KeyGen(2048)
	msg := big.NewInt(1234567890)
	encMsg := Encrypt(msg, publicKey)
	decMsg := Decrypt(encMsg, secretKey)
	fmt.Println(decMsg)
	if decMsg.Cmp(msg) != 0 {
		t.Errorf("Message wasnt the same after enc and dec")
	}
}

func TestGenerateHash(t *testing.T) {
	msg := big.NewInt(1234567890)
	msgConstant := 394652721
	hashedMsg := GenerateHash(msg)
	data := binary.BigEndian.Uint32(hashedMsg)
	// fmt.Println(data)
	if data != uint32(msgConstant) {
		t.Errorf("Nah im playing")
	}
}

func TestHashedSignature(t *testing.T) { //Test of a succesful signature
	msg := big.NewInt(1234567890)
	sk, pk := KeyGen(2048)
	signature := AssignSignature(msg, sk)
	wasVerified := VerifySignature(msg, signature, pk)
	if !wasVerified {
		t.Errorf("Signature wasnt verified")
	}
}

func TestTamperedSignatureVerification(t *testing.T) {
	msg := big.NewInt(111111199)
	sk, pk := KeyGen(2048)
	signature := AssignSignature(msg, sk)
	signature.Add(signature, big.NewInt(69))
	wasVerified := VerifySignature(msg, signature, pk)
	if wasVerified {
		t.Errorf("Signature was tampered, so shouldnt be verified")
	}
}

func TestTimeOfHash(t *testing.T) {
	msg := getRandomByteArray(10000)
	msgLength := len(msg) * 8
	msgLengthAs64Int := int64(msgLength)
	startTime := time.Now()
	sha256.Sum256(msg)
	elapsedTime := time.Since(startTime)
	fmt.Println("This is the amount of bits: ", msgLength)
	fmt.Println("This is the elapsed time in microseconds: ", elapsedTime)
	bitsPerMicroSecond := msgLengthAs64Int / elapsedTime.Microseconds()
	bitsPerSecond := bitsPerMicroSecond * 1000000
	fmt.Println("This is the number of bits per second", bitsPerSecond)
	if bitsPerSecond == 0 {
		t.Errorf("bitsPerSecond has been calculated incorrectly")
	}
}

func TestTimeOfRSASignature(t *testing.T) {
	message := big.NewInt(123456789)
	secretKey, publicKey := KeyGen(2000)
	hashedMessage := GenerateHash(message)
	hashedMsgAsBigInt := big.NewInt(0)
	hashedMsgAsBigInt = hashedMsgAsBigInt.SetBytes(hashedMessage)
	startTime := time.Now()
	hashedSignature := AssignSignature(hashedMsgAsBigInt, secretKey)
	verifiedSignature := VerifySignature(hashedMsgAsBigInt, hashedSignature, publicKey)
	elapsedTime := time.Since(startTime)
	fmt.Println("elapsed time of RSA signature: ", elapsedTime)
	if !verifiedSignature {
		t.Errorf("Signature was not verified")
	}
}

// func TimeOfHash(message []byte) int {
// 	messageAsBigInt := ConvertByteArrayToBigInt(message)
// 	startTime := time.Now()
// 	GenerateHash(messageAsBigInt)
// 	elapsedTime := time.Since(startTime)
// 	bitLength := messageAsBigInt.BitLen()
// 	if elapsedTime == 0 {
// 		return bitLength
// 	}
// 	fmt.Println("her")
// 	fmt.Println(bitLength)
// 	fmt.Println(elapsedTime)
// 	time.Sleep(time.Second * 1)
// 	bitsPerSecond := bitLength / int(elapsedTime)
// 	fmt.Println(bitsPerSecond)
// 	return bitsPerSecond
// }

func getRandomByteArray(size int) []byte {
	byteArray := make([]byte, size)
	rand.Read(byteArray)
	return byteArray
}
