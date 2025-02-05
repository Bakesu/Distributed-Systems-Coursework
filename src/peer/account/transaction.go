package account

import (
	"dissy2021-group1/Hand-in9/peer/ownRSA"
	"dissy2021-group1/Hand-in9/peer/utility"
	"fmt"
	"strconv"
	"strings"
)

type SignedTransaction struct {
	ID        string // Any string
	From      string // A verification key coded as a string
	To        string // A verification key coded as a string
	Amount    int    // Amount to transfer
	Signature string // Potential signature coded as string
}

type UnsignedTransaction struct {
	ID     string // Any string
	From   string // A verification key coded as a string
	To     string // A verification key coded as a string
	Amount int    // Amount to transfer
}

var TRANSACTIONFEE int = 1

//Processes a transaction after verifying the signature. If the verification fails, the transaction is ignored
func (l *Ledger) SignedTransaction(t *SignedTransaction) {
	l.lock.Lock()
	defer l.lock.Unlock()

	signatureAsBigInt := utility.ConvertStringToBigInt(t.Signature)
	encodedTransactionAsString := EncodeTransaction(UnsignedTransaction{ID: t.ID, From: t.From, To: t.To, Amount: t.Amount})
	fromPK := l.PublicKeyStringToStruct(t.From)
	validSignature := ownRSA.VerifySignature(encodedTransactionAsString, signatureAsBigInt, fromPK)
	if !validSignature {
		fmt.Println("Not a valid signature")
	} else if t.Amount <= 0 {
		fmt.Println("Amount was less than 1")
	} else if l.Accounts[t.From]-t.Amount < 0 {
		fmt.Println("Result would be negative for the account")
	} else {
		l.Accounts[t.From] -= t.Amount
		l.Accounts[t.To] += (t.Amount - TRANSACTIONFEE) // transaction fee of one
		// fmt.Println("Transaction was successful")
		// fmt.Println(l.Accounts[t.From])
		// fmt.Println(l.Accounts[t.To])
	}
}

//Creates a signed transaction
func (l *Ledger) CreateSignedTransaction(tempTransaction UnsignedTransaction, secretKey ownRSA.SecretKey) SignedTransaction {
	id := tempTransaction.ID
	from := tempTransaction.From
	to := tempTransaction.To
	amount := tempTransaction.Amount

	encodedTransactionAsString := EncodeTransaction(tempTransaction)
	signature := ownRSA.AssignSignature(encodedTransactionAsString, secretKey)
	signatureAsString := utility.ConvertBigIntToString(signature)
	signedTransaction := SignedTransaction{ID: id, From: from, To: to, Amount: amount, Signature: signatureAsString}

	return signedTransaction
}

//Encodes transaction to a single string (Currently very naive made)
func EncodeTransaction(transaction UnsignedTransaction) string {
	id := strings.Trim(transaction.ID, ":")
	from := strings.Split(transaction.From, ",") //This is converted a bigInt so it will never include a ":"
	to := strings.Split(transaction.To, ",")     //This is converted a bigInt so it will never include a ":"
	return id + ":" + from[0] + ":" + from[1] + ":" + to[0] + ":" + to[1] + ":" + strconv.Itoa(transaction.Amount)
}
