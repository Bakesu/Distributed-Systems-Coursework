package peer

import (
	"dissy2021-group1/Hand-in9/peer/account"
	"dissy2021-group1/Hand-in9/peer/blockchain"
	"dissy2021-group1/Hand-in9/peer/ownRSA"
	"dissy2021-group1/Hand-in9/peer/utility"
	"encoding/gob"
	"fmt"
	"math/big"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Peer struct {
	ip                   string
	port                 string
	networkSlice         []string
	encoderSlice         []*gob.Encoder
	decoderSlice         []*gob.Decoder
	envelopeChannel      chan Envelope
	prevBroadcasts       map[string]bool
	prevTransactions     map[string]bool
	prevSeenBlocks       map[string]bool
	ledger               *account.Ledger
	secretKey            ownRSA.SecretKey
	publicKey            ownRSA.PublicKey
	unqueuedTransactions []*account.SignedTransaction
	queuedTransactionIDs []string
	doneTransactions     []*account.SignedTransaction
	mutex                sync.Mutex
	blockSlice           []blockchain.Block
	slotCounter          int
	winnerFoundInSlot    bool
}

type Envelope struct { //sub structs eks til networkslice, sequencerPk og ledgerAccounts
	EnvelopeID        string
	Message           string
	NetworkSlice      []string
	IpAndPort         string
	LedgerAccounts    map[string]int
	PublicKey         ownRSA.PublicKey
	SignedTransaction *account.SignedTransaction
	Block             blockchain.Block
}

var WINNERCOUNTER int
var WINNERCOUNTERMUTEX sync.Mutex

var STARTINGMONEY int = 1000000
var SLOTLENGTH time.Duration = 1

// Constructor of Peer
func MkPeer(ip string, port string) *Peer {
	peer := new(Peer)
	peer.networkSlice = make([]string, 0)
	peer.encoderSlice = make([]*gob.Encoder, 0)
	peer.decoderSlice = make([]*gob.Decoder, 0)
	peer.envelopeChannel = make(chan Envelope)
	peer.prevBroadcasts = make(map[string]bool)
	peer.prevTransactions = make(map[string]bool)
	peer.prevSeenBlocks = make(map[string]bool)
	peer.ledger = account.MakeLedger()
	peer.unqueuedTransactions = make([]*account.SignedTransaction, 0)
	peer.queuedTransactionIDs = make([]string, 0)
	peer.blockSlice = make([]blockchain.Block, 0)
	peer.slotCounter = 0
	genesisBlock := blockchain.MakeGenesisBlock()
	peer.blockSlice = append(peer.blockSlice, *genesisBlock)
	peer.winnerFoundInSlot = false

	peer.secretKey, peer.publicKey = ownRSA.KeyGen(2048)

	var wg sync.WaitGroup

	go peer.HandleEnvelopeChannel()

	wg.Add(1)
	conn := peer.ConnectTCP(ip, port)
	if conn == nil {
		fmt.Println("No server available, creating new network")
		go peer.CreateListener(&wg)
		wg.Wait() //waits for done in createlistnere() -> peer.ip and peer.port is updated
		ipAndPort := peer.ip + ":" + peer.port
		peer.networkSlice = append(peer.networkSlice, ipAndPort)
		peer.mutex.Lock()
		peer.blockSlice[0].Data.SpecialPKSlice = append(peer.blockSlice[0].Data.SpecialPKSlice, peer.publicKey)
		peer.mutex.Unlock()
		pkAsString := peer.ledger.PublicKeyToString(peer.publicKey)
		peer.ledger.Accounts[pkAsString] = STARTINGMONEY
		//go peer.lottery()
	} else {
		peer.networkSlice, peer.ledger.Accounts, peer.blockSlice[0].Data.SpecialPKSlice = peer.ReceiveNetworkData(conn)
		go peer.ConnectionListener(conn)
		peer.connectToTenPeers()
		go peer.CreateListener(&wg)
		wg.Wait()
		ipAndPort := peer.ip + ":" + peer.port
		peer.AddEnvelopeToEnvelopeChannel(&Envelope{EnvelopeID: utility.BROADCAST, IpAndPort: ipAndPort, PublicKey: peer.publicKey}) //Broadcasting presence // Vi parser en struct
		//go peer.lottery()                                                                                                            //placeret rigtigt?
	}
	return peer
}

//Connects the peer to the last 10 peers
func (peer *Peer) connectToTenPeers() {
	sliceLen := len(peer.networkSlice)

	var loopIterations int
	if sliceLen > 10 {
		loopIterations = sliceLen - 10
	} else {
		loopIterations = 0
	}
	peer.mutex.Lock()
	for i := sliceLen; i > loopIterations; i-- {
		s := strings.Split(peer.networkSlice[i-1], ":")
		conn := peer.ConnectTCP("", s[len(s)-1]) //Crappy løsning ":::123456".
		peer.encoderSlice = append(peer.encoderSlice, gob.NewEncoder(conn))
		go peer.ConnectionListener(conn)
	}
	peer.mutex.Unlock()
}

//Hosts a tcp connection and listens for new connection. When a peer connects, it sends a networkslice and calls peer.ConnectionListener()
func (peer *Peer) CreateListener(wg *sync.WaitGroup) {
	ln, _ := net.Listen("tcp", "")
	defer ln.Close()
	ip, port, _ := net.SplitHostPort(ln.Addr().String())
	fmt.Println("hosting on: ", ip, ":", port)
	peer.ip = ip
	peer.port = port
	wg.Done()
	for {
		conn, _ := ln.Accept() //Blocking
		enc := gob.NewEncoder(conn)
		peer.mutex.Lock()
		data := blockchain.BlockData{SpecialPKSlice: peer.blockSlice[0].Data.SpecialPKSlice}
		block := blockchain.Block{Data: data}
		envelope := Envelope{NetworkSlice: peer.networkSlice, LedgerAccounts: peer.ledger.Accounts, Block: block} //not doing the EnvelopeID thing
		err := enc.Encode(envelope)
		if err != nil {
			fmt.Println("Error: " + err.Error())
			return
		}
		peer.encoderSlice = append(peer.encoderSlice, enc)
		dec := gob.NewDecoder(conn)
		peer.decoderSlice = append(peer.decoderSlice, dec)
		peer.mutex.Unlock()
		go peer.ConnectionListener(conn)
	}
}

//Recieves data about the network (networkslice, ledger and sequencerPK) and adds the newly created decoder to decoderSlice. Returns networksSlice
func (peer *Peer) ReceiveNetworkData(conn net.Conn) ([]string, map[string]int, []ownRSA.PublicKey) {
	dec := gob.NewDecoder(conn)

	var envelope Envelope
	err := dec.Decode(&envelope)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return nil, nil, nil
	}
	peer.decoderSlice = append(peer.decoderSlice, dec)
	if envelope.NetworkSlice != nil {
		return envelope.NetworkSlice, envelope.LedgerAccounts, envelope.Block.Data.SpecialPKSlice
	}
	fmt.Println("networkslice was empty (should be an error)")
	return nil, nil, nil
}

