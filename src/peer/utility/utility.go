package utility

import (
	"fmt"
	"math/big"
)

var BROADCAST = "broadcast"
var TRANSACTION = "transaction"
var BLOCK = "block"

//Iterates over a list []string to check whether a string already exists in the list. Returns bool
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

//Converts a string to big Int
func ConvertStringToBigInt(msg string) *big.Int {
	messageAsBigInt := new(big.Int)
	messageAsBigInt, ok := messageAsBigInt.SetString(msg, 10)
	if !ok {
		fmt.Println("SetString: error")
	}
	return messageAsBigInt
}

//Converts a big Int to string
func ConvertBigIntToString(bigInt *big.Int) string {
	stringifiedInt := bigInt.String()
	return stringifiedInt
}

func ConvertByteArrayToBigInt(byteSlice []byte) *big.Int {
	result := new(big.Int)
	result = result.SetBytes(byteSlice)
	return result
}
