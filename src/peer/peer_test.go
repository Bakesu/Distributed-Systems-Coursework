package peer

import (
	"dissy2021-group1/Hand-in9/peer/account"
	"dissy2021-group1/Hand-in9/peer/ownRSA"
	"dissy2021-group1/Hand-in9/peer/utility"
	"fmt"
	"math/big"
	"strconv"
	"testing"
	"time"
)

var firstPeer *Peer
var secondPeer *Peer
var thirdPeer *Peer
var fourthPeer *Peer
var fifthPeer *Peer
var sixthPeer *Peer
var seventhPeer *Peer
var eigthPeer *Peer
var ninthPeer *Peer
var tenthPeer *Peer
var eleventhPeer *Peer
var twelthPeer *Peer
var thirteenthPeer *Peer
var fourteenthPeer *Peer
var fifteenthPeer *Peer

// func TestMain(m *testing.M) {
// 	setup()
// 	time.Sleep(5 * time.Second)
// 	code := m.Run()
// 	os.Exit(code)
// }

func setup() {
	firstPeer = MkPeer("", "55555")
	secondPeer = MkPeer("", firstPeer.port)
	thirdPeer = MkPeer("", secondPeer.port)
	fourthPeer = MkPeer("", thirdPeer.port)
	fifthPeer = MkPeer("", fourthPeer.port)
	sixthPeer = MkPeer("", fifthPeer.port)
	seventhPeer = MkPeer("", sixthPeer.port)
	eigthPeer = MkPeer("", seventhPeer.port)
	ninthPeer = MkPeer("", eigthPeer.port)
	tenthPeer = MkPeer("", ninthPeer.port)
	eleventhPeer = MkPeer("", tenthPeer.port)
	twelthPeer = MkPeer("", eleventhPeer.port)
	thirteenthPeer = MkPeer("", twelthPeer.port)
	fourteenthPeer = MkPeer("", thirteenthPeer.port)
	fifteenthPeer = MkPeer("", fourteenthPeer.port)
	time.Sleep(1 * time.Second) //making sure that broadcast of last peer is flooding the whole system.
	fmt.Println("setup was succesful")
}

func createTransactionEnvelope(transactionID string, amount int, fromPK ownRSA.PublicKey, fromSK ownRSA.SecretKey, toPK ownRSA.PublicKey) Envelope {
	fifteenthPeer = MkPeer("", "")
	transaction := new(account.UnsignedTransaction)
	transaction.ID = transactionID
	transaction.Amount = amount
	from := fifteenthPeer.ledger.PublicKeyToString(fromPK)
	to := fifteenthPeer.ledger.PublicKeyToString(toPK)
	transaction.From = from
	transaction.To = to

	signedTransaction := fifteenthPeer.ledger.CreateSignedTransaction(*transaction, fromSK)
	envelope := Envelope{EnvelopeID: utility.TRANSACTION, SignedTransaction: &signedTransaction}
	return envelope
}

// func TestConnectToExistingPeer(t *testing.T) {
// 	setup()
// 	if len(secondPeer.connSlice) == 0 {
// 		t.Errorf("Should have connected to existing network (with TCP) and not create own network")
// 	}
// }

func TestNetworkSliceOfLastPeer(t *testing.T) {
	setup()
	if len(thirdPeer.networkSlice) != 15 || len(thirdPeer.ledger.Accounts) != 15 { //Måske skrevet forkert med anden condition
		fmt.Println(len(thirdPeer.ledger.Accounts))
		fmt.Println(len(firstPeer.ledger.Accounts))
		t.Errorf("networkSlice of thirdPeer should be 15")
	}
}
func TestNetworkSliceOfThirdPeer(t *testing.T) {
	// This test does not require the setup() method
	firstPeer := MkPeer("", "123456")
	secondPeer := MkPeer("", firstPeer.port)
	thirdPeer := MkPeer("", secondPeer.port)
	time.Sleep(1 * time.Second)
	if len(thirdPeer.networkSlice) != 3 {
		for _, adress := range thirdPeer.networkSlice {
			fmt.Println(adress)
		}
		fmt.Println(len(thirdPeer.networkSlice))
		t.Errorf("networkSlice of thirdPeer should be 3")
	}
}

// func TestConnectivityBetweenThreePeers(t *testing.T) {
// 	This test does not require the setup() method
// 	smallPeer1 := MkPeer("", "55000")
// 	smallPeer2 := MkPeer("", smallPeer1.port)
// 	smallPeer3 := MkPeer("", smallPeer2.port)
// 	smallPeer4 := MkPeer("", smallPeer3.port)
// 	if len(smallPeer4.connSlice) != 3 {
// 		t.Errorf("forthPeer should have three connections")
// 	}
// }

