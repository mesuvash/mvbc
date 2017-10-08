package main

import (
	"fmt"
	"mvbc/core/cryptoutils"
)

func main() {

	b := []byte{'h', 'e', 'l', 'l', 'o'}
	fmt.Println(cryptoutils.SHAString(cryptoutils.SHA256(b)))
}
