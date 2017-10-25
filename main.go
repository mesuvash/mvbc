package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"
	"mvbc/core/utils"

	"github.com/tv42/base58"
)

func bigJoin(expectedLen int, bigs ...*big.Int) *big.Int {

	bs := []byte{}
	for i, b := range bigs {

		by := b.Bytes()
		dif := expectedLen - len(by)
		if dif > 0 && i != 0 {

			by = append(utils.ArrayOfBytes(dif, 0), by...)
		}

		bs = append(bs, by...)
	}

	b := new(big.Int).SetBytes(bs)

	return b
}

func main() {
	// b := []byte{'a', 'b', 'c'}
	// fmt.Println(cryptoutils.SHAString(cryptoutils.SHA256(b)))
	pk, _ := ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	//fmt.Println(pk.PublicKey.X)
	//fmt.Println(pk.PublicKey.Y)
	//fmt.Println(pk)
	fmt.Println(pk.PublicKey.X, pk.PublicKey.Y)
	b := bigJoin(28, pk.PublicKey.X, pk.PublicKey.Y)
	fmt.Println(b)
	public := base58.EncodeBig([]byte{}, b)
	//private := base58.EncodeBig([]byte{}, pk.D)
	fmt.Println(public)
}