// func TestConnectivityAtMoreThan10Peers(t *testing.T) {
// 	setup()
// 	if len(fifteenthPeer.connSlice) != 10 {
// 		fmt.Println(len(fifteenthPeer.connSlice))
// 		t.Errorf("fifteenth should have 10 connections")
// 	}
// }

// func TestContentOfConnSlice(t *testing.T) {
// 	setup()
// 	for _, conn := range fifteenthPeer.connSlice {
// 		_, port, _ := net.SplitHostPort(conn.RemoteAddr().String())
// 		if port == firstPeer.port {
// 			t.Errorf("firstPeer.port: " + firstPeer.port + ", should not be in fifteenthPeer's connslice")
// 		}
// 	}
// }

func TestConvertingPublicKeyFromAndToString(t *testing.T) {
	setup()
	publicKeyAsString := firstPeer.ledger.PublicKeyToString(firstPeer.publicKey)
	fmt.Println(publicKeyAsString)
	publicKeyAsStruct := firstPeer.ledger.PublicKeyStringToStruct(publicKeyAsString)
	fmt.Println(publicKeyAsStruct.E)
	fmt.Println(publicKeyAsStruct.N)
	if publicKeyAsStruct.E.Cmp(big.NewInt(3)) != 0 {
		t.Errorf("Did not convert properly")
	}
}

func TestTransactionFloodingInBigNetwork(t *testing.T) {
	setup()
	fmt.Println("test string")
	transaction := new(account.UnsignedTransaction)
	transaction.ID = "1"
	transaction.Amount = 200
	from := firstPeer.ledger.PublicKeyToString(firstPeer.publicKey)
	to := secondPeer.ledger.PublicKeyToString(secondPeer.publicKey)
	transaction.From = from
	transaction.To = to

	signedTransaction := firstPeer.ledger.CreateSignedTransaction(*transaction, firstPeer.secretKey)

	envelope := Envelope{EnvelopeID: utility.TRANSACTION, SignedTransaction: &signedTransaction}

	firstPeer.AddEnvelopeToEnvelopeChannel(&envelope)

	time.Sleep(time.Second * 2)

	if fifteenthPeer.ledger.Accounts[to] != 200 {
		t.Errorf("Transaction didnt reach or at least change Alices amount to 200 i.e the network is not consistent")
	}

	transaction.ID = "2"
	signedTransaction = firstPeer.ledger.CreateSignedTransaction(*transaction, firstPeer.secretKey)
	envelope = Envelope{EnvelopeID: utility.TRANSACTION, SignedTransaction: &signedTransaction}
	firstPeer.AddEnvelopeToEnvelopeChannel(&envelope)

	time.Sleep(time.Second * 2)

	if seventhPeer.ledger.Accounts[to] != 400 {
		fmt.Println("Alices amount at seventhPeer after 2 transactions: " + strconv.Itoa(fifteenthPeer.ledger.Accounts["alice"]))
		fmt.Println("Size of seventhPeer's prevTransactions: " + strconv.Itoa(len(fifteenthPeer.prevTransactions)))
		t.Errorf("Transaction didnt reach or at least change Alices amount to 400 i.e the network is not consistent")
	}

	if fifteenthPeer.ledger.Accounts[to] != 400 {
		fmt.Println("Alices amount at fifteenthPeer after 2 transactions: " + strconv.Itoa(fifteenthPeer.ledger.Accounts["alice"]))
		fmt.Println("Size of fifteenthPeer's prevTransactions: " + strconv.Itoa(len(fifteenthPeer.prevTransactions)))
		t.Errorf("Transaction didnt reach or at least change Alices amount to 400 i.e the network is not consistent")
	}
}

func TestTamperingWithSignedTransaction(t *testing.T) {
	setup()
	transaction := new(account.UnsignedTransaction)
	transaction.ID = "1"
	transaction.Amount = 200
	from := firstPeer.ledger.PublicKeyToString(firstPeer.publicKey)
	to := secondPeer.ledger.PublicKeyToString(secondPeer.publicKey)
	transaction.From = from
	transaction.To = to

	signedTransaction := firstPeer.ledger.CreateSignedTransaction(*transaction, firstPeer.secretKey)

	//Switch of From and To field to see if it still is validated (which is shouldnt)
	temp := signedTransaction.From
	signedTransaction.From = signedTransaction.To
	signedTransaction.To = temp

	envelope := Envelope{EnvelopeID: utility.TRANSACTION, SignedTransaction: &signedTransaction}

	firstPeer.AddEnvelopeToEnvelopeChannel(&envelope)

	time.Sleep(time.Second * 7)

	if fifteenthPeer.ledger.Accounts[to] != 1000 {
		t.Errorf("Signature should have failed the verification")
	}
}

