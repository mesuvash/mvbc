package cryptoutils

import (
	//"crypto/rand"
	//"crypto/sha1"
	"crypto/sha256"
	//"encoding/base64"
	//"encoding/json"
	"fmt"
	//"io"
	//"math/big"
	//"reflect"
	//"strings"
	//"time"
)

func SHA256(data []byte) []byte {

	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}

func SHAString(data []byte) string {
	return fmt.Sprintf("%x", data)
}