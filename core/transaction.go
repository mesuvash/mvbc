package core

import (
	"debug/elf"
	"time"
	"mbc/core/cryptoutils"
)

type Transaction struct {
	Header    TransactionHeader
	Signature []byte
	data   []byte
}

type TransactionHeader struct {
	From          []byte
	To            []byte
	Timestamp     uint32
	DataHash   []byte
	DataLength uint32
	Nonce         uint32
}


func NewTransaction(from, to, data []byte) *Transaction {
	t := Transaction{Header:TransactionHeader{From:from, To: to}, data: data}
	t.Header.Timestamp = uint32(time.Now().Unix())
	t.Header.DataHash = cryptoutils.SHA256(t.data)
	t.Header.DataLength = uint32(len(t.data))

	return &t
}