//Handles the TCP connection where it listens for envelopes and adds them to the envelopeChannel
func (peer *Peer) ConnectionListener(conn net.Conn) {
	dec := peer.decoderSlice[len(peer.decoderSlice)-1]
	defer conn.Close()
	for {
		var envelope Envelope
		err := dec.Decode(&envelope)
		if err != nil {
			fmt.Println("Error: " + err.Error())
			return
		}
		peer.envelopeChannel <- envelope
	}
}

//Reads from the envelopeChannel and checks what the envelope contains and acts upon it.
func (peer *Peer) HandleEnvelopeChannel() {
	for {
		envelope := <-peer.envelopeChannel //blocking

		switch envelope.EnvelopeID {
		case utility.BROADCAST:
			address := envelope.IpAndPort
			if !peer.prevBroadcasts[address] || !utility.StringInSlice(address, peer.networkSlice) {
				peer.prevBroadcasts[address] = true
				peer.networkSlice = append(peer.networkSlice, address)
				pkAsString := peer.ledger.PublicKeyToString(envelope.PublicKey)
				if len(peer.networkSlice) <= 10 {
					peer.ledger.Accounts[pkAsString] = STARTINGMONEY
					peer.blockSlice[0].Data.SpecialPKSlice = append(peer.blockSlice[0].Data.SpecialPKSlice, envelope.PublicKey)
					if len(peer.blockSlice[0].Data.SpecialPKSlice) == 10 {
						go peer.lottery()
					}
				} else {
					peer.ledger.Accounts[pkAsString] = 0
				}
				peer.broadcastOnEncoders(envelope)
			}
		case utility.TRANSACTION:
			if !peer.prevTransactions[envelope.SignedTransaction.ID] {
				peer.prevTransactions[envelope.SignedTransaction.ID] = true
				peer.unqueuedTransactions = append(peer.unqueuedTransactions, envelope.SignedTransaction)
				peer.broadcastOnEncoders(envelope)
			}
		case utility.BLOCK:
			block := envelope.Block
			encodedBlock := peer.encodeBlock(block)
			if !ownRSA.VerifySignature(encodedBlock, block.Signature, block.VerificationKey) {
				fmt.Println("Block was not verified")
				break
			}
			if !peer.prevSeenBlocks[peer.BlockAsHashedString(block)] {
				peer.prevSeenBlocks[peer.BlockAsHashedString(block)] = true
				lastBlock := peer.blockSlice[len(peer.blockSlice)-1]
				lastBlockHashed := peer.BlockAsHashedString(lastBlock)
				pkAsString := peer.ledger.PublicKeyToString(block.VerificationKey)
				peer.ledger.Accounts[pkAsString] = peer.ledger.Accounts[pkAsString] + 10 + len(block.Data.TransactionIDs)

				if peer.winnerFoundInSlot {
					fmt.Println("More than one winner was found in a slot")
					peer.insertIntoSecondLastIndexInBlockSlice(block)
					break
				}
				peer.winnerFoundInSlot = true

				if lastBlockHashed != block.ParentHashed {
					//Check for longest chain
					fmt.Println("Dang sumting wrong - rollback?")
					fmt.Println("lastblock: ", lastBlock.SlotNumber)
					fmt.Println("Parenthas: ", block.SlotNumber)
					fmt.Println(len(peer.blockSlice))
					parentIndex := peer.findBlockSliceIndex(block.ParentHashed)
					if parentIndex == -1 {
						fmt.Println("ParentHash of block doesnt exist in BlockSlice - which is bad")
						break
					}
					peer.rollback(parentIndex)
				}

				peer.blockSlice = append(peer.blockSlice, block)
				var DoneTransactionsIndex []int
				blockTransactionsIDs := append(peer.queuedTransactionIDs, envelope.Block.Data.TransactionIDs...)
				for indexID, transactionID := range blockTransactionsIDs {
					matchFound := false
					for indexT, unqueuedTransaction := range peer.unqueuedTransactions {
						if transactionID == unqueuedTransaction.ID {
							fmt.Println("Match!")
							matchFound = true
							peer.ledger.SignedTransaction(peer.unqueuedTransactions[indexT])
							DoneTransactionsIndex = append(DoneTransactionsIndex, indexT)
						}
					}
					if !matchFound {
						fmt.Println("No match was found")
						//append resterende tal to queuedTransaction som skal køres første gang i næste blok.
						for i := indexID; i < len(blockTransactionsIDs); i++ {
							peer.queuedTransactionIDs = append(peer.queuedTransactionIDs, blockTransactionsIDs[i])
						}
						break
					}
					for _, index := range DoneTransactionsIndex {
						peer.doneTransactions = append(peer.unqueuedTransactions, peer.unqueuedTransactions[index])
						removeFromSignedTransactionSlice(peer.unqueuedTransactions, index)
					}
				}
				peer.broadcastOnEncoders(envelope)
			}
		default:
			fmt.Println("No valid envelopeID")
		}
	}
}

