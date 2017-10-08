package core

import (
	"bytes"
	"encoding/binary"
	"errors"
	"mvbc/core/cryptoutils"
	"mvbc/core/utils"
	"reflect"
	"time"
)

type Transaction struct {
	Header    TransactionHeader
	Signature []byte
	data      []byte
}

type TransactionHeader struct {
	From       []byte
	To         []byte
	Timestamp  uint32
	DataHash   []byte
	DataLength uint32
	Nonce      uint32
}

func NewTransaction(from, to, data []byte) *Transaction {
	t := Transaction{Header: TransactionHeader{From: from, To: to}, data: data}
	t.Header.Timestamp = uint32(time.Now().Unix())
	t.Header.DataHash = cryptoutils.SHA256(t.data)
	t.Header.DataLength = uint32(len(t.data))

	return &t
}

func (t *Transaction) Hash() []byte {
	headerBytes, _ := t.Header.MarshalBinary()
	return cryptoutils.SHA256(headerBytes)
}

func (t *Transaction) Sign(keypair *Keypair) []byte {
	s, _ := keypair.Sign(t.Hash())
	return s
}

func (t *Transaction) VerifyTransaction(pow []byte) bool {
	headerHash := t.Hash()
	dataHash := cryptoutils.SHA256(t.data)
	return reflect.DeepEqual(dataHash, t.Header.DataHash) &&
		CheckProofOfWork(pow, headerHash) &&
		SignatureVerify(t.Header.From, t.Signature, headerHash)
}

func (t *Transaction) GenerateNonce(prefix []byte) uint32 {

	newT := t
	for {

		if CheckProofOfWork(prefix, newT.Hash()) {
			break
		}

		newT.Header.Nonce++
	}

	return newT.Header.Nonce
}

func (t *Transaction) MarshalBinary() ([]byte, error) {

	headerBytes, _ := t.Header.MarshalBinary()

	if len(headerBytes) != TRANSACTION_HEADER_SIZE {
		return nil, errors.New("Header marshalling error")
	}

	return append(append(headerBytes, utils.FitBytesInto(t.Signature, NETWORK_KEY_SIZE)...), t.data...), nil
}

func (t *Transaction) UnmarshalBinary(d []byte) ([]byte, error) {

	buf := bytes.NewBuffer(d)

	if len(d) < TRANSACTION_HEADER_SIZE+NETWORK_KEY_SIZE {
		return nil, errors.New("Insuficient bytes for unmarshalling transaction")
	}

	header := &TransactionHeader{}
	if err := header.UnmarshalBinary(buf.Next(TRANSACTION_HEADER_SIZE)); err != nil {
		return nil, err
	}

	t.Header = *header

	t.Signature = utils.StripByte(buf.Next(NETWORK_KEY_SIZE), 0)
	t.data = buf.Next(int(t.Header.DataLength))

	return buf.Next(MaxInt), nil

}

func (th *TransactionHeader) MarshalBinary() ([]byte, error) {

	buf := new(bytes.Buffer)

	buf.Write(utils.FitBytesInto(th.From, NETWORK_KEY_SIZE))
	buf.Write(utils.FitBytesInto(th.To, NETWORK_KEY_SIZE))
	binary.Write(buf, binary.LittleEndian, th.Timestamp)
	buf.Write(utils.FitBytesInto(th.DataHash, 32))
	binary.Write(buf, binary.LittleEndian, th.DataLength)
	binary.Write(buf, binary.LittleEndian, th.Nonce)

	return buf.Bytes(), nil
}

func (th *TransactionHeader) UnmarshalBinary(d []byte) error {

	buf := bytes.NewBuffer(d)
	th.From = utils.StripByte(buf.Next(NETWORK_KEY_SIZE), 0)
	th.To = utils.StripByte(buf.Next(NETWORK_KEY_SIZE), 0)
	binary.Read(bytes.NewBuffer(buf.Next(4)), binary.LittleEndian, &th.Timestamp)
	th.DataHash = buf.Next(32)
	binary.Read(bytes.NewBuffer(buf.Next(4)), binary.LittleEndian, &th.DataLength)
	binary.Read(bytes.NewBuffer(buf.Next(4)), binary.LittleEndian, &th.Nonce)

	return nil
}

type TransactionSlice []Transaction

func (slice TransactionSlice) Len() int {
	return len(slice)
}

func (slice TransactionSlice) Exists(tr Transaction) bool {

	for _, t := range slice {
		if reflect.DeepEqual(t.Signature, tr.Signature) {
			return true
		}
	}
	return false
}

func (slice TransactionSlice) AddTransaction(t Transaction) TransactionSlice {

	// Inserted sorted by timestamp
	for i, tr := range slice {
		if tr.Header.Timestamp >= t.Header.Timestamp {
			return append(append(slice[:i], t), slice[i:]...)
		}
	}

	return append(slice, t)
}

func (slice *TransactionSlice) MarshalBinary() ([]byte, error) {

	buf := new(bytes.Buffer)

	for _, t := range *slice {

		bs, err := t.MarshalBinary()

		if err != nil {
			return nil, err
		}

		buf.Write(bs)
	}

	return buf.Bytes(), nil
}

func (slice *TransactionSlice) UnmarshalBinary(d []byte) error {

	remaining := d

	for len(remaining) > TRANSACTION_HEADER_SIZE+NETWORK_KEY_SIZE {
		t := new(Transaction)
		rem, err := t.UnmarshalBinary(remaining)

		if err != nil {
			return err
		}
		(*slice) = append((*slice), *t)
		remaining = rem
	}
	return nil
}
