package core

import (
	"bytes"
	"encoding/binary"
	"mvbc/core/cryptoutils"
	"mvbc/core/utils"
	"reflect"
)

type Block struct {
	*BlockHeader
	Signature []byte
	*TransactionSlice
}

type BlockHeader struct {
	Origin     []byte
	PrevBlock  []byte
	MerkelRoot []byte
	Timestamp  uint32
	Nonce      uint32
}

func NewBlock(previousBlock []byte) Block {
	header := &BlockHeader{PrevBlock: previousBlock}
	return Block{header, nil, new(TransactionSlice)}
}

func (b *Block) AddTransaction(t *Transaction) {
	newSlice := b.TransactionSlice.AddTransaction(*t)
	b.TransactionSlice = &newSlice
}

func (b *Block) Sign(keypair *keypair) []byte {
	s, _ := keypair.Sign(b.Hash())
	return s
}

func (b *Block) VerifyBlock(prefix []byte) bool {

	headerHash := b.Hash()
	merkel := b.GenerateMerkelRoot()

	return reflect.DeepEqual(merkel, b.BlockHeader.MerkelRoot) && CheckProofOfWork(prefix, headerHash) && SignatureVerify(b.BlockHeader.Origin, b.Signature, headerHash)
}

func (b *Block) GenerateNonce(prefix []byte) uint32 {

	newB := b
	for {

		if CheckProofOfWork(prefix, newB.Hash()) {
			break
		}

		newB.BlockHeader.Nonce++
	}

	return newB.BlockHeader.Nonce
}

func (b *Block) GenerateMerkelRoot() []byte {

	var merkell func(hashes [][]byte) []byte
	merkell = func(hashes [][]byte) []byte {

		l := len(hashes)
		if l == 0 {
			return nil
		}
		if l == 1 {
			return hashes[0]
		} else {

			if l%2 == 1 {
				return merkell([][]byte{merkell(hashes[:l-1]), hashes[l-1]})
			}

			bs := make([][]byte, l/2)
			for i, _ := range bs {
				j, k := i*2, (i*2)+1
				bs[i] = cryptoutils.SHA256(append(hashes[j], hashes[k]...))
			}
			return merkell(bs)
		}
	}

	ts := functional.Map(func(t Transaction) []byte { return t.Hash() }, []Transaction(*b.TransactionSlice)).([][]byte)
	return merkell(ts)

}

func (b *Block) Hash() []byte {
	headerHash, _ := b.BlockHeader.MarshalBinary()
	return cryptoutils.SHA256(headerHash)
}

func (h *BlockHeader) MarshalBinary() ([]byte, error) {

	buf := new(bytes.Buffer)

	buf.Write(utils.FitBytesInto(h.Origin, NETWORK_KEY_SIZE))
	binary.Write(buf, binary.LittleEndian, h.Timestamp)
	buf.Write(utils.FitBytesInto(h.PrevBlock, 32))
	buf.Write(utils.FitBytesInto(h.MerkelRoot, 32))
	binary.Write(buf, binary.LittleEndian, h.Nonce)

	return buf.Bytes(), nil
}

func (h *BlockHeader) UnmarshalBinary(d []byte) error {

	buf := bytes.NewBuffer(d)
	h.Origin = utils.StripByte(buf.Next(NETWORK_KEY_SIZE), 0)
	binary.Read(bytes.NewBuffer(buf.Next(4)), binary.LittleEndian, &h.Timestamp)
	h.PrevBlock = buf.Next(32)
	h.MerkelRoot = buf.Next(32)
	binary.Read(bytes.NewBuffer(buf.Next(4)), binary.LittleEndian, &h.Nonce)

	return nil
}

func (b *Block) MarshalBinary() ([]byte, error) {

	bhb, err := b.BlockHeader.MarshalBinary()
	if err != nil {
		return nil, err
	}
	sig := utils.FitBytesInto(b.Signature, NETWORK_KEY_SIZE)
	tsb, err := b.TransactionSlice.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return append(append(bhb, sig...), tsb...), nil
}

func (b *Block) UnmarshalBinary(d []byte) error {

	buf := bytes.NewBuffer(d)

	header := new(BlockHeader)
	err := header.UnmarshalBinary(buf.Next(BLOCK_HEADER_SIZE))
	if err != nil {
		return err
	}

	b.BlockHeader = header
	b.Signature = utils.StripByte(buf.Next(NETWORK_KEY_SIZE), 0)

	ts := new(TransactionSlice)
	err = ts.UnmarshalBinary(buf.Next(utils.MaxInt))
	if err != nil {
		return err
	}

	b.TransactionSlice = ts

	return nil
}

func (b *Block) GenerateMerkelRoot() []byte {

	var merkell func(hashes [][]byte) []byte
	merkell = func(hashes [][]byte) []byte {

		l := len(hashes)
		if l == 0 {
			return nil
		}
		if l == 1 {
			return hashes[0]
		} else {

			if l%2 == 1 {
				return merkell([][]byte{merkell(hashes[:l-1]), hashes[l-1]})
			}

			bs := make([][]byte, l/2)
			for i, _ := range bs {
				j, k := i*2, (i*2)+1
				bs[i] = helpers.SHA256(append(hashes[j], hashes[k]...))
			}
			return merkell(bs)
		}
	}

	ts := functional.Map(func(t Transaction) []byte { return t.Hash() }, []Transaction(*b.TransactionSlice)).([][]byte)
	return merkell(ts)

}