func (peer *Peer) lottery() {
	fmt.Println("lottery() was called")
	for {
		time.Sleep(time.Second * SLOTLENGTH)
		peer.mutex.Lock()
		peer.slotCounter += 1
		peer.winnerFoundInSlot = false
		peer.mutex.Unlock()
		pkAsString := peer.ledger.PublicKeyToString(peer.publicKey)
		tickets := big.NewInt(int64(peer.ledger.Accounts[pkAsString]))
		seed := utility.ConvertBigIntToString(peer.blockSlice[0].Data.Seed)
		msg := strconv.Itoa(peer.slotCounter) + seed
		drawAsBigInt := ownRSA.AssignSignature(msg, peer.secretKey)
		draw := utility.ConvertBigIntToString(drawAsBigInt)
		hashedDrawAsByteArray := ownRSA.GenerateHash(draw)
		hashedDraw := utility.ConvertByteArrayToBigInt(hashedDrawAsByteArray)
		result := big.NewInt(0)
		result.Mul(tickets, hashedDraw) //vær opmærksom på at den gemmer det rigtige sted. We dont trust bigint
		//fmt.Println("Det her er result: ", result)
		//fmt.Println("Det her er Hardness: ", peer.blockSlice[0].Data.Hardness)
		if result.Cmp(peer.blockSlice[0].Data.Hardness) == 1 {
			fmt.Println("aaaaaah vinder, kan jeg lige låne nogle penge")
			WINNERCOUNTERMUTEX.Lock()
			WINNERCOUNTER++
			WINNERCOUNTERMUTEX.Unlock()
			peer.mutex.Lock()
			blockData := blockchain.BlockData{}
			for _, signedTransaction := range peer.unqueuedTransactions {
				blockData.TransactionIDs = append(blockData.TransactionIDs, signedTransaction.ID)
			}
			peer.unqueuedTransactions = nil //Fint bare med = nil yes?

			parent := peer.blockSlice[len(peer.blockSlice)-1]
			parentHashed := peer.BlockAsHashedString(parent)

			block := blockchain.Block{
				VerificationKey: peer.publicKey,
				SlotNumber:      peer.slotCounter,
				DrawSignature:   drawAsBigInt,
				Data:            blockData,
				ParentHashed:    parentHashed}

			encodedBlock := peer.encodeBlock(block)
			signature := ownRSA.AssignSignature(encodedBlock, peer.secretKey)
			block.Signature = signature
			envelope := Envelope{EnvelopeID: utility.BLOCK, Block: block}
			peer.AddEnvelopeToEnvelopeChannel(&envelope)
			peer.mutex.Unlock()
		} else {
			fmt.Println("lol taber")
		}

	}
}

