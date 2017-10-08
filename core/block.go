package core

type Block struct {
	*BlockHeader
	Signature []byte
	// *TransactionSlice
}

type BlockHeader struct {
	Origin     []byte
	PrevBlock  []byte
	MerkelRoot []byte
	Timestamp  uint32
	Nonce      uint32
}
