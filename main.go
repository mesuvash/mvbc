// package main

// import (
// 	"crypto/ecdsa"
// 	"crypto/elliptic"
// 	"crypto/rand"
// 	"fmt"
// 	"math/big"
// 	"mvbc/core/utils"

// 	"github.com/tv42/base58"
// )

// func bigJoin(expectedLen int, bigs ...*big.Int) *big.Int {

// 	bs := []byte{}
// 	for i, b := range bigs {

// 		by := b.Bytes()
// 		dif := expectedLen - len(by)
// 		if dif > 0 && i != 0 {

// 			by = append(utils.ArrayOfBytes(dif, 0), by...)
// 		}

// 		bs = append(bs, by...)
// 	}

// 	b := new(big.Int).SetBytes(bs)

// 	return b
// }

// func main() {
// 	// b := []byte{'a', 'b', 'c'}
// 	// fmt.Println(cryptoutils.SHAString(cryptoutils.SHA256(b)))
// 	pk, _ := ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
// 	//fmt.Println(pk.PublicKey.X)
// 	//fmt.Println(pk.PublicKey.Y)
// 	//fmt.Println(pk)
// 	fmt.Println(pk.PublicKey.X, pk.PublicKey.Y)
// 	b := bigJoin(28, pk.PublicKey.X, pk.PublicKey.Y)
// 	fmt.Println(b)
// 	public := base58.EncodeBig([]byte{}, b)
// 	//private := base58.EncodeBig([]byte{}, pk.D)
// 	fmt.Println(public)
// }
package main

import (
	"fmt"
	"mvbc/core/cryptoutils"
	"strconv"
)

//Staoshi style merkel tree generation
func GenerateMerkelRoot(hashes [][]byte) []byte {
	newhashes := [][]byte{}
	l := len(hashes)
	if l == 1 {
		return hashes[0]
	} else if l%2 == 1 {
		hashes = append(hashes, hashes[l-1])
	}
	for i := 0; i < len(hashes); i += 2 {
		thash := cryptoutils.SHA256(append(hashes[i], hashes[i+1]...))
		// thash := append(hashes[i], hashes[i+1]...)
		newhashes = append(newhashes, thash)
	}

	return GenerateMerkelRoot(newhashes)
}

func main() {
	var gridB []string
	grid := [][]byte{}
	ctr := 0
	for i := 0; i < 5; i++ {
		gridB = append(gridB, strconv.Itoa(ctr))
		grid = append(grid, []byte(strconv.Itoa(ctr)))
		ctr++

	}
	fmt.Println(GenerateMerkelRoot(grid))
}