// func TestAllPeersHaveSequencerPK(t *testing.T) {
// 	setup()
// 	if firstPeer.publicKey.N.Cmp(fifteenthPeer.sequencerPK.N) != 0 {
// 		t.Errorf("FifteenthPeer's sequencerPK was not first peer's pk (Which is the sequencer)")
// 	}
// }

func TestUnqueuedTransactions(t *testing.T) {
	setup()
	transaction := new(account.UnsignedTransaction)
	transaction.ID = "1"
	transaction.Amount = 200
	from := fifthPeer.ledger.PublicKeyToString(fifthPeer.publicKey)
	to := secondPeer.ledger.PublicKeyToString(secondPeer.publicKey)
	transaction.From = from
	transaction.To = to

	signedTransaction := fifthPeer.ledger.CreateSignedTransaction(*transaction, fifthPeer.secretKey)

	envelope := Envelope{EnvelopeID: utility.TRANSACTION, SignedTransaction: &signedTransaction}

	fifthPeer.AddEnvelopeToEnvelopeChannel(&envelope)

	time.Sleep(time.Second * 2)

	if len(fifteenthPeer.unqueuedTransactions) != 1 {
		t.Errorf("FifteenthPeer should have 1 signedTransaction in unqueuedTransaction")
	}
	if len(firstPeer.unqueuedTransactions) != 1 {
		t.Errorf("firstPeer should have 1 signedTransaction in unqueuedTransaction")
	}
}

func TestSequencerPeerMakesTransaction(t *testing.T) {
	setup()
	transaction := new(account.UnsignedTransaction)
	transaction.ID = "1"
	transaction.Amount = 200
	from := firstPeer.ledger.PublicKeyToString(firstPeer.publicKey)
	to := secondPeer.ledger.PublicKeyToString(secondPeer.publicKey)
	transaction.From = from
	transaction.To = to

	signedTransaction := firstPeer.ledger.CreateSignedTransaction(*transaction, firstPeer.secretKey)

	envelope := Envelope{EnvelopeID: utility.TRANSACTION, SignedTransaction: &signedTransaction}

	fmt.Println("Adding transaction envelope")
	firstPeer.AddEnvelopeToEnvelopeChannel(&envelope)

	time.Sleep(time.Second * 7)

	if firstPeer.ledger.Accounts[from] != 800 {
		fmt.Println(firstPeer.ledger.Accounts[from])
		t.Errorf("In FirstPeer's ledger: firstPeer should have 800")
	}

	if firstPeer.ledger.Accounts[to] != 1200 {
		fmt.Println(firstPeer.ledger.Accounts[to])
		t.Errorf("In FirstPeer's ledger: SecondPeer should have 1200")
	}

	if len(firstPeer.unqueuedTransactions) != 0 {
		t.Errorf("UnqueuedTransactions should be empty")
	}
}

func TestBlockExecutionAtFifteenthPeer(t *testing.T) {
	setup()
	envelope := createTransactionEnvelope("1", 200, tenthPeer.publicKey, tenthPeer.secretKey, fifthPeer.publicKey)

	tenthPeer.AddEnvelopeToEnvelopeChannel(&envelope)

	time.Sleep(time.Second * 7)

	if fifteenthPeer.ledger.Accounts[envelope.SignedTransaction.To] != 1200 {
		fmt.Println(fifteenthPeer.ledger.Accounts[envelope.SignedTransaction.To])
		t.Errorf("In fifteenthPeer's ledger: FifthPeer should have 1200")
	}
}

func TestInvalidNegativeTransaction(t *testing.T) {
	setup()
	envelope := createTransactionEnvelope("1", 500, tenthPeer.publicKey, tenthPeer.secretKey, fifthPeer.publicKey)
	tenthPeer.AddEnvelopeToEnvelopeChannel(&envelope)

	envelope = createTransactionEnvelope("2", 500, tenthPeer.publicKey, tenthPeer.secretKey, fifthPeer.publicKey)
	tenthPeer.AddEnvelopeToEnvelopeChannel(&envelope)

	envelope = createTransactionEnvelope("3", 500, tenthPeer.publicKey, tenthPeer.secretKey, fifthPeer.publicKey)
	tenthPeer.AddEnvelopeToEnvelopeChannel(&envelope)

	time.Sleep(time.Second * 12)

	from := tenthPeer.ledger.PublicKeyToString(tenthPeer.publicKey)

	fmt.Println(fifteenthPeer.ledger.Accounts[from])
	if fifteenthPeer.ledger.Accounts[from] != 0 {
		t.Errorf("In fifteenthPeer's ledger: TenthPeer should have 0, because the third transaction is invalid")
	}
}

