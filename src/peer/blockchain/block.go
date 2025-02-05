package blockchain

import (
	"dissy2021-group1/Hand-in9/peer/ownRSA"
	"math/big"
)

type Block struct {
	VerificationKey ownRSA.PublicKey
	SlotNumber      int
	DrawSignature   *big.Int
	Data            BlockData
	ParentHashed    string
	Signature       *big.Int
}

type BlockData struct {
	TransactionIDs []string
	Seed           *big.Int
	Hardness       *big.Int
	SpecialPKSlice []ownRSA.PublicKey
}

func MakeGenesisBlock() *Block {
	genesisBlock := new(Block)
	genesisBlock.Data.Seed = big.NewInt(42069)
	genesisBlock.Data.Hardness = big.NewInt(8273482937589649447)
	h := big.NewInt(0)
	// two := big.NewInt(2)
	// twofiftysix := big.NewInt(256)
	// ninety := big.NewInt(90)
	// hundred := big.NewInt(100)
	// h.Exp(two, twofiftysix, nil)
	// h.Div(h, hundred)
	// h.Mul(h, ninety)
	//"111006789098742178647812648156921847891275983765892138219571892389125123789927149810"
	h.SetString("111996789098742178647812648156921847891275983765892138219571892389125123789927149810", 10)
	genesisBlock.Data.Hardness = h
	return genesisBlock
}