func (peer *Peer) rollback(newLeafParentIndex int) {
	//Reversed slice of all block's index in the new tree
	var indexOfRollBack []int
	indexOfRollBack = append(indexOfRollBack, newLeafParentIndex)
	for {
		parentHash := peer.blockSlice[indexOfRollBack[0]].ParentHashed
		if parentHash == "" {
			break
		}
		parentHashIndex := peer.findBlockSliceIndex(parentHash)
		indexOfRollBack = append([]int{parentHashIndex}, indexOfRollBack...) //Making it reversed
	}
	//RESET
	genesisBlock := peer.blockSlice[0]
	specialPKSlice := genesisBlock.Data.SpecialPKSlice
	for _, PK := range specialPKSlice {
		pkAsString := peer.ledger.PublicKeyToString(PK)
		peer.ledger.Accounts[pkAsString] = STARTINGMONEY
	}
	//Reward all winners
	for _, block := range peer.blockSlice {
		pkAsString := peer.ledger.PublicKeyToString(block.VerificationKey)
		peer.ledger.Accounts[pkAsString] = peer.ledger.Accounts[pkAsString] + 10 + len(block.Data.TransactionIDs)
	}
	//Execute all transactions from new tree
	for _, blockSliceIndex := range indexOfRollBack {
		transactionIDSlice := peer.blockSlice[blockSliceIndex].Data.TransactionIDs
		for _, transactionID := range transactionIDSlice {
			for indexT, doneTransaction := range peer.doneTransactions {
				if transactionID == doneTransaction.ID {
					peer.ledger.SignedTransaction(peer.doneTransactions[indexT])
				}
			}
		}
	}
	//OPDATERER LASTBLOCK
	newLastBlock := peer.blockSlice[newLeafParentIndex]
	peer.removeFromBlockSlice(newLeafParentIndex)
	peer.blockSlice = append(peer.blockSlice, newLastBlock)
}

