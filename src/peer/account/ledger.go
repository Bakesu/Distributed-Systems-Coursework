package account

import (
	"dissy2021-group1/Hand-in9/peer/ownRSA"
	"fmt"
	"math/big"
	"strings"
	"sync"
)

type Ledger struct {
	Accounts map[string]int
	lock     sync.Mutex
}

//Constructor of ledger
func MakeLedger() *Ledger {
	ledger := new(Ledger)
	ledger.Accounts = make(map[string]int)
	return ledger
}

//Converts the PublicKey struct to String
func (ledger *Ledger) PublicKeyToString(publicKey ownRSA.PublicKey) string {
	return publicKey.N.String() + "," + publicKey.E.String()
}

//Converts the PublicKey as string to a struct
func (ledger *Ledger) PublicKeyStringToStruct(s string) ownRSA.PublicKey {
	publicKey := strings.Split(s, ",")
	n := big.NewInt(0)
	e := big.NewInt(0)
	N, ok1 := n.SetString(publicKey[0], 10)
	E, ok2 := e.SetString(publicKey[1], 10)
	if !ok1 || !ok2 {
		fmt.Println("Couldnt convert from string to ownRSA.PublicKey struct")
		panic(-1)
	}
	return ownRSA.PublicKey{N: N, E: E}
}