func TestConcurrentTransactions(t *testing.T) {
	setup()
	go executeTransactionsFromSecondPeer(1000, thirdPeer.publicKey)
	go executeTransactionsFromSecondPeer(1000, fourthPeer.publicKey)
	time.Sleep(time.Second * 50)

	from := secondPeer.ledger.PublicKeyToString(secondPeer.publicKey)
	to := secondPeer.ledger.PublicKeyToString(thirdPeer.publicKey)
	to1 := secondPeer.ledger.PublicKeyToString(fourthPeer.publicKey)
	if firstPeer.ledger.Accounts[from] != 0 {
		fmt.Println(firstPeer.ledger.Accounts[from])
		fmt.Println(firstPeer.ledger.Accounts[to])
		fmt.Println(secondPeer.ledger.Accounts[from])
		fmt.Println(secondPeer.ledger.Accounts[to])
		fmt.Println(secondPeer.ledger.Accounts[to1])
		t.Errorf("In firstPeers's ledger: Secondpeer should have 0")
	}
}

func executeTransactionsFromSecondPeer(iterations int, toPK ownRSA.PublicKey) {
	for i := 1; i <= iterations; i++ {
		envelope := createTransactionEnvelope(strconv.Itoa(i), 1, secondPeer.publicKey, secondPeer.secretKey, toPK)
		secondPeer.AddEnvelopeToEnvelopeChannel(&envelope)
		time.Sleep(time.Millisecond * 10)
	}
}

//var peerSlice []*Peer

func LefancyMKPeerhelperMethod() []*Peer {
	var peerSlice []*Peer
	peer1 := MkPeer("", "")
	peerSlice = append(peerSlice, peer1)
	for i := 1; i < 10; i++ {
		peerSlice = append(peerSlice, MkPeer("", peerSlice[i-1].port))
	}
	time.Sleep(time.Second * 3)
	return peerSlice
}

func TestAllTenPeersShouldHaveGenesis(t *testing.T) {
	peersToTestSlice := LefancyMKPeerhelperMethod()
	for _, peer := range peersToTestSlice {
		if len(peer.blockSlice) != 1 {
			t.Errorf("Block length should be 1")
		}
		if len(peersToTestSlice) != 10 {
			fmt.Println(len(peersToTestSlice))
			t.Errorf("Peerslice should contain 10 peers")
		}
	}
}

func TestSpecialPKSliceHasLength10ForAllPeers(t *testing.T) {
	peersToTestSlice := LefancyMKPeerhelperMethod()
	time.Sleep(time.Second * 8)
	if len(peersToTestSlice[0].blockSlice[0].Data.SpecialPKSlice) != 10 {
		fmt.Println(len(peersToTestSlice[0].blockSlice[0].Data.SpecialPKSlice))
		t.Errorf("In Peer1: GenesisBlock should have 10 vk in SpecialPKSlice")
	}
}

func Test1WinnerIn10Slots(t *testing.T) {
	LefancyMKPeerhelperMethod()
	time.Sleep(time.Second * 8)
	if WINNERCOUNTER > 3 {
		fmt.Println(WINNERCOUNTER)
		t.Errorf("There was more than 1 winner during 10 slots")
	}
}

func TestTransactionIsInvalidIfSmallerThan1(t *testing.T) {
	peerSlice := LefancyMKPeerhelperMethod()
	peer1 := peerSlice[0]
	peer2 := peerSlice[1]
	envelope := createTransactionEnvelope("1", 0, peer1.publicKey, peer1.secretKey, peer2.publicKey)
	peer1.AddEnvelopeToEnvelopeChannel(&envelope)
	time.Sleep(time.Second * 8)
	peer2PkAsString := peer2.ledger.PublicKeyToString(peer2.publicKey)
	if peer1.ledger.Accounts[peer2PkAsString] != 1000000 {
		t.Errorf("Transaction of 0 shouldnt have been done (neither should its transaction fee of 1)")
	}
}

func TestTransactionFee(t *testing.T) {
	peerSlice := LefancyMKPeerhelperMethod()
	peer1 := peerSlice[0]
	peer2 := peerSlice[1]
	envelope := createTransactionEnvelope("1", 10, peer1.publicKey, peer1.secretKey, peer2.publicKey)
	peer1.AddEnvelopeToEnvelopeChannel(&envelope)
	time.Sleep(time.Second * 8)
	peer1PkAsString := peer2.ledger.PublicKeyToString(peer1.publicKey)
	peer2PkAsString := peer2.ledger.PublicKeyToString(peer2.publicKey)
	if peer1.ledger.Accounts[peer1PkAsString] != 999990 || peer1.ledger.Accounts[peer2PkAsString] != 1000009 {
		fmt.Println(peer1.ledger.Accounts[peer1PkAsString])
		fmt.Println(peer1.ledger.Accounts[peer2PkAsString])
		t.Errorf("")
	}
}

func TestBlockCreationReward(t *testing.T) {

}