// func (peer *Peer) findAndExecuteTransaction(transactionID string) bool {
// 	matchFound := false
// 	for indexT, unqueuedTransaction := range peer.unqueuedTransactions {
// 		if transactionID == unqueuedTransaction.ID {
// 			fmt.Println("Match!")
// 			matchFound = true
// 			peer.ledger.SignedTransaction(peer.unqueuedTransactions[indexT])
// 		}
// 	}
// 	return matchFound
// }

func (peer *Peer) encodeBlock(block blockchain.Block) string {
	vkAsString := peer.ledger.PublicKeyToString(block.VerificationKey)
	signatureAsString := utility.ConvertBigIntToString(block.DrawSignature)
	parentHashedAsString := string(block.ParentHashed)
	encodedBlock := vkAsString + ":" + strconv.Itoa(block.SlotNumber) + ":" + signatureAsString
	for _, transactionID := range block.Data.TransactionIDs {
		encodedBlock += ":" + transactionID
	}
	encodedBlock += ":" + parentHashedAsString
	return encodedBlock
}

//Broadcasts envelope given in parameter to all encoder in peer.encoderslice
func (peer *Peer) broadcastOnEncoders(envelope Envelope) {
	for _, enc := range peer.encoderSlice {
		err := enc.Encode(envelope)
		if err != nil {
			fmt.Println("Error: " + err.Error())
			return
		}
	}
}

//Adds an envelope to the envelopeChannel
func (peer *Peer) AddEnvelopeToEnvelopeChannel(envelope *Envelope) {
	peer.envelopeChannel <- *envelope
}

//Establish connection to the given ip and port
func (peer *Peer) ConnectTCP(ip string, port string) net.Conn {
	ipAndPort := strings.TrimSpace(ip) + ":" + strings.TrimSpace(port)
	conn, err := net.Dial("tcp", ipAndPort)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return nil // Måske ikke optimalt med nil
	}
	return conn
}

func (peer *Peer) BlockAsHashedString(block blockchain.Block) string {
	blockAsString := peer.ledger.PublicKeyToString(block.VerificationKey) + ":" + strconv.Itoa(block.SlotNumber)
	parentHashed := ownRSA.GenerateHash(blockAsString)
	return string(parentHashed)
}

func removeFromSignedTransactionSlice(slice []*account.SignedTransaction, s int) []*account.SignedTransaction {
	return append(slice[:s], slice[s+1:]...)
}

func (peer *Peer) removeFromBlockSlice(s int) {
	peer.blockSlice = append(peer.blockSlice[:s], peer.blockSlice[s+1:]...)
}

func (peer *Peer) insertIntoSecondLastIndexInBlockSlice(value blockchain.Block) {
	lastBlock := peer.blockSlice[len(peer.blockSlice)-1]
	peer.blockSlice[len(peer.blockSlice)-1] = value
	peer.blockSlice = append(peer.blockSlice, lastBlock)
}

func (peer *Peer) findBlockSliceIndex(hashedValue string) int {
	for i := len(peer.blockSlice) - 1; i >= 0; i-- {
		if hashedValue == peer.BlockAsHashedString(peer.blockSlice[i]) {
			return i
		}
	}
	return -1
}
